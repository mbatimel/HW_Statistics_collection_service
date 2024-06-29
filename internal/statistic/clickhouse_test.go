package statistic

import (
	"testing"
	"time"

	"github.com/mbatimel/HW_Statistics_collection_service/internal/config"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/model"
)

func TestStatisticsService_GetOrderBook(t *testing.T) {
	cfg := config.ClickHouse{
		Host:     "localhost", // замените на ваш хост ClickHouse
		Port:     "9006",      // замените на ваш порт ClickHouse
		DB:       "my_database",   // замените на вашу тестовую базу данных ClickHouse
		Username: "my_user",
		Password: "my_password",
	}
	service, err := NewStatisticsService(cfg)
	if err != nil {
		t.Fatalf("failed to create StatisticsService: %v", err)
	}
	defer service.Close()

	

	exchangeName := "test_exchange"
	pair := "BTC/USD"

	_, err = service.GetOrderBook(exchangeName, pair)
	if err != nil {
		t.Errorf("GetOrderBook() error = %v", err)
	}
}

func TestStatisticsService_SaveOrderBook(t *testing.T) {
	cfg := config.ClickHouse{
		Host:     "localhost", // замените на ваш хост ClickHouse
		Port:     "9006",      // замените на ваш порт ClickHouse
		DB:       "my_database",   // замените на вашу тестовую базу данных ClickHouse
		Username: "my_user",
		Password: "my_password",
	}
	service, err := NewStatisticsService(cfg)
	if err != nil {
		t.Fatalf("failed to create StatisticsService: %v", err)
	}
	defer service.Close()


	exchangeName := "test_exchange"
	pair := "BTC/USD"
	orderBook := []*model.DepthOrder{
		{Price: 10000, BaseQty: 1},
		{Price: 10100, BaseQty: 2},
	}

	err = service.SaveOrderBook(exchangeName, pair, orderBook)
	if err != nil {
		t.Errorf("SaveOrderBook() error = %v", err)
	}
}

func TestStatisticsService_GetOrderHistory(t *testing.T) {
	cfg := config.ClickHouse{
		Host:     "localhost", // замените на ваш хост ClickHouse
		Port:     "9006",      // замените на ваш порт ClickHouse
		DB:       "my_database",   // замените на вашу тестовую базу данных ClickHouse
		Username: "my_user",
		Password: "my_password",
	}
	service, err := NewStatisticsService(cfg)
	if err != nil {
		t.Fatalf("failed to create StatisticsService: %v", err)
	}
	defer service.Close()


	client := &model.Client{
		ClientName:   "test_client",
		ExchangeName: "test_exchange",
		Label:        "test_label",
		Pair:         "BTC/USD",
	}

	_, err = service.GetOrderHistory(client)
	if err != nil {
		t.Errorf("GetOrderHistory() error = %v", err)
	}
}

func TestStatisticsService_SaveOrder(t *testing.T) {
	cfg := config.ClickHouse{
		Host:     "localhost", // замените на ваш хост ClickHouse
		Port:     "9006",      // замените на ваш порт ClickHouse
		DB:       "my_database",   // замените на вашу тестовую базу данных ClickHouse
		Username: "my_user",
		Password: "my_password",
	}
	service, err := NewStatisticsService(cfg)
	if err != nil {
		t.Fatalf("failed to create StatisticsService: %v", err)
	}
	defer service.Close()


	client := &model.Client{
		ClientName:   "test_client",
		ExchangeName: "test_exchange",
		Label:        "test_label",
		Pair:         "BTC/USD",
	}
	order := &model.HistoryOrder{
		ClientName:          "test_client",
		ExchangeName:        "test_exchange",
		Label:               "test_label",
		Pair:                "BTC/USD",
		Side:                "buy",
		TypeOrder:           "limit",
		BaseQty:             1,
		Price:               10000,
		AlgorithmNamePlaced: "test_algo",
		LowestSellPrc:       10100,
		HighestBuyPrc:       9900,
		CommissionQuoteQty:  10,
		TimePlaced:          time.Now(),
	}

	err = service.SaveOrder(client, order)
	if err != nil {
		t.Errorf("SaveOrder() error = %v", err)
	}
}
