package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/config"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/statistic"
)
var ErrChannelClosed = errors.New("channel is closed")
type Server interface {
	Run(ctx context.Context) error
	Close() error

}
type server struct {
	srv *http.Server
	statistic statistic.IStatistics
}

func (s *server) Run(ctx context.Context) error{
	ch:=make(chan error, 1)
	defer close(ch)
	go func(){
		ch <- s.srv.ListenAndServe()
	}()
	select  {
	case err, ok := <-ch:
		if !ok{
			return ErrChannelClosed
		}
		if err != nil{
			return fmt.Errorf("failed to listen and serve: %w", err)
		}
	case <-ctx.Done():
		if err:=ctx.Err();err!=nil{
			return fmt.Errorf("context faild: %w", err)
		}
			
	}
	return nil
}
func (s *server) Close() error{
	return nil
}

func NewServerConfig(cfg config.Config) (Server, error){
	srv:= http.Server{
		Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
	}
	var cacheMemory cache.ICache
	var err error
	if !cfg.Cache.InMemory{ 
		cacheMemory, err = cache.NewRedisCache(cfg.Cache)
		if err != nil {
			
			return nil,fmt.Errorf("failed to create cache: %w", err)
		}
	}else{
		cacheMemory = cache.NewInMemoryCache(cfg.Cache.Cap)
	}
	sv := server{
		srv :&srv,
		cache: cacheMemory,
	}
	sv.setupRoutes()
	return &sv, nil
}

func (s *server)setupRoutes(){
	mx :=http.NewServeMux()

	mx.HandleFunc("/add",s.handleAdd)
	mx.HandleFunc("/clear",s.handleClear)
	mx.HandleFunc("/addttl",s.handleAddWithTTL)
	mx.HandleFunc("/get",s.handleGet)
	mx.HandleFunc("/remove",s.handleRemove)

	s.srv.Handler = mx
	
}
