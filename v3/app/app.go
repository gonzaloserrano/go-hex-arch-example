package app

type Counter struct {
	ID    string `db:"id" json:"id"`
	Value int    `db:"value" json:"value"`
}

type CounterRepository interface {
	FindByID(ID string) Counter
	Upsert(c Counter)
}
