package order

import (
	"errors"
	"log"
	"wbl0/internal/storage"
)

type Storage interface {
	Create(orderUid string, data string) error
	Get(uid string) (string, error)
	GetAll() []storage.OrderPair
}

type CacheStorage interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type Service struct {
	storage Storage
	cache   CacheStorage
}

func New(storage Storage, cache CacheStorage) *Service {
	orders := storage.GetAll()
	for _, order := range orders {
		cache.Set(order.OrderUID, order.Data)
	}
	return &Service{
		storage: storage,
		cache:   cache,
	}
}

func (svc *Service) Create(orderUID, data string) error {
	if err := svc.storage.Create(orderUID, data); err != nil {
		log.Printf("[order] [1H] [Create]: %s", err)
		return err
	}
	svc.cache.Set(orderUID, data)
	return nil
}

func (svc *Service) Get(uid string) (string, error) {
	log.Printf("[order][0][Get] trying to get: %s", uid)
	cacheData, ok := svc.cache.Get(uid)
	if ok {
		log.Printf("[order][1H][Get] cacheData: %s", cacheData)
		return cacheData.(string), nil
	}
	data, err := svc.storage.Get(uid)
	if err == nil {
		log.Printf("[order][2H] cacheData: %s", data)
		return data, nil
	}
	return "", errors.New("not found")
}
