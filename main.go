package main

import (
	"fmt"
	"io/ioutil"
	"lzero/internal/data"
	"lzero/internal/db"
	"lzero/internal/nats"
	"lzero/internal/serv"
	"lzero/internal/utils"
	"net/http"
	"sync"
)

func main() {
	l := utils.NewLogger()
	var wg sync.WaitGroup

	cache := data.NewCache()

	output := make(chan []byte)

	db, err := db.OpenDB() // Open db
	if err != nil {
		l.ErrorLog.Fatal(err)
	}
	defer db.DB.Close()

	db.Recovery(cache) // Get orders from db

	natsConn, err := nats.NewConnection() // Connect to NATS server
	if err != nil {
		l.ErrorLog.Fatal(err)
	}
	defer natsConn.STANConn.Close()

	subscription, err := natsConn.Subscribe(output) // Subscribe to NATS streaming
	if err != nil {
		l.ErrorLog.Fatal(err)
	}
	defer subscription.Close()

	Publisher(&natsConn)

	wg.Add(1)
	go func() {
		if !subscription.IsValid() {
			wg.Done()
			return
		}

		var order data.ReceivedOrder

		for byteOrder := range output {
			order, err = data.NewOrder(byteOrder)
			if err != nil {
				l.ErrorLog.Printf("Cant form order struct: %s\n", err)
			}

			err = order.CheckForMissingFields()
			if err != nil {
				l.ErrorLog.Printf("Broken order: %s\n", err)
			} else {
				go db.UploadOrder(order)
				go cache.AddOrder(order)
			}
		}
		wg.Done()
	}()

	server := serv.NewServer(cache)
	l.ErrorLog.Println(http.ListenAndServe(serv.HTTP_PORT, server)) // Run http server

	wg.Wait()
}

func Publisher(nc *nats.NATSConn) {
	l := utils.NewLogger()

	files, err := ioutil.ReadDir("./test/examples")
    if err != nil {
        return
    }

	for _, f := range files {
		fullName := fmt.Sprintf("./test/examples/%s", f.Name())
		byteOrder, err := ioutil.ReadFile(fullName)
		if err != nil {
			l.ErrorLog.Fatal(err)
		}
		nc.Publish(byteOrder)
	}
}
