package handlers_test

import (
	"context"
	"errors"
	"shortened_links_service_on_grpc/internal/handlers"
	pb "shortened_links_service_on_grpc/proto"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetShortLinkUnit проверяет работу метода GetShortLink хендлера с использованием мок-сервиса
func TestGetShortLinkUnit(t *testing.T) {
	mockService := new(handlers.MockShortenerService)
	handler := handlers.RegisterShortenerHandler(mockService)

	testOriginalLink := "https://example.com"
	expectedShortLink := "abcd123456"

	mockService.On("GetShortLink", testOriginalLink).Return(expectedShortLink, nil)

	req := &pb.GetShortLinkRequest{OriginalLink: testOriginalLink}
	resp, err := handler.GetShortLink(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expectedShortLink, resp.ShortLink)

	testOriginalLink = ""

	mockService.On("GetShortLink", testOriginalLink).Return("", errors.New("invalid URL format"))
	req = &pb.GetShortLinkRequest{OriginalLink: testOriginalLink}
	_, err = handler.GetShortLink(context.Background(), req)
	require.Error(t, err)

	testOriginalLink = "https:::|/example.com"

	mockService.On("GetShortLink", testOriginalLink).Return("", errors.New("invalid URL format"))
	req = &pb.GetShortLinkRequest{OriginalLink: testOriginalLink}
	_, err = handler.GetShortLink(context.Background(), req)
	require.Error(t, err)

	mockService.AssertExpectations(t)
}

// TestGetOriginalLinkUnit проверяет работу метода GetOriginalLink хендлера с использованием мок-сервиса
func TestGetOriginalLinkUnit(t *testing.T) {
	mockService := new(handlers.MockShortenerService)
	handler := handlers.RegisterShortenerHandler(mockService)

	testShortLink := "abcd123456"
	expectedOriginalLink := "https://example.com"

	mockService.On("GetOriginalLink", testShortLink).Return(expectedOriginalLink, nil)

	req := &pb.GetOriginalLinkRequest{ShortLink: testShortLink}
	resp, err := handler.GetOriginalLink(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expectedOriginalLink, resp.OriginalLink)

	testShortLink = "nonexistent"

	mockService.On("GetOriginalLink", testShortLink).Return("", errors.New("not found"))

	_, err = handler.GetOriginalLink(context.Background(), &pb.GetOriginalLinkRequest{ShortLink: testShortLink})
	require.Error(t, err)

	testShortLink = ""

	mockService.On("GetOriginalLink", testShortLink).Return("", errors.New("not found"))

	_, err = handler.GetOriginalLink(context.Background(), &pb.GetOriginalLinkRequest{ShortLink: testShortLink})
	require.Error(t, err)

	mockService.AssertExpectations(t)
}
