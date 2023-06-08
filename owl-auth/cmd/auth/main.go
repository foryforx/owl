package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/foryforx/owl/owl-auth/api"
	"github.com/foryforx/owl/owl-auth/internal/app"
	"github.com/foryforx/owl/owl-auth/internal/conf"
	"github.com/foryforx/owl/owl-auth/internal/db/pg"
	"github.com/foryforx/owl/owl-auth/internal/transport/httpapi"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

var isProduction bool

const (
	PROD = "production"
)

func init() {

	// Init env variables
	if os.Getenv("ENV") == PROD {
		isProduction = true
	}

	conf.Initialize(isProduction)
}

func main() {
	log.Infoln("owl-go started")

	conf := conf.GetConfiguration()
	conf.Print()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, conf); err != nil {
		log.Fatalf("failed running app: %v", err)
	}

	log.Println("owl-go stopped")
}

func run(ctx context.Context, config *conf.Configuration) error {

	connPool, err := conf.GetConnPool(config)
	if err != nil {
		log.Fatalf("Error creating connType: %v", err)
	}

	if err := connPool.PingContext(ctx); err != nil {
		log.Fatalf("db ping: %w", err)
	}
	log.Infoln("starting app http server ", config.Port)

	var userRouter httpapi.Router
	{
		userStore := pg.NewUserStore(connPool)
		userService := app.NewUserService(userStore)
		accountsStore := pg.NewAccountStore(connPool)
		accountsService := app.NewAccountService(accountsStore)
		userRouter = &httpapi.UserRouter{UserService: userService, AccountService: accountsService}
	}

	var accountRouter httpapi.Router
	{
		accountsStore := pg.NewAccountStore(connPool)
		accountsService := app.NewAccountService(accountsStore)
		accountRouter = &httpapi.AccountRouter{AccountService: accountsService}
	}

	apiServer := &httpapi.APIServer{
		ConnPool:      connPool,
		Conf:          config,
		HealthChecker: doHealthCheck(connPool),
		Middleware:    httpapi.ChainMiddleware(httpapi.WithLogging()),
		Authenticator: httpapi.ChainMiddleware(httpapi.WithAuthenticator(app.NewUserService(pg.NewUserStore(connPool)))),
	}
	serveMux := apiServer.CreateMux(accountRouter, userRouter)
	return startHTTPServer(ctx, config.Port, serveMux)
}

func startHTTPServer(ctx context.Context, addr string, handler http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	errc := make(chan error, 1)

	log.Infoln("starting app http server ", addr)

	go func(errc chan<- error) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- fmt.Errorf("http server shutdonwn: %w", err)
		}
	}(errc)

	shutdown := func(timeout time.Duration) error {
		log.Infoln("received context cancellation; shutting down server")

		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(shutCtx); err != nil {
			return fmt.Errorf("http server shutdonwn: %w", err)
		}
		return nil
	}

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return shutdown(time.Second * 5)
	}
}

func doHealthCheck(db *sqlx.DB) httpapi.HealthCheckerFunc {
	return func(ctx context.Context) ([]api.Health, error) {
		if err := db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("db ping: %w", err)
		}

		var i int
		if err := db.Get(&i, "SELECT 1"); err != nil {
			return nil, fmt.Errorf("db read: %w", err)
		}

		return []api.Health{
			{Service: "owl-go", Status: "OK", Time: time.Now().Local().String()},
			{Service: "owl-db", Status: "OK", Time: time.Now().Local().String(), Details: db.Stats()},
		}, nil
	}
}
