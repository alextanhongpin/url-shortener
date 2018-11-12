package service

// Shortener provides method to shorten the url and retrieve back the long url
// from the shortened url.
type Shortener interface {
	// Shortens a long url.
	Shorten(longURL string) (string, error)

	// Retrieve the long url back from the shortened url.
	Retrieve(shortURL string) (string, error)
}
