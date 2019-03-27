package shortensvc

import (
	"errors"
	"strconv"
	"strings"

	"github.com/alextanhongpin/go-base62"
)

type (
	GetShortenedURLUseCase    func(shortURL string) (*ShortURL, error)
	GetShortenedURLRepository interface {
		WithID(shortURL string) (*ShortURL, error)
	}
)

func GetShortenedURL(urls GetShortenedURLRepository) GetShortenedURLUseCase {
	return func(shortURL string) (*ShortURL, error) {
		if len(strings.TrimSpace(shortURL)) == 0 {
			return nil, errors.New("url is required")
		}
		id := base62.Decode(shortURL)
		u, err := urls.WithID(strconv.FormatUint(id, 10))
		if err != nil {
			return nil, err
		}
		return u, nil
	}
}
