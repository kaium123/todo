package cmd

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	log "github.com/zuu-development/fullstack-examination-2024/internal/log"
	"github.com/zuu-development/fullstack-examination-2024/internal/server"
)

func init() {
	rootCmd.AddCommand(NewServerCmd())
}

// NewServerCmd returns a new `server` command to be used as a sub-command to root
func NewServerCmd() *cobra.Command {
	serverCmd := cobra.Command{
		Use:     "server",
		Short:   "Print version information",
		Example: `  # Print the full version of client and server to stdout todo-cli server `,
		Run: func(_ *cobra.Command, _ []string) {
			var servers []server.Server

			logger := log.New()
			ctx := context.Background()

			initNewAPI := &server.InitNewAPI{
				TodoAPIServerOpts: server.TodoAPIServerOpts{
					ListenPort: cfg.APIServer.Port,
					Config:     cfg,
				},
				Log: logger,
			}

			apiServer, err := server.NewAPI(ctx, initNewAPI)
			if err != nil {
				logger.Fatal(ctx, "failed to init api.", zap.Error(err))
			}
			servers = append(servers, apiServer)

			if cfg.SwaggerServer.Enable {
				initNewSwagger := &server.InitNewSwagger{
					SwaggerServerOpts: server.SwaggerServerOpts{
						ListenPort: cfg.SwaggerServer.Port,
					},
					Log: logger,
				}

				swagServer := server.NewSwagger(ctx, initNewSwagger)
				servers = append(servers, swagServer)
			}

			ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
			defer stop()

			for _, s := range servers {
				server := s
				go func() {
					if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						logger.Fatal(ctx, fmt.Sprintf("shutting down %s. ", server.Name()), zap.Error(err))
					}
				}()
			}

			logger.Info(ctx, "server started")
			// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
			<-ctx.Done()
			logger.Info(ctx, "server shutting down")
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			for _, s := range servers {
				if err := s.Shutdown(ctx); err != nil {
					logger.Fatal(ctx, "error while shutting down.", zap.Error(err))
				}
			}
			logger.Info(ctx, "server shutdown gracefully")
		},
	}
	return &serverCmd
}
