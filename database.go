package main

import (
	"fmt"

	f "github.com/fauna/faunadb-go/v4/faunadb"
	"github.com/google/uuid"
)

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

	*key = secret

	fmt.Printf("Secret: %s\n", *key)

	return nil
}
