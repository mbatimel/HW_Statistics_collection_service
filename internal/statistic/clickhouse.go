package statistic

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/config"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/model"
)

type StatisticsService struct {
	conn driver.Conn
}

func NewStatisticsService(cfg config.ClickHouse) (*StatisticsService, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.DB,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		DialTimeout:      5 * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenRoundRobin,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %v", err)
	}

	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse after connection: %v", err)
	}

	return &StatisticsService{
		conn: conn,
	}, nil
}

func (s *StatisticsService) Close() error {
	return s.conn.Close()
}

func (s *StatisticsService) GetOrderBook(exchangeName, pair string) ([]*model.DepthOrder, error) {
	ctx := context.Background()
	query := `
		SELECT asks, bids
		FROM OrderBook
		WHERE exchange = ? AND pair = ?
	`
	rows, err := s.conn.Query(ctx, query, exchangeName, pair)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for order book: %v", err)
	}
	defer rows.Close()

	var orderBook []*model.DepthOrder
	for rows.Next() {
		var asks, bids []struct {
			Price    float64
			BaseQty  float64
		}
		if err := rows.Scan(&asks, &bids); err != nil {
			return nil, fmt.Errorf("failed to scan row for order book: %v", err)
		}
		for _, ask := range asks {
			orderBook = append(orderBook, &model.DepthOrder{Price: ask.Price, BaseQty: ask.BaseQty})
		}
		for _, bid := range bids {
			orderBook = append(orderBook, &model.DepthOrder{Price: bid.Price, BaseQty: bid.BaseQty})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over order book rows: %v", err)
	}

	return orderBook, nil
}

func (s *StatisticsService) SaveOrderBook(exchangeName, pair string, orderBook []*model.DepthOrder) error {
	ctx := context.Background()
	batch, err := s.conn.PrepareBatch(ctx, "INSERT INTO OrderBook (exchange, pair, asks, bids)")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %v", err)
	}

	var asks, bids []struct {
		Price   float64
		BaseQty float64
	}
	for _, order := range orderBook {
		if order.Price > 0 {
			asks = append(asks, struct {
				Price   float64
				BaseQty float64
			}{order.Price, order.BaseQty})
		} else {
			bids = append(bids, struct {
				Price   float64
				BaseQty float64
			}{order.Price, order.BaseQty})
		}
	}

	if err := batch.Append(exchangeName, pair, asks, bids); err != nil {
		return fmt.Errorf("failed to append to batch: %v", err)
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %v", err)
	}

	return nil
}

func (s *StatisticsService) GetOrderHistory(client *model.Client) ([]*model.HistoryOrder, error) {
	ctx := context.Background()
	query := `
		SELECT client_name, exchange_name, label, pair, side, type_order,
			   base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc,
			   commission_quote_qty, time_placed
		FROM HistoryOrder
		WHERE client_name = ? AND exchange_name = ? AND label = ? AND pair = ?
	`
	rows, err := s.conn.Query(ctx, query, client.ClientName, client.ExchangeName, client.Label, client.Pair)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for order history: %v", err)
	}
	defer rows.Close()

	var orderHistory []*model.HistoryOrder
	for rows.Next() {
		var historyOrder model.HistoryOrder
		if err := rows.Scan(
			&historyOrder.ClientName, &historyOrder.ExchangeName, &historyOrder.Label,
			&historyOrder.Pair, &historyOrder.Side, &historyOrder.TypeOrder,
			&historyOrder.BaseQty, &historyOrder.Price, &historyOrder.AlgorithmNamePlaced,
			&historyOrder.LowestSellPrc, &historyOrder.HighestBuyPrc, &historyOrder.CommissionQuoteQty,
			&historyOrder.TimePlaced,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row for order history: %v", err)
		}
		orderHistory = append(orderHistory, &historyOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over order history rows: %v", err)
	}

	return orderHistory, nil
}

func (s *StatisticsService) SaveOrder(client *model.Client, order *model.HistoryOrder) error {
	ctx := context.Background()
	query := `
		INSERT INTO HistoryOrder (client_name, exchange_name, label, pair, side, type_order,
								   base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc,
								   commission_quote_qty, time_placed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	if err := s.conn.Exec(ctx, query,
		client.ClientName, client.ExchangeName, client.Label, client.Pair,
		order.Side, order.TypeOrder, order.BaseQty, order.Price, order.AlgorithmNamePlaced,
		order.LowestSellPrc, order.HighestBuyPrc, order.CommissionQuoteQty, order.TimePlaced,
	); err != nil {
		return fmt.Errorf("failed to execute insert query for order history: %v", err)
	}

	return nil
}
