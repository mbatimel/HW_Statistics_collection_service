package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/mbatimel/HW_Statistics_collection_service/internal/config"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/model"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/statistic"
)

var ErrChannelClosed = errors.New("channel is closed")

type Server interface {
	Run(ctx context.Context) error
	Close() error
}

type server struct {
	srv       *http.Server
	statistic statistic.IStatistics
}

func (s *server) Run(ctx context.Context) error {
	ch := make(chan error, 1)
	defer close(ch)
	go func() {
		ch <- s.srv.ListenAndServe()
	}()
	select {
	case err, ok := <-ch:
		if !ok {
			return ErrChannelClosed
		}
		if err != nil {
			return fmt.Errorf("failed to listen and serve: %w", err)
		}
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context failed: %w", err)
		}
	}
	return nil
}

func (s *server) Close() error {
	return nil
}

func NewServerConfig(cfg config.Config) (Server, error) {
	srv := http.Server{
		Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
	}
	statisticservic, err := statistic.NewStatisticsService(cfg.ClickHouse)
	if err != nil {
		return nil, fmt.Errorf("failed to create statistic: %w", err)
	}

	sv := server{
		srv:       &srv,
		statistic: statisticservic,
	}
	sv.setupRoutes()
	return &sv, nil
}

func (s *server) setupRoutes() {
	mx := http.NewServeMux()

	mx.HandleFunc("/get-order-book", s.handleGetOrderBook)
	mx.HandleFunc("/save-order-book", s.handleSaveOrderBook)
	mx.HandleFunc("/get-order-history", s.handleGetOrderHistory)
	mx.HandleFunc("/save-order-history", s.handleSaveOrderHistory)

	s.srv.Handler = mx
}

func (s *server) handleGetOrderBook(w http.ResponseWriter, r *http.Request) {
	exchangeName := r.URL.Query().Get("exchange_name")
	pair := r.URL.Query().Get("pair")

	orderBook, err := s.statistic.GetOrderBook(exchangeName, pair)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get order book: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderBook)
}

func (s *server) handleSaveOrderBook(w http.ResponseWriter, r *http.Request) {
    var orderBook []*model.DepthOrder
    if err := json.NewDecoder(r.Body).Decode(&orderBook); err != nil {
        http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
        return
    }

    exchangeName := r.URL.Query().Get("exchange_name")
    pair := r.URL.Query().Get("pair")

    // Example: Assuming statistic.SaveOrderBook takes []*model.DepthOrder
    if err := s.statistic.SaveOrderBook(exchangeName, pair, orderBook); err != nil {
        http.Error(w, fmt.Sprintf("failed to save order book: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}



func (s *server) handleGetOrderHistory(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	orderHistory, err := s.statistic.GetOrderHistory(&client)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get order history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderHistory)
}

func (s *server) handleSaveOrderHistory(w http.ResponseWriter, r *http.Request) {
	var order model.HistoryOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	var client model.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode client from request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := s.statistic.SaveOrder(&client, &order); err != nil {
		http.Error(w, fmt.Sprintf("failed to save order: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
