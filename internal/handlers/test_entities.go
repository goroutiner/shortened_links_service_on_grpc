package handlers

import "github.com/stretchr/testify/mock"

// MockShortenerService — мок для ShortenerService
type MockShortenerService struct {
	mock.Mock
}

// GetShortLink — заглушка для метода получения сокращённой ссылки
func (m *MockShortenerService) GetShortLink(originalLink string) (string, error) {
	args := m.Called(originalLink)
	return args.String(0), args.Error(1)
}

// GetOriginalLink — заглушка для метода получения оригинального URL
func (m *MockShortenerService) GetOriginalLink(shortLink string) (string, error) {
	args := m.Called(shortLink)
	return args.String(0), args.Error(1)
}
