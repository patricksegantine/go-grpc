# Go gRPC

This is a simple example that shows how to use gRPC with Go.

## Services implemented on Server


```
service UserService {
    rpc AddUser (User) returns (User);
    rpc AddUserVerbose (User) returns (stream UserResponse);
    rpc AddUsers (stream User) returns (Users);
    rpc AddUsersBothStream(stream User) returns (stream UserResponse);
}
```


## How to start?

### Server

```
go run cmd/server/server.go
```


### client

```
go run cmd/client/client.go
```

To run an specific example uncomment at least one these methods AddUser(s) below:

```go
func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserBothStream(client)
}
```