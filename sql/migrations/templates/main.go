package main

import (
	"os"
	"text/template"
	"time"
)

func createMonthlyDates(start, end time.Time) []time.Time {
	t := []time.Time{}

	// Change day to the first of the month
	curr := start.AddDate(0, 0, -start.Day()+1)

	for curr.Before(end) {
		t = append(t, curr)
		curr = curr.AddDate(0, 1, 0) // + 1 month
	}

	return t
}

// Run with `make generate-sql` from the root
func main() {
	templateFile := "sql/migrations/templates/001_create_initial_table.up.tmpl.sql"

	// Create 2 years of partitions starting from the current month
	start := time.Now()
	end := start.AddDate(2, 0, 0)
	dates := createMonthlyDates(start, end)

	t := template.Must(template.ParseFiles(templateFile))
	t.Execute(os.Stdout, dates)
}
