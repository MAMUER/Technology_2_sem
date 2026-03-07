package grpc

import (
	"context"

	pb "tech-ip-sem2/proto/gen/go/auth"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
	log         *logger.Logger
}

func NewAuthServer(authService *service.AuthService, log *logger.Logger) *AuthServer {
	return &AuthServer{
		authService: authService,
		log:         log,
	}
}

func (s *AuthServer) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	var requestID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get("x-request-id"); len(values) > 0 {
			requestID = values[0]
		}
	}

	log := s.log.WithRequestID(requestID)

	if req.Token == "" {
		log.Warn("empty token in request")
		return nil, status.Error(codes.Unauthenticated, "token is required")
	}

	log.Debug("verifying token")

	valid, subject := s.authService.ValidateToken(req.Token)
	if !valid {
		log.Info("invalid token")
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	log.Info("token verified", zap.String("subject", subject))

	return &pb.VerifyResponse{
		Valid:   true,
		Subject: subject,
	}, nil
}

func RegisterAuthServiceServer(s *grpc.Server, authService *service.AuthService, log *logger.Logger) {
	pb.RegisterAuthServiceServer(s, NewAuthServer(authService, log))
}
