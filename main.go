package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alextanhongpin/base62"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type ShortenUrlRequest struct {
	URL string `json:"url,omitempty"`
}

type ShortenUrlResponse struct {
	URL string `json:"url,omitempty"`
}

type URLEntity struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	URLCRC    uint64    `json:"url_crc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func main() {
	var (
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASS")
		dbName = os.Getenv("DB_NAME")
		cname  = os.Getenv("CNAME")
		port   = os.Getenv("PORT")
	)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPass, dbName))
	if err != nil {
		log.Fatal(err)
	}
	// Set the number of open and idle connection to a maximum total of 3.
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	stmtCreate, err := db.Prepare("INSERT INTO url (url, url_crc) VALUES (?, CRC32(?))")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtCreate.Close()

	stmtGet, err := db.Prepare("SELECT (url) FROM url WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtGet.Close()

	createURL := func(longurl string) (int64, error) {

		res, err := stmtCreate.Exec(longurl, longurl)
		if err != nil {
			return -1, err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return -1, err
		}
		if rows != 1 {
			return -1, errors.New("insert fail")
		}
		return res.LastInsertId()
	}

	router := httprouter.New()
	router.GET("/v1/urls/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		enc := ps.ByName("id")
		id := base62.Decode(enc)
		var response URLEntity
		if err := stmtGet.QueryRow(id).Scan(&response.URL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, response.URL, http.StatusFound)
	})
	router.POST("/v1/urls", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		var request ShortenUrlRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Check if the url is valid.
		u, err := url.Parse(request.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: Check if the url exists using bloom filter.
		id, err := createURL(u.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortened := base62.Encode(uint64(id))
		response := ShortenUrlResponse{
			URL: fmt.Sprintf("https://%s/%s", cname, shortened),
		}
		json.NewEncoder(w).Encode(response)
	})

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		log.Printf("listening to port *%s. press ctrl + c to cancel.\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("gracefully shutting down server")
}
