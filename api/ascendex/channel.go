package ascendex

import (
	"fmt"
	"strings"
)

type BboChannel struct {
	Token string
	Asset string
}

func ParseBboChannel(symbol string) BboChannel {
	tokens := strings.Split(symbol, "_")
	return BboChannel{
		tokens[0],
		tokens[1],
	}
}

func (c *BboChannel) String() string {
	return fmt.Sprintf("bbo:%s/%s", c.Token, c.Asset)
}
