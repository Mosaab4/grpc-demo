package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"pcbook/pb"
	"pcbook/sample"
	"pcbook/serializer"
	"testing"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.Newlaptop()
	expectedId := laptop.Id

	req := &pb.CreateLaptopRequest{Laptop: laptop}

	res, err := laptopClient.CreateLaptop(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.Id, expectedId)

	other, err := laptopServer.Store.Find(res.Id)

	require.NoError(t, err)
	require.NotNil(t, other)

	requireSameLaptop(t, laptop, other)
}

func newTestLaptopClient(t *testing.T, address string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)

	return pb.NewLaptopServiceClient(conn)
}

func startTestLaptopServer(t *testing.T) (*LaptopServer, string) {
	laptopServer := NewLaptopServer(NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")

	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}

func requireSameLaptop(t *testing.T, laptop1, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtobufToJson(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJson(laptop2)

	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
