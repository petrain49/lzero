package nats

import (
	"lzero/internal/utils"

	"github.com/nats-io/stan.go"
)

const (
	NATS_CLUSTER = "test-cluster"
	NATS_SUBJECT = "orders"
	NATS_CLIENT  = "client"
)

type NATSConn struct {
	STANConn stan.Conn
}

func NewConnection() (NATSConn, error) {
	l := utils.NewLogger()

	l.InfoLog.Println("New connection to NATS streaming")

	nc := new(NATSConn)
	var err error

	l.InfoLog.Printf("Connect to nats streaming, cluster: %s, client: %s\n", NATS_CLUSTER, NATS_CLIENT)
	nc.STANConn, err = stan.Connect(NATS_CLUSTER, NATS_CLIENT)

	return *nc, err
}

func (nc *NATSConn) Publish(jsonName []byte) {
	nc.STANConn.Publish(NATS_SUBJECT, jsonName)
}

func (nc *NATSConn) Subscribe(output chan<- []byte) (stan.Subscription, error) {
	l := utils.NewLogger()

	l.InfoLog.Println("Subscribe to NATS streaming")

	sub, err := nc.STANConn.Subscribe(NATS_SUBJECT, func(msg *stan.Msg) {
		l.InfoLog.Println("Received message")
		output <- msg.Data

	}, stan.DeliverAllAvailable())
	if err != nil {
		l.ErrorLog.Printf("Failed subscription: %s\n", err)
	}

	return sub, err
}
