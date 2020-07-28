package infra

import (
	"github.com/gonzaloserrano/go-hex-arch-example/v3/app"
	upper "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type CounterPostgreSQLRepository struct {
	DB sqlbuilder.Database
}

var _ app.CounterRepository = CounterPostgreSQLRepository{}

func (r CounterPostgreSQLRepository) FindByID(ID string) app.Counter {
	var c app.Counter
	r.DB.Collection("counter").Find(upper.Cond{"id": c.ID}).One(&c)
	return c
}

func (r CounterPostgreSQLRepository) Upsert(c app.Counter) {
	r.DB.InsertInto("counter").Values(c).Amend(onConflictUpdate).Exec()
}

func onConflictUpdate(q string) string {
	return q + " ON CONFLICT(id) DO UPDATE SET value = EXCLUDED.value"
}
