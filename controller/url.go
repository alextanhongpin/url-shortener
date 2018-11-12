package controller

import (
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/url-shortener/service"
	"github.com/julienschmidt/httprouter"
)

// URL represents the url controller and the corresponding endpoints.
type URL struct {
	service service.Shortener
}

// NewURL returns a new URL Controller.
func NewURL(svc service.Shortener) *URL {
	return &URL{
		service: svc,
	}
}

// Setup initializes the endpoints with the controller on the given router.
func (u *URL) Setup(r *httprouter.Router) {
	r.GET("/v1/urls/:id", u.GetURL)
	r.POST("/v1/urls", u.PostURL)
}

// GetURL represents the GET endpoint for url by id.
func (u *URL) GetURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	request := ps.ByName("id")
	response, err := u.service.Retrieve(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, response, http.StatusFound)
}

// PostURL represents the POST endpoint to create a new short url.
func (u *URL) PostURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL, err := u.service.Shorten(request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := ShortenURLResponse{
		URL: shortURL,
	}
	json.NewEncoder(w).Encode(response)
}

type ShortenURLRequest struct {
	URL string `json:"url,omitempty"`
}

type ShortenURLResponse struct {
	URL string `json:"url,omitempty"`
}
