package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"booking-service/internal/api/rest"
	"booking-service/internal/api/rest/auth"
	"booking-service/internal/api/rest/bookings"
	"booking-service/internal/api/rest/business-account"
	"booking-service/internal/api/rest/middlewares"
	"booking-service/internal/api/rest/specialists"
	"booking-service/internal/api/rest/user-account"
	"booking-service/internal/store/users"
	"booking-service/pkg/db"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := LoadConfig()

	errChan := make(chan error)

	ctx := context.Background()

	dbConn, err := db.NewDB(ctx, cfg.DBReadURL, cfg.DBWriteURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to the database")
		return
	}

	defer dbConn.Close()

	usersStore := users.NewStore(dbConn.ReadPool, dbConn.WritePool)

	go func() {
		server := &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.Port),
			Handler:           setUpRouter(cfg, usersStore),
			ReadHeaderTimeout: 2 * time.Second,
		}
		log.Info().Msgf("Starting api service at port %d", cfg.Port)

		errChan <- server.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		log.Warn().Str("sig", sig.String()).Msg("got termination signal, exiting")
	case err := <-errChan:
		log.Error().Err(err).Msg("received error from error channel")
	}
}

func setUpRouter(cnf *Config, usersStore users.Store) *mux.Router {
	authMiddleware := middlewares.NewJWTMiddleware(cnf.JWTSecret)
	authHandler := auth.NewHandler(cnf.GoogleLoginConfig, cnf.GoogleRandomState, cnf.JWTSecret, cnf.JWTExpPeriod, usersStore)
	specialistsHandler := specialists.NewHandler()
	specialistsRouter := specialists.NewRouter(specialistsHandler, authMiddleware.Middleware)

	authRouter := auth.NewRouter(authHandler)
	bookingsHandler := bookings.NewHandler()
	bookingsRouter := bookings.NewRouter(bookingsHandler, authMiddleware.Middleware)

	businessAccountHandler := business_account.NewHandler()
	businessAccountRouter := business_account.NewRouter(businessAccountHandler, authMiddleware.Middleware)

	userAccountHandler := user_account.NewHandler(usersStore)
	userAccountRouter := user_account.NewRouter(userAccountHandler, authMiddleware.Middleware)

	routes := []rest.Register{
		authRouter,
		specialistsRouter,
		bookingsRouter,
		businessAccountRouter,
		userAccountRouter,
	}
	return rest.NewRouter(routes)
}
