package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/patricksegantine/go-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	AddUsers(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Name:  "Patrick Segantine",
		Email: "patricksegantine@email.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

// AddUserVerbose server stream
func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receibe the msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, stream.GetUser())
	}
}

// AddUsers implementa um client stream
func AddUsers(client pb.UserServiceClient) {
	req := []*pb.User{
		{
			Name:  "Patrick",
			Email: "patrick@email.com",
		},
		{
			Name:  "Raquel",
			Email: "raquel@email.com",
		},
		{
			Name:  "Antonella",
			Email: "antonella@email.com",
		},
		{
			Name:  "Gabriel",
			Email: "biel@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, v := range req {
		stream.Send(v)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserBothStream(client pb.UserServiceClient) {
	users := []*pb.User{
		{},
	}

	stream, err := client.AddUsersBothStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	wait := make(chan bool)

	go func() {
		for _, v := range users {
			fmt.Println("Sending user:", v.GetName())
			stream.Send(v)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status %v", res.GetUser(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
