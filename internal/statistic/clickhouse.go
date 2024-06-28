package statistic

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/config"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/model"
)
type StatisticsService struct {
    db *sql.DB
}
func NewStatisticsService(cfg config.ClickHouse) (*StatisticsService, error) {
    connStr := fmt.Sprintf("tcp://%s:%s?database=%s", cfg.Host, cfg.Port, cfg.DB)
    fmt.Println("Connecting to ClickHouse with:", connStr)
    db, err := sql.Open("clickhouse", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to ClickHouse: %v", err)
    }

    // Check if the connection is successful
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping ClickHouse after connection: %v", err)
    }

    return &StatisticsService{
        db: db,
    }, nil
}

func (s *StatisticsService) Close() error {
    return s.db.Close()
}

func (s *StatisticsService) GetOrderBook(exchangeName, pair string) ([]*model.DepthOrder, error) {
    query := `
        SELECT price, base_qty
        FROM order_book
        WHERE exchange = ? AND pair = ?
    `
    
    rows, err := s.db.Query(query, exchangeName, pair)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query for order book: %v", err)
    }
    defer rows.Close()
    
    var orderBook []*model.DepthOrder
    for rows.Next() {
        var depthOrder model.DepthOrder
        if err := rows.Scan(&depthOrder.Price, &depthOrder.BaseQty); err != nil {
            return nil, fmt.Errorf("failed to scan row for order book: %v", err)
        }
        orderBook = append(orderBook, &depthOrder)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over order book rows: %v", err)
    }
    
    return orderBook, nil
}



func (s *StatisticsService) SaveOrderBook(exchangeName, pair string, orderBook []*model.DepthOrder) error {
    // Begin transaction
    tx, err := s.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to start transaction: %v", err)
    }
    
    // Insert each DepthOrder into order_book table
    for _, order := range orderBook {
        query := `
            INSERT INTO order_book (exchange, pair, price, base_qty)
            VALUES (?, ?, ?, ?)
        `
        _, err := tx.Exec(query, exchangeName, pair, order.Price, order.BaseQty)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to execute insert query: %v", err)
        }
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }
    
    return nil
}


func (s *StatisticsService) GetOrderHistory(client *model.Client) ([]*model.HistoryOrder, error) {
    query := `
        SELECT client_name, exchange_name, label, pair, side, type_order,
               base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc,
               commission_quote_qty, time_placed
        FROM order_history
        WHERE client_name = ? AND exchange_name = ? AND label = ? AND pair = ?
    `
    
    rows, err := s.db.Query(query, client.ClientName, client.ExchangeName, client.Label, client.Pair)
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
    query := `
        INSERT INTO order_history (client_name, exchange_name, label, pair, side, type_order,
                                   base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc,
                                   commission_quote_qty, time_placed)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := s.db.Exec(
        query, client.ClientName, client.ExchangeName, client.Label, client.Pair,
        order.Side, order.TypeOrder, order.BaseQty, order.Price, order.AlgorithmNamePlaced,
        order.LowestSellPrc, order.HighestBuyPrc, order.CommissionQuoteQty, order.TimePlaced,
    )
    if err != nil {
        return fmt.Errorf("failed to execute insert query for order history: %v", err)
    }
    
    return nil
}