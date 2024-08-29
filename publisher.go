package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

// Структуры для заказа
type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Orders struct {
	ValidOrders   []Order `json:"valid_orders"`
	InvalidOrders []Order `json:"invalid_orders"`
}

func sendOrder(sc stan.Conn, order Order) {
	// Сериализация заказа в JSON
	dataBytes, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order: %v", err)
	}

	// Отправка сообщения в канал "orders"
	err = sc.Publish("orders", dataBytes)
	if err != nil {
		log.Fatalf("Error during publishing: %v", err)
	}

	log.Printf("Order %s published successfully\n", order.OrderUID)
}

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-client")
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}
	defer sc.Close()

	file, err := os.Open("orders.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var orders Orders
	if err := json.Unmarshal(jsonData, &orders); err != nil {
		log.Fatalf("Error unmarshaling JSON data: %v", err)
	}

	for _, order := range orders.ValidOrders {
		sendOrder(sc, order)
	}

	for _, order := range orders.InvalidOrders {
		sendOrder(sc, order)
	}
}
