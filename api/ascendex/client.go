package ascendex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"test/api"
)

const baseUrl = "wss://ascendex.com/%d/api/pro/v1/stream"

type client struct {
	userGroup  int
	apiKey     string
	apiSecret  string
	conn       *websocket.Conn
	connDone   chan struct{}
	input      chan []byte
	subscriber chan<- api.BestOrderBook
	output     chan []byte
}

func NewClient(apiKey string, apiSecret string, userGroup int) *client {
	return &client{
		userGroup:  userGroup,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		input:      make(chan []byte),
		subscriber: nil,
		output:     make(chan []byte),
	}
}

func (c *client) Connection() error {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(baseUrl, c.userGroup), nil)
	if err != nil {
		return err
	}

	c.conn = conn

	go c.handleConnection()

	err = c.expectMessage(connectedType)
	if err != nil {
		c.Disconnect()
		return err
	}

	authMes := NewAuthMessage(c.apiSecret, c.apiKey)
	err = c.sendMessage(authMes)
	if err != nil {
		c.Disconnect()
		return err
	}

	return c.expectMessage(authType)
}

func (c *client) Disconnect() {
	c.connDone <- struct{}{}
	c.conn.Close()
}

func (c *client) SubscribeToChannel(symbol string) error {
	channel := ParseBboChannel(symbol)

	err := c.sendMessage(NewSubMessage(channel))
	if err != nil {
		return err
	}

	return c.expectMessage(subType)
}

func (c *client) handleConnection() {
	for {
		select {
		case <-c.connDone:
			return
		default:

		}
		_, mData, err := c.conn.ReadMessage()

		if err != nil {
			log.Println(err)
			continue
		}

		inMessage, err := ParseInputMessage(mData)
		if err != nil {
			log.Println(err)
			continue
		}

		if inMessage.M == pingType {
			c.sendMessage(NewPongMessage())
		}

		if inMessage.M != bboType {
			c.input <- mData
			continue
		}

		if c.subscriber == nil {
			continue
		}

		bboMessage, err := ParseBboInputMessage(mData)
		order, err := bboMessage.ToBestOrderBook()
		if err != nil {
			log.Println(err)
			continue
		}

		select {
		case c.subscriber <- *order:
		default:
		}

	}
}

func (c *client) WriteMessagesToChannel() {
	for {
		data := <-c.output
		err := c.conn.WriteMessage(typeMessage, data)
		if err != nil {
			log.Println(err)
		}
	}

}

func (c *client) ReadMessagesFromChannel(ch chan<- api.BestOrderBook) {
	c.subscriber = ch
}

func (c *client) sendMessage(mes any) error {
	data, err := json.Marshal(&mes)
	if err != nil {
		return err
	}
	c.output <- data

	return nil
}

func (c *client) expectMessage(messageType string) error {
	mData := <-c.input

	message, err := ParseInputMessage(mData)
	if err != nil {
		return err
	}

	if message.M != messageType || message.Code != successCode {
		return errors.New(string(mData))
	}

	return nil
}
