package test

import (
	"lzero/internal/data"
	"lzero/internal/db"
	"lzero/internal/utils"
	"testing"
)

func TestConnection(t *testing.T) {
	db, err := db.OpenDB()
	if err != nil {
		t.Fail()
	}

	err = db.DB.Ping()
	if err != nil {
		t.Fail()
	}
}

func TestUploadOrder(t *testing.T) {
	l := utils.NewLogger()

	db, err := db.OpenDB()
	if err != nil {
		t.FailNow()
	}
	defer db.DB.Close()

	l.InfoLog.Println("Get test data")
	jsonEx, err := utils.GetData("./examples/model.json")
	if err != nil {
		t.FailNow()
	}

	l.InfoLog.Println("Upload test data to the DB")
	err = db.UploadOrder(jsonEx)	
	if err != nil {
		t.FailNow()
	}
}

func TestRecovery(t *testing.T) {
	l := utils.NewLogger()

	l.InfoLog.Println("Open the DB")
	db, err := db.OpenDB()
	if err != nil {
		t.FailNow()
	}
	defer db.DB.Close()

	cache := data.NewCache()
	l.InfoLog.Println("Recovery")
	err = db.Recovery(cache)
	if err != nil {
		t.Fail()
	}
	
	for _, f := range cache.GetOrders() {
		if err = f.CheckForMissingFields(); err != nil {
			t.Fail()
		}
		t.Log(*f.OrderUID)
	}
}


func TestCleanTables(t *testing.T) {
	l := utils.NewLogger()

	l.InfoLog.Println("Open the DB")
	db, err := db.OpenDB()
	if err != nil {
		t.FailNow()
	}
	defer db.DB.Close()

	_, err = db.DB.Exec("DELETE FROM deliveries")
	if err != nil {
		t.Fail()
	}
	_, err = db.DB.Exec("DELETE FROM payments")
	if err != nil {
		t.Fail()
	}
	_, err = db.DB.Exec("DELETE FROM items")
	if err != nil {
		t.Fail()
	}
	_, err = db.DB.Exec("DELETE FROM orders")
	if err != nil {
		t.Fail()
	}
}
