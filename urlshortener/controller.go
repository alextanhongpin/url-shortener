package urlshortener

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Controller struct {
	// There's no need to pass an interface here, since we are not gonna
	// mock the Service implementation.
	service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{svc}
}

func (ctl *Controller) GetShortURLByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := ps.ByName("id")
	res, err := ctl.service.GetShortURL(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, res.LongURL, http.StatusFound)
}

func (ctl *Controller) PostShortURLs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type request struct {
		URL       string    `json:"long_url,omitempty"`
		ExpiresAt time.Time `json:"expires_at,omitempty"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := ctl.service.ShortenURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}
