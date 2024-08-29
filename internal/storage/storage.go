package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	db *sqlx.DB
}

type OrderPair struct {
	OrderUID string `db:"order_uid"`
	Data     string `db:"data"`
}

func New(connectionString string) *Storage {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = migrate(db); err != nil {
		return nil
	}

	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(orderUid string, data string) error {
	query := `INSERT INTO orders (order_uid, data) VALUES ($1, $2)`
	_, err := s.db.Exec(query, orderUid, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Get(uid string) (string, error) {
	var data string
	query := `SELECT data FROM orders WHERE order_uid = $1`
	err := s.db.Get(&data, query, uid)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (s *Storage) GetAll() []OrderPair {
	var orders []OrderPair
	query := `SELECT order_uid, data FROM orders`
	err := s.db.Select(&orders, query)
	if err != nil {
		log.Printf("Failed to get all orders: %v", err)
		return nil
	}
	return orders
}
