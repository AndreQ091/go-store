package main

import (
	"go-store/core"
	"go-store/frontend"
	"go-store/transact"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	tl, err := transact.NewTransactionLogger("postgres")
	if err != nil {
		log.Fatal(err)
	}

	store := core.NewKeyValueStore(tl)
	store.Restore()

	fe, err := frontend.NewFrontend("rest")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(fe.Start(store))
}
