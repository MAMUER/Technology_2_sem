package grpc

import (
	"context"
	"errors"

	pb "tech-ip-sem2/proto/gen/go/auth"
	"tech-ip-sem2/services/auth/internal/service"
	
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthServer(authService *service.AuthService) *AuthServer {
	return &AuthServer{
		authService: authService,
	}
}

func (s *AuthServer) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is required")
	}

	valid, subject := s.authService.ValidateToken(req.Token)
	if !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &pb.VerifyResponse{
		Valid:   true,
		Subject: subject,
	}, nil
}

func RegisterAuthServiceServer(s *grpc.Server, authService *service.AuthService) {
	pb.RegisterAuthServiceServer(s, NewAuthServer(authService))
}