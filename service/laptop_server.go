package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"pcbook/pb"
)

type LaptopServer struct {
	Store LaptopStore
	pb.UnimplementedLaptopServiceServer
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}
func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()

		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}

		laptop.Id = id.String()
	}

	if ctx.Err() == context.Canceled {
		log.Println("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	err := server.Store.Save(laptop)

	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrorAlreadyExist) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save laptop to the server: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.Id)

	return &pb.CreateLaptopResponse{Id: laptop.Id}, nil
}
