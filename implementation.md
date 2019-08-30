## URL Shortener

Different implementation using md5 hash instead. This allows user to customize the short url, since it is not dependent on the sequential id that is generated from the database.

```go
package main

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

var ErrAlreadyExists = errors.New("already exists")

type URLShortener interface {
	Shorten(longURL string) (code string)
}

type Shortener struct{}

func NewShortener() *Shortener {
	return &Shortener{}
}

func (s *Shortener) Shorten(longURL string) (code string) {
	h := md5.New()
	h.Write([]byte(longURL))
	msg := h.Sum(nil)
	code = base64.URLEncoding.EncodeToString(msg[:6])
	return
}

type ShortenerService interface {
	Get(code string) (longURL string, err error)
	Put(code, longURL string, expireAt time.Time) (string, error)
	CheckExists(code string) bool
}

type Repository interface {
	Get(code string) (longURL string, err error)
	Put(code, longURL string, expireAt time.Time) error
	CheckExists(code string) bool
}

type ShortenerServiceImpl struct {
	baseURL    string
	shortener  URLShortener
	repository Repository
}

func NewShortenerService(baseURL string, shortener URLShortener, repository Repository) *ShortenerServiceImpl {
	return &ShortenerServiceImpl{
		baseURL:    baseURL,
		shortener:  shortener,
		repository: repository,
	}
}

func (s *ShortenerServiceImpl) Get(code string) (string, error) {
	return s.repository.Get(code)
}

func (s *ShortenerServiceImpl) Put(code, longURL string, expireAt time.Time) (string, error) {
	// User customize the shortURL.
	if code != "" {
		if exists := s.repository.CheckExists(code); exists {
			return "", ErrAlreadyExists
		}
	} else {
		// Generate the short url if none is provided.
		code = s.shortener.Shorten(longURL)
	}
	err := s.repository.Put(code, longURL, expireAt)
	for err == ErrAlreadyExists {
		// Keep shortening until the URL is found.
		code = s.shortener.Shorten(code)
		err = s.repository.Put(code, longURL, expireAt)
	}
	return code, err
}

func (s *ShortenerServiceImpl) CheckExists(code string) bool {
	return s.repository.CheckExists(code)
}

func main() {
	shortener := NewShortener()
	shortURL := shortener.Shorten("http://www.google.com")
	fmt.Println(shortURL)
	shortURL = shortener.Shorten(shortURL)
	fmt.Println(shortURL)
}
```
