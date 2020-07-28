package main

import (
	"encoding/json"
	"net/http"

	upper "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	dialect "upper.io/db.v3/postgresql"
)

// Note: error handling omited for presentation purposes
func main() {
	db, _ := dialect.Open(dialect.ConnectionURL{
		Database: "counter",
		User:     "counter",
		Host:     "localhost:54321",
	})

	http.HandleFunc("/counter/add", newAddHandler(db))
	http.HandleFunc("/counter/get", newGetHandler(db))

	http.ListenAndServe(":8080", nil)
}

type counter struct {
	ID    string `db:"id" json:"id"`
	Value int    `db:"value" json:"value"`
}

func newAddHandler(db sqlbuilder.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var httpCounter counter
		json.NewDecoder(r.Body).Decode(&httpCounter)

		var dbCounter counter
		db.Collection("counter").Find(upper.Cond{"id": httpCounter.ID}).One(&dbCounter)

		dbCounter.ID = httpCounter.ID
		dbCounter.Value += httpCounter.Value

		db.InsertInto("counter").Values(dbCounter).Amend(onConflictUpdate).Exec()
	}
}

func newGetHandler(db sqlbuilder.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var c counter
		json.NewDecoder(r.Body).Decode(&c)
		db.Collection("counter").Find(upper.Cond{"id": c.ID}).One(&c)
		json.NewEncoder(w).Encode(c)
	}
}

func onConflictUpdate(q string) string {
	return q + " ON CONFLICT(id) DO UPDATE SET value = EXCLUDED.value"
}
