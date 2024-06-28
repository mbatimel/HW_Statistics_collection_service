package statistic

import (
	"github.com/mbatimel/HW_Statistics_collection_service/internal/model"
)
type IStatistics interface{
	
	GetOrderBook(exchange_name, pair string) ([]*model.DepthOrder, error)
	SaveOrderBook(exchange_name, pair string, orderBook []*model.DepthOrder) error
	GetOrderHistory(client *model.Client)  ([]*model.HistoryOrder, error)
	SaveOrder(client *model.Client, order *model.HistoryOrder) error

}