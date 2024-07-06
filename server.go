package main

import (
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "grpc-user-service/proto/github.com/360Ritik/user-service/proto"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// In-memory user data
var users = map[int32]pb.User{
	1: {Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
	2: {Id: 2, Fname: "Alice", City: "NYC", Phone: 9876543210, Height: 5.5, Married: false},
	3: {Id: 3, Fname: "John", City: "Chicago", Phone: 5551234567, Height: 6.1, Married: true},
	4: {Id: 4, Fname: "Emily", City: "Seattle", Phone: 9998887777, Height: 5.6, Married: false},
	5: {Id: 5, Fname: "Michael", City: "San Francisco", Phone: 1112223333, Height: 6.0, Married: true},
	6: {Id: 6, Fname: "Sophia", City: "Boston", Phone: 4445556666, Height: 5.7, Married: false},
	7: {Id: 7, Fname: "David", City: "Denver", Phone: 7778889999, Height: 5.9, Married: true},
	8: {Id: 8, Fname: "Emma", City: "Houston", Phone: 2223334444, Height: 5.4, Married: false},
}

// GetUserById retrieves a user by ID
func (s *UserServiceServer) GetUserById(ctx context.Context, req *pb.UserByIdRequest) (*pb.UserResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "User ID must be positive")
	}
	user, ok := users[req.UserId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "User with ID %d not found", req.UserId)
	}
	return &pb.UserResponse{User: &user}, nil
}

// GetUsersByIds retrieves multiple users by their IDs
func (s *UserServiceServer) GetUsersByIds(ctx context.Context, req *pb.UsersByIdsRequest) (*pb.UsersResponse, error) {
	if len(req.UserIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "User IDs cannot be empty")
	}

	var result []*pb.User
	for _, id := range req.UserIds {
		if user, ok := users[id]; ok {
			// Create a new pb.User object and assign values from users[id]
			newUser := &pb.User{
				Id:      user.Id,
				Fname:   user.Fname,
				City:    user.City,
				Phone:   user.Phone,
				Height:  user.Height,
				Married: user.Married,
			}
			result = append(result, newUser)
		}
	}

	// Check if any users were found
	if len(result) == 0 {
		return nil, status.Errorf(codes.NotFound, "No users found for the provided IDs")
	}

	return &pb.UsersResponse{Users: result}, nil
}

// SearchUsers searches for users based on criteria
func (s *UserServiceServer) SearchUsers(ctx context.Context, req *pb.SearchRequest) (*pb.UsersResponse, error) {
	if req.Criteria == "" || req.Value == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Criteria and Value must be specified")
	}
	var result []*pb.User
	for _, user := range users {
		switch req.Criteria {
		case "city":
			if user.City == req.Value {
				result = append(result, &user)
			}
		case "phone":
			value, err := strconv.ParseInt(req.Value, 10, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "Invalid phone number format")
			}
			if user.Phone == value {
				result = append(result, &user)
			}
		case "married":
			if (req.Value == "true" && user.Married) || (req.Value == "false" && !user.Married) {
				result = append(result, &user)
			}
		case "fname":
			if req.Value == user.Fname {
				result = append(result, &user)
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "Unknown search criteria")
		}
	}
	return &pb.UsersResponse{Users: result}, nil
}

// AddNewUser adds a new user
func (s *UserServiceServer) AddNewUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "User ID must be positive")
	}
	if req.Fname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "First name cannot be empty")
	}
	if len(req.Fname) < 5 {
		return nil, status.Errorf(codes.InvalidArgument, "First name must be at least 5 characters long")
	}
	if req.City == "" {
		return nil, status.Errorf(codes.InvalidArgument, "City cannot be empty")
	}
	if req.Phone <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Phone number must be positive")
	}
	if req.Height <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Height must be positive")
	}

	// Check if the user ID already exists
	if _, exists := users[req.Id]; exists {
		return nil, status.Errorf(codes.AlreadyExists, "User with ID %d already exists", req.Id)
	}

	// Add the user to the in-memory map
	users[req.Id] = *req

	return req, nil
}

// Main function to start the gRPC server
func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &UserServiceServer{})

	log.Println("Starting gRPC server on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
