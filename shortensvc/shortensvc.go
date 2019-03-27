package shortensvc

type Service struct {
	ShortenURL      ShortenURLUseCase
	GetShortenedURL GetShortenedURLUseCase
}

func NewService(repo Repository) *Service {
	return &Service{
		ShortenURL:      ShortenURL(repo),
		GetShortenedURL: GetShortenedURL(repo),
	}
}
