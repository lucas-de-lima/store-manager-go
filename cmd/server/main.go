package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lucas-de-lima/store-manager-go/handler"
	"github.com/lucas-de-lima/store-manager-go/repository"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/StoreManager"
	}

	db, err := repository.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	saleRepo := repository.NewSaleRepository(db)

	router := handler.NewRouter(productRepo, saleRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}