package ascendex

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"test/api"
	"time"
)

const connectedType = "connected"
const pingType = "pong"
const bboType = "bbo"
const authType = "auth"
const subType = "sub"
const typeMessage = 1
const successCode = 0
const authPath = "stream"

type BaseOutputMessage struct {
	Op string `json:"op"`
}

type AuthOutputMessage struct {
	BaseOutputMessage
	TimeStamp string `json:"t"`
	ApiKey    string `json:"key"`
	Signature string `json:"sig"`
}

type SubOutputMessage struct {
	BaseOutputMessage
	Channel string `json:"ch"`
}

type InputMessage struct {
	M    string `json:"m"`
	Code int    `json:"code"`
	Ch   string `json:"ch"`
}

type BboInputMessage struct {
	M      string `json:"m"`
	Symbol string `json:"symbol"`
	Data   struct {
		Ts  int64    `json:"ts"`
		Bid []string `json:"bid"`
		Ask []string `json:"ask"`
	} `json:"data"`
}

func (message *BboInputMessage) ToBestOrderBook() (*api.BestOrderBook, error) {
	ask, err := makeOrderFromArray(message.Data.Ask)
	if err != nil {
		return nil, err
	}

	bid, err := makeOrderFromArray(message.Data.Bid)
	if err != nil {
		return nil, err
	}

	return &api.BestOrderBook{
		Ask: *ask,
		Bid: *bid,
	}, nil
}

func ParseInputMessage(data []byte) (InputMessage, error) {
	mes := InputMessage{}
	err := json.Unmarshal(data, &mes)

	return mes, err
}

func ParseBboInputMessage(data []byte) (BboInputMessage, error) {
	mes := BboInputMessage{}
	err := json.Unmarshal(data, &mes)

	return mes, err
}

func NewAuthMessage(secret, apiKey string) AuthOutputMessage {
	t := strconv.Itoa(int(time.Now().Unix()))
	mes := strings.Join([]string{t, authPath}, "+")
	sig := sign(mes, secret)

	return AuthOutputMessage{
		BaseOutputMessage{authType},
		t,
		apiKey,
		sig,
	}
}

func NewSubMessage(channel BboChannel) SubOutputMessage {
	return SubOutputMessage{
		BaseOutputMessage{subType},
		channel.String(),
	}
}

func NewPongMessage() BaseOutputMessage {
	return BaseOutputMessage{pingType}
}

func sign(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
