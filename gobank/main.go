package main

import "log"

func main() {
	listenAddr := "127.0.0.1:3000"
	dbConStr := "postgres://postgres:gobank@localhost:5432/postgres?sslmode=disable"
	store, err := NewPostgresDB(dbConStr)
	if err != nil {
		log.Fatalln("Unable to open DB connection", err)
		return
	}

	initErr := store.Init()
	if initErr != nil {
		log.Fatalln("Unable to initialize DB", initErr)
		return
	}

	apiServer := NewApiServer(listenAddr, store)

	apiServer.Run()
}
