package data

import (
	"encoding/json"
	"errors"
)

type ReceivedOrder struct {
	OrderUID          *string      `json:"order_uid"`
	TrackNumber       *string      `json:"track_number"`
	Entry             *string      `json:"entry"`
	Delivery          DeliveryInfo `json:"delivery"`
	Payment           Payment      `json:"payment"`
	Items             []Item       `json:"items"`
	Locale            *string      `json:"locale"`
	InternalSignature *string      `json:"internal_signature"`
	CustomerID        *string      `json:"customer_id"`
	DeliveryService   *string      `json:"delivery_service"`
	ShardKey          *string      `json:"shardkey"`
	SmID              *int         `json:"sm_id"`
	DateCreated       *string      `json:"date_created"`
	OOFShard          *string      `json:"oof_shard"`
}

type DeliveryInfo struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	ZIP     *string `json:"zip"`
	City    *string `json:"city"`
	Address *string `json:"address"`
	Region  *string `json:"region"`
	Email   *string `json:"email"`
}

type Payment struct {
	Transaction  *string `json:"transaction"`
	RequestID    *string `json:"request_id"`
	Currency     *string `json:"currency"`
	Provider     *string `json:"provider"`
	Amount       *int    `json:"amount"`
	PaymentDt    *int    `json:"payment_dt"`
	Bank         *string `json:"bank"`
	DeliveryCost *int    `json:"delivery_cost"`
	GoodsTotal   *int    `json:"goods_total"`
	CustomFee    *int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      *int    `json:"chrt_id"`
	TrackNumber *string `json:"track_number"`
	Price       *int    `json:"price"`
	RID         *string `json:"rid"`
	Name        *string `json:"name"`
	Sale        *int    `json:"sale"`
	Size        *string `json:"size"`
	TotalPrice  *int    `json:"total_price"`
	NmID        *int    `json:"nm_id"`
	Brand       *string `json:"brand"`
	Status      *int `json:"status"`
}

func NewOrder(byteOrder []byte) (ReceivedOrder, error) {
	res := new(ReceivedOrder)

	err := json.Unmarshal(byteOrder, res)
	
	return *res, err
}

func (r *ReceivedOrder) Marshal() ([]byte, error) {
	res, err := json.Marshal(r)
	return res, err
}

func (r *ReceivedOrder) CheckForMissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case r.OrderUID == nil:
		return
	case r.TrackNumber == nil:
		return
	case r.Entry == nil:
		return
	case r.Locale == nil:
		return
	case r.InternalSignature == nil:
		return
	case r.CustomerID == nil:
		return
	case r.DeliveryService == nil:
		return
	case r.ShardKey == nil:
		return
	case r.SmID == nil:
		return
	case r.DateCreated == nil:
		return
	case r.OOFShard == nil:
		return
	case r.Delivery.CheckForMissingFields() != nil:
		return
	case r.Payment.CheckForMissingFields() != nil:
		return
	}

	for _, item := range r.Items {
		if item.CheckForMissingFields() != nil {
			return
		}
	}

	return nil
}

func (d *DeliveryInfo) CheckForMissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case d.Name == nil:
		return
	case d.Phone == nil:
		return
	case d.ZIP == nil:
		return
	case d.City == nil:
		return
	case d.Address == nil:
		return
	case d.Region == nil:
		return
	case d.Email == nil:
		return
	default:
		return nil
	}
}

func (p *Payment) CheckForMissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case p.Transaction == nil:
		return
	case p.RequestID == nil:
		return
	case p.Currency == nil:
		return
	case p.Provider == nil:
		return
	case p.Amount == nil:
		return
	case p.PaymentDt == nil:
		return
	case p.Bank == nil:
		return
	case p.DeliveryCost == nil:
		return
	case p.GoodsTotal == nil:
		return
	case p.CustomFee == nil:
		return
	default:
		return nil
	}
}

func (item *Item) CheckForMissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case item.ChrtID == nil:
		return
	case item.TrackNumber == nil:
		return
	case item.Price == nil:
		return
	case item.RID == nil:
		return
	case item.Name == nil:
		return
	case item.Sale == nil:
		return
	case item.Size == nil:
		return
	case item.TotalPrice == nil:
		return
	case item.NmID == nil:
		return
	case item.Brand == nil:
		return
	case item.Status == nil:
		return
	default:
		return nil
	}
}
