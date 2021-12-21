package test

import (
	"lzero/internal/data"
	"lzero/internal/utils"
	"strconv"
	"sync"
	"testing"
)

func TestCheck(t *testing.T) {
	o := data.ReceivedOrder{}
	d := data.DeliveryInfo{}
	p := data.Payment{}
	i := data.Item{}

	if o.CheckForMissingFields() == nil {
		t.FailNow()
	}
	if d.CheckForMissingFields() == nil {
		t.FailNow()
	}
	if p.CheckForMissingFields() == nil {
		t.FailNow()
	}
	if i.CheckForMissingFields() == nil {
		t.FailNow()
	}
}

func TestCacheAddMethod(t *testing.T) {
	var wg sync.WaitGroup

	c := *data.NewCache()

	var orderUID string
	for x := 0; x < 1000000; x++ {
		orderUID = strconv.Itoa(x)

		wg.Add(1)
		go func() {
			c.AddOrder(data.ReceivedOrder{OrderUID: &orderUID})
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestCacheGetMethod(t *testing.T) {
	var wg sync.WaitGroup

	c := *data.NewCache()

	var orderUID string
	for x := 0; x < 1000000; x++ {
		orderUID = strconv.Itoa(x)
		c.AddOrder(data.ReceivedOrder{OrderUID: &orderUID})
	}

	for x := 0; x < 1000000; x++ {
		orderUID = strconv.Itoa(x)
		wg.Add(1)
		go func() {
			_, ok := c.GetOrder(orderUID)
			if !ok {
				t.Fail()
			}
			wg.Done()
		}()
	}
}

func TestMarshal(t *testing.T) {
	l := utils.NewLogger()

	l.InfoLog.Println("Get test data")
	jsonEx, err := utils.GetData("./examples/model.json")
	if err != nil {
		t.FailNow()
	}

	resByte, _ := jsonEx.Marshal()
	t.Log(string(resByte))
}
