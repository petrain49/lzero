package db

import (
	"database/sql"
	"fmt"
	"lzero/internal/data"
	"lzero/internal/utils"
	"sync"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "wb_petr"
	DB_PASSWORD = "test"
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_NAME     = "test"
)

type Database struct {
	DB *sql.DB
}

func OpenDB() (*Database, error) {
	l := utils.NewLogger()

	pgConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	l.InfoLog.Printf("Open the database: %s", pgConn)
	db, err := sql.Open("postgres", pgConn)
	if err != nil {
		l.ErrorLog.Printf("Fail to open DB: %s", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		l.ErrorLog.Printf("DB connection is dead: %s", err)
		return nil, err
	}

	return &Database{DB: db}, err
}

func (db *Database) UploadOrder(wg *sync.WaitGroup, order data.ReceivedOrder) error {
	l := utils.NewLogger()

	orderStatement := "INSERT INTO orders VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	deliveryStatement := "INSERT INTO deliveries VALUES($1, $2, $3, $4, $5, $6, $7, $8)"
	paymentStatement := "INSERT INTO payments VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	itemStatement := "INSERT INTO items VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	l.InfoLog.Printf("Start of the transaction")
	tx, err := db.DB.Begin()
	if err != nil {
		l.ErrorLog.Printf("Fail to start the transaction: %s", err)
		tx.Rollback()
		return err
	}

	insertItems, err := tx.Prepare(itemStatement)
	if err != nil {
		l.ErrorLog.Printf("Fail to prepare statement %s: %s", itemStatement, err)
		tx.Rollback()
		return err
	}
	defer insertItems.Close()

	_, err = tx.Exec(orderStatement,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OOFShard,
	)
	if err != nil {
		l.ErrorLog.Printf("Fail to execute statement %s: %s", orderStatement, err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(deliveryStatement,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.ZIP,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		l.ErrorLog.Printf("Fail to execute statement %s: %s", deliveryStatement, err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(paymentStatement,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		l.ErrorLog.Printf("Fail to execute statement %s: %s", paymentStatement, err)
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = insertItems.Exec(
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			l.ErrorLog.Printf("Fail to execute statement %s: %s", itemStatement, err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	l.InfoLog.Printf("End of the transaction: %s, err: %s", *order.OrderUID, err)
	wg.Done()
	return err
}

func (db *Database) Recovery(orderSet *data.Cache) error {
	l := utils.NewLogger()

	l.InfoLog.Printf("Start the cache recovery from DB")

	orderStatement := "SELECT * FROM orders"
	deliveryStatement := "SELECT name, phone, zip, city, address, region, email FROM deliveries where order_uid = $1"
	paymentStatement := "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payments WHERE transaction = $1"
	itemStatement := "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE track_number = $1"

	//prepared statements
	getFromDeliveries, err := db.DB.Prepare(deliveryStatement)
	if err != nil {
		l.ErrorLog.Printf("Fail to prepare statement %s: %s", deliveryStatement, err)
		return err
	}
	defer getFromDeliveries.Close()

	getFromPayments, err := db.DB.Prepare(paymentStatement)
	if err != nil {
		l.ErrorLog.Printf("Fail to prepare statement %s: %s", paymentStatement, err)
		return err
	}
	defer getFromPayments.Close()

	getFromItem, err := db.DB.Prepare(itemStatement)
	if err != nil {
		l.ErrorLog.Printf("Fail to prepare statement %s: %s", itemStatement, err)
		return err
	}
	defer getFromItem.Close()

	// get every order, every delivery info and payment for every order, all items for every order
	orderRows, err := db.DB.Query(orderStatement)
	if err != nil {
		l.ErrorLog.Printf("Fail to execute statement %s: %s", orderStatement, err)
		return err
	}
	defer orderRows.Close()

	for orderRows.Next() {
		order := new(data.ReceivedOrder)

		err = orderRows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OOFShard,
		)
		if err != nil {
			l.ErrorLog.Printf("Fail to scan a row from orders table: %s", err)
			return err
		}

		err = getFromDeliveries.QueryRow(order.OrderUID).Scan(
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.ZIP,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
		)
		if err != nil {
			l.ErrorLog.Printf("Fail to scan a row from deliveries table: %s", err)
			return err
		}

		err = getFromPayments.QueryRow(order.OrderUID).Scan(
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDt,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			l.ErrorLog.Printf("Fail to scan a row from payments table: %s", err)
			return err
		}

		itemRows, err := getFromItem.Query(order.TrackNumber)
		if err != nil {
			l.ErrorLog.Printf("Fail to execute prepared statement %s: %s", itemStatement, err)
			return err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			item := new(data.Item)

			err = itemRows.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.RID,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				l.ErrorLog.Printf("Fail to scan a row from items table: %s", err)
				return err
			}

			order.Items = append(order.Items, *item)
		}

		orderSet.AddOrder(*order)
	}

	l.InfoLog.Printf("End of the cache recovery")
	return err
}
