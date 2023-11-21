package gapi

import (
	"context"
	"database/sql"
	"db"
	"pb"
	"utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	if err := utils.VerifyHashPassword(req.Password, user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	mtdt := server.extractMetadata(ctx)
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpireAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	resp := &pb.LoginUserResponse{
		SessionId:    refreshPayload.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken}
	return resp, nil
}
