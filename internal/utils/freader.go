package utils

import (
	"encoding/json"
	"io/ioutil"
	"lzero/internal/data"
)

func GetData(name string) (data.ReceivedOrder, error) {
	l := NewLogger()

	order := new(data.ReceivedOrder)

	l.InfoLog.Printf("Read json file: %s", name)
	jsonFile, err := ioutil.ReadFile(name)
	if err != nil {
		return data.ReceivedOrder{}, err
	}

	l.InfoLog.Println("Unmarshal json file")
	err = json.Unmarshal(jsonFile, order)
	if err != nil {
		return data.ReceivedOrder{}, err
	}
	
	err = order.CheckForMissingFields()
	if err != nil {
		return data.ReceivedOrder{}, err
	}

	return *order, err
}