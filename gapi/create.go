package gapi

import (
	"context"
	"db"
	"pb"
	"utils"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := utils.CreateHashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error in creating hash password")
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "could not create user")
	}
	resp := &pb.CreateUserResponse{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
	}
	return resp, nil
}
