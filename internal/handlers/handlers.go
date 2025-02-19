package handlers

import (
	"context"
	"shortened_links_service_on_grpc/internal/services"
	pb "shortened_links_service_on_grpc/proto"
)

type ShortenerHandler struct {
	service services.ShortenerServiceInterface
	pb.UnimplementedShortenerServiceServer
}

func RegisterShortenerHandler(service services.ShortenerServiceInterface) *ShortenerHandler {
	return &ShortenerHandler{service: service}
}

// GetShortLink возвращает и сохраняет созданную сокращёную ссылку
func (s *ShortenerHandler) GetShortLink(ctx context.Context, req *pb.GetShortLinkRequest) (*pb.GetShortLinkResponse, error) {
	originalLink := req.OriginalLink
	shortLink, err := s.service.GetShortLink(originalLink)

	return &pb.GetShortLinkResponse{ShortLink: shortLink}, err
}

// GetOriginalLink возвращает оригинальную ссылку
func (s *ShortenerHandler) GetOriginalLink(ctx context.Context, req *pb.GetOriginalLinkRequest) (*pb.GetOriginalLinkResponse, error) {
	shortLink := req.ShortLink
	originalLink, err := s.service.GetOriginalLink(shortLink)

	return &pb.GetOriginalLinkResponse{OriginalLink: originalLink}, err
}
