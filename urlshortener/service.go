package urlshortener

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/alextanhongpin/go-base62"
)

type (
	repository interface {
		WithID(shortURL string) (*ShortURL, error)
		Create(ShortURL) (int64, error)
	}
	Service struct {
		urls repository
	}
)

func NewService(urls repository) *Service {
	return &Service{urls}
}

func (s *Service) ShortenURL(longURL string) (string, error) {
	if isNilString(longURL) {
		return "", errors.New("url is required")
	}
	u, err := url.Parse(longURL)
	if err != nil {
		return "", err
	}
	id, err := s.urls.Create(ShortURL{LongURL: u.String()})
	if err != nil {
		return "", err
	}
	shortID := base62.Encode(uint64(id))
	return shortID, nil
}

func (s *Service) GetShortURL(shortURL string) (*ShortURL, error) {
	if isNilString(shortURL) {
		return nil, errors.New("url is required")
	}
	id := base62.Decode(shortURL)
	u, err := s.urls.WithID(strconv.FormatUint(id, 10))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func isNilString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
