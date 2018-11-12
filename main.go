package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alextanhongpin/url-shortener/controller"
	"github.com/alextanhongpin/url-shortener/internal/shortener"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

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

	router := httprouter.New()

	// Initialize the Shortener service dependencies.
	{
		repo := shortener.NewRepository(db)
		defer repo.Close()

		svc := shortener.NewService(repo, cname)
		ctl := controller.NewURL(svc)
		ctl.Setup(router)
	}

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
