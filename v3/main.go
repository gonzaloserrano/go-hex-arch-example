package main

import (
	"encoding/json"
	"net/http"

	"github.com/gonzaloserrano/go-hex-arch-example/v3/app"
	"github.com/gonzaloserrano/go-hex-arch-example/v3/infra"
	dialect "upper.io/db.v3/postgresql"
)

// Note: error handling omited for presentation purposes
func main() {
	db, _ := dialect.Open(dialect.ConnectionURL{
		Database: "counter",
		User:     "counter",
		Host:     "localhost:54321",
	})

	repo := infra.CounterPostgreSQLRepository{DB: db}

	http.HandleFunc("/counter/add", newAddHandler(repo))
	http.HandleFunc("/counter/get", newGetHandler(repo))

	http.ListenAndServe(":8080", nil)
}

func newAddHandler(repo app.CounterRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var httpCounter app.Counter
		json.NewDecoder(r.Body).Decode(&httpCounter)

		counter := repo.FindByID(httpCounter.ID)

		counter.ID = httpCounter.ID
		counter.Value += httpCounter.Value

		repo.Upsert(counter)
	}
}

func newGetHandler(repo app.CounterRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var c app.Counter
		json.NewDecoder(r.Body).Decode(&c)

		c = repo.FindByID(c.ID)

		json.NewEncoder(w).Encode(c)
	}
}
