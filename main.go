package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/url-shortener/shortensvc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var shutdowns grace.Shutdowns
	// cname  = os.Getenv("CNAME")
	port := os.Getenv("PORT")

	db := initDB()
	defer db.Close()

	router := httprouter.New()

	// Initialize the Shortener service dependencies.
	{
		repo := shortensvc.NewRepository(db)
		svc := shortensvc.NewService(repo)
		ctl := shortensvc.NewController(svc)

		router.GET("/v1/urls/:id", ctl.GetShortURLByID)
		router.POST("/v1/urls", ctl.PostShortURLs)
	}

	shutdowns.Append(grace.New(router, port))
	<-grace.Signal()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdowns.Close(ctx)
}

func initDB() *sql.DB {
	var (
		user = os.Getenv("DB_USER")
		pass = os.Getenv("DB_PASS")
		name = os.Getenv("DB_NAME")
	)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", user, pass, name))
	if err != nil {
		log.Fatal(err)
	}
	// Set the number of open and idle connection to a maximum total of 3.
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	for i := 0; i < 3; i++ {
		if err := db.Ping(); err != nil {
			log.Println("retrying db connection in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	return db
}
