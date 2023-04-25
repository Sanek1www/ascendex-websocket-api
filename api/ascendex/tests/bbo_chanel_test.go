package tests

import (
	"test/api/ascendex"
	"testing"
)

func TestParseBboChannel(t *testing.T) {
	cases := []struct {
		name string
		data string
		exp  string
	}{
		{
			name: "case1",
			data: "BTC_USDT",
			exp:  "bbo:BTC/USDT",
		},
		{
			name: "case1",
			data: "ETH_DOGE",
			exp:  "bbo:ETH/DOGE",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			channel := ascendex.ParseBboChannel(tCase.data)

			if channel.String() != tCase.exp {
				t.Error("Incorrect parsed!")
			}
		})
	}
}
