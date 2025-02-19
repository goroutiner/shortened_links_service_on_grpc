package handlers_test

import (
	"context"
	"errors"
	"net"
	"shortened_links_service_on_grpc/internal/handlers"
	pb "shortened_links_service_on_grpc/proto"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestGetShortLinkUnit проверяет работу метода GetShortLink хендлера с подключением к gRPC серверу
func TestGetShortLinkIntegration(t *testing.T) {
	listener, err := net.Listen("tcp", ":8080")
	require.NoError(t, err)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handlers.RateLimitInterceptor()))
	mockService := new(handlers.MockShortenerService)
	handler := handlers.RegisterShortenerHandler(mockService)

	pb.RegisterShortenerServiceServer(grpcServer, handler)

	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()
	defer grpcServer.Stop()

	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewShortenerServiceClient(conn)

	testOriginalLink := "https://example.com"
	expectedShortLink := "abcd123456"

	mockService.On("GetShortLink", testOriginalLink).Return(expectedShortLink, nil)

	req := &pb.GetShortLinkRequest{OriginalLink: testOriginalLink}
	resp, err := client.GetShortLink(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expectedShortLink, resp.ShortLink)

	testOriginalLink = ""

	mockService.On("GetShortLink", testOriginalLink).Return("", errors.New("invalid URL format"))
	req = &pb.GetShortLinkRequest{OriginalLink: testOriginalLink}
	_, err = client.GetShortLink(context.Background(), req)
	require.Error(t, err)

	testOriginalLink = "https:::|/example.com"

	mockService.On("GetShortLink", testOriginalLink).Return("", errors.New("invalid URL format"))
	req = &pb.GetShortLinkRequest{OriginalLink: testOriginalLink}
	_, err = client.GetShortLink(context.Background(), req)
	require.Error(t, err)

	mockService.AssertExpectations(t)
}

// TestGetOriginalLinkUnit проверяет работу метода GetOriginalLink хендлера с подключением к gRPC серверу
func TestGetOriginalLinkIntegration(t *testing.T) {
	listener, err := net.Listen("tcp", ":8080")
	require.NoError(t, err)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handlers.RateLimitInterceptor()))
	mockService := new(handlers.MockShortenerService)
	handler := handlers.RegisterShortenerHandler(mockService)

	pb.RegisterShortenerServiceServer(grpcServer, handler)

	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()
	defer grpcServer.Stop()

	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewShortenerServiceClient(conn)

	testShortLink := "abcd123456"
	expectedOriginalLink := "https://example.com"

	mockService.On("GetOriginalLink", testShortLink).Return(expectedOriginalLink, nil)

	req := &pb.GetOriginalLinkRequest{ShortLink: testShortLink}
	resp, err := client.GetOriginalLink(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expectedOriginalLink, resp.OriginalLink)

	testShortLink = "nonexistent"

	mockService.On("GetOriginalLink", testShortLink).Return("", errors.New("not found"))

	_, err = client.GetOriginalLink(context.Background(), &pb.GetOriginalLinkRequest{ShortLink: testShortLink})
	require.Error(t, err)

	testShortLink = ""

	mockService.On("GetOriginalLink", testShortLink).Return("", errors.New("not found"))

	_, err = client.GetOriginalLink(context.Background(), &pb.GetOriginalLinkRequest{ShortLink: testShortLink})
	require.Error(t, err)

	mockService.AssertExpectations(t)
}
