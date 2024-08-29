package stan

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"log"
	"wbl0/internal/config"
)

type OrderService interface {
	Create(orderUID, data string) error
	Get(orderUID string) (string, error)
}

type Service struct {
	orderSvc    OrderService
	NatsConnect *nats.Conn
	StanConnect stan.Conn
	Sub         stan.Subscription

	url       string
	clusterID string
	clientID  string
	subject   string
}

func New(cfg *config.Stan, orderSvc OrderService) *Service {
	return &Service{
		orderSvc:  orderSvc,
		url:       cfg.URL,
		clusterID: cfg.ClusterID,
		clientID:  cfg.ClientID,
		subject:   cfg.Subject,
	}
}

func (s *Service) Start() {
	// Подключение к NATS
	nc, err := nats.Connect(s.url, nats.Name("Orders reader"))
	if err != nil {
		log.Fatal(err)
	}
	// Подключение к NATS Streaming (Stan)
	sc, err := stan.Connect(
		s.clusterID, s.clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}),
	)
	if err != nil {
		log.Fatalf("[stan.go] Can't connect: %v. \nMake sure a Nats Streaming Server is running at: %s", err, s.url)
	}
	log.Printf("[stan.go] Connected to %s clusterID: [%s] clientID: [%s]\n", s.url, s.clusterID, s.clientID)

	// Подписка на канал с обработчиком сообщений
	sub, err := sc.QueueSubscribe(
		s.subject,
		"", // Группа очереди (можно оставить пустым)
		s.handleMessage,
		stan.StartAt(pb.StartPosition_NewOnly),
		stan.DurableName(""),
	)
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}
	log.Printf("[stan.go] Listening on [%s], clientID=[%s]\n", s.subject, s.clientID)

	// Сохранение соединений и подписки
	s.Sub = sub
	s.NatsConnect = nc
	s.StanConnect = sc
}

func orderIsValid(order *Order) bool {
	validate := validator.New()
	err := validate.Struct(order)
	return err == nil
}

func (s *Service) handleMessage(m *stan.Msg) {
	var order Order
	log.Printf("Received: %s\n", m)
	if err := json.Unmarshal(m.Data, &order); err != nil {
		log.Printf("[1H] err:\n")
		log.Println(err)
		return
	}
	log.Printf("[1] err:\n")
	if !orderIsValid(&order) {
		log.Printf("[2H] err:\n")
		return
	}
	log.Printf("[2] err:\n")
	dataBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("[3H] err:\n")
		return
	}
	log.Printf("[3] err:\n")
	data := string(dataBytes)
	log.Printf("data [%s]\n\n", data)
	s.orderSvc.Create(order.OrderUID, data)
}

func (s *Service) Stop() {
	if err := s.Sub.Unsubscribe(); err != nil {
		log.Printf("Error during unsubscribe: %v", err)
	}
	if err := s.StanConnect.Close(); err != nil {
		log.Printf("Error during Stan connection close: %v", err)
	}
	if s.NatsConnect != nil && !s.NatsConnect.IsClosed() {
		s.NatsConnect.Close()
	}
	log.Println("[stan.go] Service stopped")
}
