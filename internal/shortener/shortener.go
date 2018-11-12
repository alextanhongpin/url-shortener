package shortener

import (
	"net/url"

	"github.com/alextanhongpin/base62"
	"github.com/alextanhongpin/url-shortener/repository"
)

// Shortener represents the url shortener service.
type Shortener struct {
	repo        repository.URL
	formatterFn func(string) string
}

// NewService returns a new Shortener service.
func NewService(repo repository.URL, formatterFn func(string) string) *Shortener {
	return &Shortener{
		repo:        repo,
		formatterFn: formatterFn,
	}
}

// Shorten takes a long url and returns a shortened url.
func (s *Shortener) Shorten(longURL string) (string, error) {
	// Validate the url.
	u, err := url.Parse(longURL)
	if err != nil {
		return "", err
	}
	id, err := s.repo.Insert(u.String())
	if err != nil {
		return "", err
	}
	shortID := base62.Encode(id)
	return s.formatterFn(shortID), nil
}

// Retrieve takes a short url and find the corresponding long url.
func (s *Shortener) Retrieve(shortURL string) (string, error) {
	id := base62.Decode(shortURL)
	u, err := s.repo.Get(id)
	if err != nil {
		return "", err
	}
	return u.URL, nil
}
