package tests

import (
	"strconv"
	"test/api"
	"test/api/ascendex"
	"testing"
	"time"
)

func TestToBestOrderBook(t *testing.T) {
	cases := []struct {
		name string
		data ascendex.BboInputMessage
		exp  api.BestOrderBook
	}{
		{
			name: "case1",
			data: ascendex.BboInputMessage{
				"bbo",
				"BTC/USDT",
				struct {
					Ts  int64    `json:"ts"`
					Bid []string `json:"bid"`
					Ask []string `json:"ask"`
				}{Ts: 123333, Bid: []string{"5454.1", "0.1"}, Ask: []string{"5455", "0.2"}},
			},
			exp: api.BestOrderBook{
				Ask: api.Order{
					0.2,
					5455,
				},
				Bid: api.Order{
					0.1,
					5454.1,
				},
			},
		},
		{
			name: "case2",
			data: ascendex.BboInputMessage{
				"bbo",
				"BTC/USDT",
				struct {
					Ts  int64    `json:"ts"`
					Bid []string `json:"bid"`
					Ask []string `json:"ask"`
				}{Ts: 123333, Bid: []string{"9999999", "9"}, Ask: []string{"99.999", "0.992"}},
			},
			exp: api.BestOrderBook{
				Ask: api.Order{
					0.992,
					99.999,
				},
				Bid: api.Order{
					9,
					9999999,
				},
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			convertedOrderBook, err := tCase.data.ToBestOrderBook()
			if err != nil {
				t.Error(err)
			}

			if convertedOrderBook.Ask != tCase.exp.Ask || convertedOrderBook.Bid != tCase.exp.Bid {
				t.Error("Incorrect convertation: ", tCase.exp, convertedOrderBook)
			}

		})
	}
}

func TestNewAuthMessage(t *testing.T) {
	cases := []struct {
		name string
		data [2]string
		exp  ascendex.AuthOutputMessage
	}{
		{
			name: "case1",
			data: [2]string{"sdgdfhsfdhsdgKprrr436345rew=ew1", "q0923r9dvfkosdjg22"},
			exp: ascendex.AuthOutputMessage{
				BaseOutputMessage: ascendex.BaseOutputMessage{
					Op: "auth",
				},
				TimeStamp: strconv.Itoa(int(time.Now().Unix())),
				ApiKey:    "q0923r9dvfkosdjg22",
				Signature: "CqEsp6vNK+E3R3JpsHyy+iEfDyAwK2O+nlBFX7ymGNk=",
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			generatedAuthMessage := ascendex.NewAuthMessage(tCase.data[0], tCase.data[1])

			if generatedAuthMessage.Signature != tCase.exp.Signature {
				t.Error("Incorrect signature generated!")
			}
		})
	}
}
