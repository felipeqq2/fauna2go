package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	f "github.com/fauna/faunadb-go/v4/faunadb"
	"github.com/google/uuid"
)

func init() {
	time.Sleep(15 * time.Second)
}

func main() {
	client := f.NewFaunaClient("secret", f.Endpoint("http://localhost:8443"))

	var secret string

	if err := createDatabase(client, &secret); err != nil {
		fatal(err)
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprint(rw, secret)
		} else {
			if err := createDatabase(client, &secret); err != nil {
				fatal(err)
			}
			fmt.Fprint(rw, secret)
		}
	})

	fatal(http.ListenAndServe(":1000", nil))
}

func createDatabase(client *f.FaunaClient, key *string) error {
	fmt.Println("\nCreating new database")

	db := uuid.NewString()

	dbRes, err := client.Query(f.CreateDatabase(f.Obj{"name": db}))
	if err != nil {
		return err
	}

	fmt.Printf("Database: %s\n", db)

	var ref f.RefV
	if err := dbRes.At(f.ObjKey("ref")).Get(&ref); err != nil {
		return err
	}

	keyRes, err := client.Query(f.CreateKey(f.Obj{"role": "admin", "database": ref}))
	if err != nil {
		return err
	}

	var secret string
	if err := keyRes.At(f.ObjKey("secret")).Get(&secret); err != nil {
		return err
	}

	fmt.Printf("Secret: %s\n", secret)

	*key = secret

	return nil
}

func fatal(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	})

	log.Fatal(http.ListenAndServe(":1000", nil))
}
