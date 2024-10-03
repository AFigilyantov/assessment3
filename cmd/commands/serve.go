package commands

import (
	"chitests/config"
	"chitests/internal/handlers"
	"chitests/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"

	// "github.com/go-chi/render"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	c := &cobra.Command{
		Use:        "serve",
		Aliases:    []string{"s"},
		SuggestFor: []string{},
		Short:      "Start API server",

		RunE: func(cmd *cobra.Command, args []string) error {
			//TODO start logger like slog

			log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			slog.SetDefault(log)

			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			router := chi.NewRouter()
			router.Use(middleware.RequestID) //registration of middlewares REALY NEED TODO
			router.Use(middleware.Recoverer)
			router.Use(middleware.Logger) // switcth off to production transfer to proxy server

			cfg, err := config.Parse("C:/Users/afigi/Desktop/Education/for_Chi/config.yaml")
			slog.Info("config", slog.Any("cfg", cfg))
			if err != nil {
				return err
			}

			s, err := storage.New("./storage.db")

			if err != nil {
				return err
			}
			// rate limiter by chi
			router.Use(httprate.Limit(
				3,             // requests
				1*time.Second, // per duration
				httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
			))

			router.Post("/register", handlers.RegisterHandler(&s))

			router.Get("/login", handlers.LoginHandler(&s))

			httpServer := http.Server{
				Addr:         cfg.HTTPServer.Address,
				ReadTimeout:  cfg.HTTPServer.ReadTimeout,
				WriteTimeout: cfg.HTTPServer.WriteTimeout,
				Handler:      router,
			}

			go func() {
				if err := httpServer.ListenAndServe(); err != nil {
					log.Info("ListenAndServe", slog.Any("err", err))

				}

			}()

			log.Info("sever start on:", slog.Any("sever start on :", cfg.HTTPServer.Address))

			<-ctx.Done()

			closeCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
			if err := httpServer.Shutdown(closeCtx); err != nil {
				return fmt.Errorf("serever close with error %s", err)

			}
			//close db connection
			s.CloseDb()
			//etc

			return nil
		},
	}
	return c
}
