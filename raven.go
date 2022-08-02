package main

import (
	ravendb "github.com/ravendb/ravendb-go-client"
)

var store, _ = getDocumentStore("NCAAF")

func getDocumentStore(databaseName string) (*ravendb.DocumentStore, error) {
	serverNodes := []string{"http://dell-4290.local:5050"}
	store := ravendb.NewDocumentStore(serverNodes, databaseName)
	if err := store.Initialize(); err != nil {
		return nil, err
	}
	return store, nil
}

func openSession() *ravendb.DocumentSession {

	var session, err = store.OpenSession("")
	if err != nil {
		panic(err)
	}

	return session
}
