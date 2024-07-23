package main

import (
	"context"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"testing"
)

func TestLambdaHandler(t *testing.T) {
	testEvent := Event{
		Query: `CREATE TABLE IF NOT EXISTS contacts (
	contact_id INTEGER PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	email TEXT NOT NULL,
	phone TEXT NOT NULL
);`,
		Action: "INSERT",
	}
	_, err := LambdaHandler(context.TODO(), &testEvent)
	if err != nil {
		log.Fatal(err)
	}

}
