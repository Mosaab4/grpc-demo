package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"os"
	"path/filepath"
	"pcbook/pb"
	"pcbook/sample"
	"pcbook/serializer"
	"testing"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()
	laptopStore := NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id

	req := &pb.CreateLaptopRequest{Laptop: laptop}

	res, err := laptopClient.CreateLaptop(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.Id, expectedId)

	other, err := laptopStore.Find(res.Id)

	require.NoError(t, err)
	require.NotNil(t, other)

	requireSameLaptop(t, laptop, other)
}

func newTestLaptopClient(t *testing.T, address string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)

	return pb.NewLaptopServiceClient(conn)
}

func startTestLaptopServer(t *testing.T, lapStore LaptopStore, imageStore ImageStore) string {
	laptopServer := NewLaptopServer(lapStore, imageStore)

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")

	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

func requireSameLaptop(t *testing.T, laptop1, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtobufToJson(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJson(laptop2)

	require.NoError(t, err)

	require.Equal(t, json1, json2)
}

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	laptopStore := NewInMemoryLaptopStore()
	expectedIds := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Memory = &pb.Memory{
				Value: 4096,
				Unit:  pb.Memory_GIGABYTE,
			}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Memory = &pb.Memory{
				Value: 16,
				Unit:  pb.Memory_GIGABYTE,
			}
			expectedIds[laptop.GetId()] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 5
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = 5.0
			laptop.Memory = &pb.Memory{
				Value: 64,
				Unit:  pb.Memory_GIGABYTE,
			}
			expectedIds[laptop.GetId()] = true
		}

		err := laptopStore.Save(laptop)
		require.NoError(t, err)
	}

	serverAddress := startTestLaptopServer(t, laptopStore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)

	require.NoError(t, err)

	found := 0

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIds, res.GetLaptop().GetId())
		found += 1
	}

	require.Equal(t, len(expectedIds), found)
}

func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	testTempFolder := "../tmp"

	laptopStore := NewInMemoryLaptopStore()
	imageStore := NewDiskImageStore(testTempFolder)

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)

	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, imageStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/laptop.jpg", testTempFolder)

	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())

	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptop.GetId(),
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()

	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s/%s%s", testTempFolder, res.GetId(), imageType)
	require.FileExists(t, savedImagePath)

	require.NoError(t, os.Remove(savedImagePath))
}
