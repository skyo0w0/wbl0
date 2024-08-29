package stan

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" validate:"required"`
	SmId              int64     `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	// add ID
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required"`
}

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestId    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       int64   `json:"amount" validate:"required"`
	PaymentDt    uint64  `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"required"`
	GoodsTotal   int64   `json:"goods_total" validate:"required"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ChrtId      uint64  `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Rid         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        float64 `json:"sale"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float64 `json:"total_price" validate:"required,gte=0"`
	NmId        uint64  `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      uint64  `json:"status" validate:"required"`
}
