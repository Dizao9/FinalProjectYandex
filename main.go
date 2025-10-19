package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v4"
	"main.go/pkg/api"
	"main.go/pkg/database"
	_ "modernc.org/sqlite"
)

func getPort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return "7540"
}

func getDBFile() string {
	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		return dbFile
	}
	return "./pkg/database/scheduler.db"
}

func main() {
	dbFile := getDBFile()
	if err := database.Init(dbFile); err != nil {
		log.Fatalf("Ошибка при открытии базы данных: %v", err)
	}
	defer database.DB.Close()

	r := chi.NewRouter()
	port := getPort()
	api.Init(r)

	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fileServer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	log.Printf("Server started and listening on http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server startup error: %v", err)
	}

}
