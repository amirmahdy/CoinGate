package gapi

import (
	"context"
	"db"
	"pb"
	"utils"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
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

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := utils.ValidateFullname(req.FullName); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := utils.ValidateUsername(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := utils.ValidateUsername(req.Password); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return violations
}
