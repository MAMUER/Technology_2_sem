package authclient

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "tech-ip-sem2/proto/gen/go/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn       *grpc.ClientConn
	authClient pb.AuthServiceClient
	timeout    time.Duration
}

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	log.Printf("Connecting to auth gRPC server at %s", addr)

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)
	log.Printf("Successfully connected to auth gRPC server")

	return &Client{
		conn:       conn,
		authClient: client,
		timeout:    timeout,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) VerifyToken(ctx context.Context, token string) (bool, string, error) {
	log.Printf("Calling gRPC verify for token: %s...", token[:10])

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.authClient.Verify(ctx, &pb.VerifyRequest{
		Token: token,
	})

	if err != nil {
		log.Printf("gRPC verify error: %v", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				return false, "", fmt.Errorf("auth service timeout")
			case codes.Unavailable:
				return false, "", fmt.Errorf("auth service unavailable")
			case codes.Unauthenticated:
				log.Printf("Token is invalid")
				return false, "", nil
			default:
				return false, "", fmt.Errorf("auth service error: %v", st.Message())
			}
		}
		return false, "", fmt.Errorf("failed to verify token: %w", err)
	}

	log.Printf("gRPC verify success: valid=%v, subject=%s", resp.Valid, resp.Subject)
	return resp.Valid, resp.Subject, nil
}
