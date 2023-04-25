package ascendex

import (
	"strconv"
	"test/api"
)

func makeOrderFromArray(data []string) (*api.Order, error) {

	amount, err := strconv.ParseFloat(data[1], 64)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(data[0], 64)
	if err != nil {
		return nil, err
	}

	return &api.Order{
		Amount: amount,
		Price:  price,
	}, nil
}
