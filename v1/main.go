package main

import (
	"encoding/json"
	"net/http"

	upper "upper.io/db.v3"
	dialect "upper.io/db.v3/postgresql"
)

// Note: error handling omited for presentation purposes
func main() {
	db, _ := dialect.Open(dialect.ConnectionURL{
		Database: "counter",
		User:     "counter",
		Host:     "localhost:54321",
	})

	type counter struct {
		ID    string `db:"id" json:"id"`
		Value int    `db:"value" json:"value"`
	}

	http.HandleFunc("/counter/add", func(w http.ResponseWriter, r *http.Request) {
		var httpCounter counter
		json.NewDecoder(r.Body).Decode(&httpCounter)

		var dbCounter counter
		db.Collection("counter").Find(upper.Cond{"id": httpCounter.ID}).One(&dbCounter)

		dbCounter.ID = httpCounter.ID
		dbCounter.Value += httpCounter.Value

		upsertFunc := func(q string) string { return q + " ON CONFLICT(id) DO UPDATE SET value = EXCLUDED.value" }
		_, err := db.InsertInto("counter").Values(dbCounter).Amend(upsertFunc).Exec()
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/counter/get", func(w http.ResponseWriter, r *http.Request) {
		var c counter
		json.NewDecoder(r.Body).Decode(&c)
		db.Collection("counter").Find(upper.Cond{"id": c.ID}).One(&c)
		json.NewEncoder(w).Encode(c)
	})

	http.ListenAndServe(":8080", nil)
}
