package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/patricksegantine/go-grpc/pb"
	uuid "github.com/satori/go.uuid"
)

// UserServiceImpl a implementation of the gRPC for UserService Server
type UserServiceImpl struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (*UserServiceImpl) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	id := getID()
	// Insert - Database
	fmt.Printf("ID: %v | Name: %v\n", id, req.Name)

	return &pb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

func (*UserServiceImpl) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResponse{
		Status: "Init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResponse{
		Status: "Inserting",
		User:   req,
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResponse{
		Status: "Completed",
		User: &pb.User{
			Id:    getID(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	return nil
}

func (*UserServiceImpl) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				Users: users,
			})
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		user := &pb.User{
			Id:    getID(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		}
		fmt.Println("Adding", user)

		users = append(users, user)
	}
}

func (*UserServiceImpl) AddUsersBothStream(stream pb.UserService_AddUsersBothStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error receiving stream from the client: %v", err)
		}

		err = stream.Send(&pb.UserResponse{
			Status: "Added",
			User:   req,
		})
		if err != nil {
			log.Fatalf("Error sending stream to the cliente?: %v", err)
		}
	}
}

func getID() string {
	var mu sync.Mutex

	mu.Lock()
	id := uuid.NewV4().String()
	mu.Unlock()

	return id
}
