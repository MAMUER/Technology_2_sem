package authclient

import (
	"context"
	"fmt"
	"time"

	pb "tech-ip-sem2/proto/gen/go/auth"
	"tech-ip-sem2/shared/logger"
	"tech-ip-sem2/shared/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn       *grpc.ClientConn
	authClient pb.AuthServiceClient
	timeout    time.Duration
	log        *logger.Logger
}

func NewClient(addr string, timeout time.Duration, log *logger.Logger) (*Client, error) {
	log.Info("Connecting to auth gRPC server", zap.String("addr", addr))

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)
	log.Info("Successfully connected to auth gRPC server")

	return &Client{
		conn:       conn,
		authClient: client,
		timeout:    timeout,
		log:        log,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) VerifyToken(ctx context.Context, token string) (bool, string, error) {
	requestID := middleware.GetRequestID(ctx)
	log := c.log.WithRequestID(requestID)

	log.Debug("Calling gRPC verify", zap.String("token_prefix", token[:min(10, len(token))]))

	ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", requestID)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.authClient.Verify(ctx, &pb.VerifyRequest{
		Token: token,
	})

	if err != nil {
		log.Error("gRPC verify error", zap.Error(err))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				return false, "", fmt.Errorf("auth service timeout")
			case codes.Unavailable:
				return false, "", fmt.Errorf("auth service unavailable")
			case codes.Unauthenticated:
				log.Info("Token is invalid")
				return false, "", nil
			default:
				return false, "", fmt.Errorf("auth service error: %v", st.Message())
			}
		}
		return false, "", fmt.Errorf("failed to verify token: %w", err)
	}

	log.Info("gRPC verify success",
		zap.Bool("valid", resp.Valid),
		zap.String("subject", resp.Subject),
	)

	return resp.Valid, resp.Subject, nil
}
