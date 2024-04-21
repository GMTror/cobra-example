package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"cobra-example/pkg/config"

	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xmux"
)

type Server struct {
	sync   sync.WaitGroup
	server http.Server
	logger xlog.Logger
	config *config.Config
}

func New(cfg *config.Config) *Server {
	l := xlog.New(xlog.Config{
		Level: cfg.LogLevel,
		Fields: xlog.F{
			"role": "http-server",
		},
		Output: xlog.NewLogfmtOutput(os.Stdout),
	})
	log.SetFlags(0)
	log.SetOutput(l)

	c := xhandler.Chain{}

	c.Use(xlog.RemoteAddrHandler("ip"))
	c.Use(xlog.UserAgentHandler("user_agent"))
	c.Use(xlog.RefererHandler("referer"))
	c.Use(xlog.RequestIDHandler("req_id", "Request-Id"))

	mux := xmux.New()

	mux.GET("/", xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	}))

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: c.HandlerCtx(context.Background(), mux),
	}

	return &Server{
		sync:   sync.WaitGroup{},
		server: srv,
		logger: l,
		config: cfg,
	}
}

func (s Server) Run() error {
	s.logger.Info("Run HTTP server")

	s.logger.Debugf("%#v", s.config)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
		s.server.Shutdown(context.Background())
	}()

	if err := s.server.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			s.logger.Error("HTTP server: ", err)
		}
	}
	s.sync.Wait()

	s.logger.Info("Stop")

	return nil
}
