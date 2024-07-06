package main

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "grpc-user-service/proto/github.com/360Ritik/user-service/proto"
)

var server *grpc.Server
var listener net.Listener

// Setup the gRPC server for testing
func setup() {
	listener, _ = net.Listen("tcp", ":50051")
	server = grpc.NewServer()
	pb.RegisterUserServiceServer(server, &UserServiceServer{})
	go server.Serve(listener)
}

// Tear down the gRPC server after testing
func tearDown() {
	server.Stop()
	listener.Close()
}

func TestGetUserById(t *testing.T) {
	setup()
	defer tearDown()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	tests := []struct {
		name      string
		userId    int32
		expectErr bool
		code      codes.Code
	}{
		{"ValidID", 1, false, codes.OK},
		{"InvalidID", -1, true, codes.InvalidArgument},
		{"NonExistentID", 999, true, codes.NotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.UserByIdRequest{UserId: tc.userId}
			resp, err := client.GetUserById(context.Background(), req)

			if tc.expectErr {
				assert.NotNil(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.code, st.Code())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.userId, resp.User.Id)
			}
		})
	}
}

func TestGetUsersByIds(t *testing.T) {
	setup()
	defer tearDown()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	tests := []struct {
		name      string
		userIds   []int32
		expectErr bool
		code      codes.Code
	}{
		{"ValidIDs", []int32{1, 2}, false, codes.OK},
		{"EmptyIDs", []int32{}, true, codes.InvalidArgument},
		{"NonExistentID", []int32{999}, false, codes.OK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.UsersByIdsRequest{UserIds: tc.userIds}
			resp, err := client.GetUsersByIds(context.Background(), req)

			if tc.expectErr {
				assert.NotNil(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.code, st.Code())
			} else {
				assert.Nil(t, err)
				if len(tc.userIds) == 0 {
					assert.Empty(t, resp.Users)
				} else {
					assert.NotEmpty(t, resp.Users)
				}
			}
		})
	}
}

func TestSearchUsers(t *testing.T) {
	setup()
	defer tearDown()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	tests := []struct {
		name      string
		criteria  string
		value     string
		expectErr bool
		code      codes.Code
	}{
		{"ValidCity", "city", "LA", false, codes.OK},
		{"InvalidCriteria", "unknown", "value", true, codes.InvalidArgument},
		{"EmptyValue", "city", "", true, codes.InvalidArgument},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.SearchRequest{Criteria: tc.criteria, Value: tc.value}
			resp, err := client.SearchUsers(context.Background(), req)

			if tc.expectErr {
				assert.NotNil(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.code, st.Code())
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, resp.Users)
			}
		})
	}
}

func TestAddNewUser(t *testing.T) {
	setup()
	defer tearDown()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	tests := []struct {
		name      string
		user      *pb.User
		expectErr bool
		code      codes.Code
	}{
		{"ValidUser", &pb.User{Id: 3, Fname: "John", City: "SF", Phone: 1122334455, Height: 5.9, Married: false}, false, codes.OK},
		{"InvalidID", &pb.User{Id: -1, Fname: "Invalid", City: "SF", Phone: 1122334455, Height: 5.9, Married: false}, true, codes.InvalidArgument},
		{"ShortName", &pb.User{Id: 4, Fname: "J", City: "SF", Phone: 1122334455, Height: 5.9, Married: false}, true, codes.InvalidArgument},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.AddNewUser(context.Background(), tc.user)

			if tc.expectErr {
				assert.NotNil(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.code, st.Code())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.user.Id, resp.Id)
			}
		})
	}
}
