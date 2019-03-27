package shortensvc

import (
	"errors"
	"net/url"
	"strings"

	"github.com/alextanhongpin/go-base62"
)

type (
	ShortenURLUseCase    func(longURL string) (string, error)
	ShortenURLRepository interface {
		Create(ShortURL) (int64, error)
	}
)

func ShortenURL(urls ShortenURLRepository) ShortenURLUseCase {
	return func(longURL string) (string, error) {
		if len(strings.TrimSpace(longURL)) == 0 {
			return "", errors.New("url is required")
		}
		u, err := url.Parse(longURL)
		if err != nil {
			return "", err
		}
		id, err := urls.Create(ShortURL{LongURL: u.String()})
		if err != nil {
			return "", err
		}
		shortID := base62.Encode(uint64(id))
		return shortID, nil
	}
}
