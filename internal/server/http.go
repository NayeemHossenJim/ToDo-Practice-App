package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/NayeemHossenJim/ToDo-Practice-App/internal/config"
)

func NewHTTPServer(
	applicationConfig config.Config,
	router *gin.Engine,
) *http.Server {
	return &http.Server{
		Addr:              applicationConfig.HTTPAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func RegisterLifecycle(
	lifecycle fx.Lifecycle,
	httpServer *http.Server,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen(
				"tcp",
				httpServer.Addr,
			)
			if err != nil {
				return fmt.Errorf(
					"listen on %s: %w",
					httpServer.Addr,
					err,
				)
			}

			log.Printf(
				"Listening on %s",
				httpServer.Addr,
			)

			go func() {
				err := httpServer.Serve(listener)
				if err != nil &&
					!errors.Is(err, http.ErrServerClosed) {
					log.Printf(
						"HTTP server failed: %v",
						err,
					)
				}
			}()

			return nil
		},
		OnStop: func(context context.Context) error {
			log.Println("Stopping HTTP server")

			if err := httpServer.Shutdown(context); err != nil {
				return fmt.Errorf(
					"shutdown HTTP server: %w",
					err,
				)
			}

			return nil
		},
	})
}
