package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/dimasmir03/web-calculator-server/docs"

	"github.com/dimasmir03/web-calculator-server/internal/logging"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http_server"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http_server/router"
	"github.com/dimasmir03/web-calculator-server/pkg/calculator/cmd/calculator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки файла .env: %v", err)
	}

	// config := new(Config)
	// config.Addr = os.Getenv("PORT")
	// if config.Addr == "" {
	// 	config.Addr = "5353"
	// }
	return nil, nil
}

type Application struct {
	config *Config
	log    *logrus.Logger
	r      *echo.Echo
}

func New() *Application {
	cfg, err := ConfigFromEnv()
	if err != nil {
		log.Printf("Err New App: %v", err)
	}
	return &Application{
		config: cfg,
		log:    logrus.New(),
		r:      echo.New(),
	}
}

func (a *Application) Run() {
	// a.r.Use(Logging)
	// a.r.HandleFunc("/api/v1/calculate", CalculationHandler)
	// a.log.Infof("Starting server on :%s", a.config.Addr)
	// if err := http.ListenAndServe(":"+a.config.Addr, a.r); err != nil {
	// 	a.log.Errorf("Failed start server: %s", err.Error())
	// }

	// Create context
	ctx := context.Background()

	// Create channel for gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	calc := calculator.NewCalculator()

	// Create router
	// router := echo.NewRouter()
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_unix}, host=${host}, method=${method}, uri=${uri}, status=${status}, error=${error} user_agent=${user_agent} latency_human=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Static("/docs", "docs")
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	router.InternalRouter(e, calc)
	router.ApiRouter(e, calc)
	router.SwaggerRouter(e)
	fmt.Println(e.Routers())
	// API subroute
	// apiRouter := router.PathPrefix("/api").Subrouter()
	// apiRouter.Use(logging.LoggingMiddleware)

	// Create new server and start
	apiServer := http_server.NewServer(fmt.Sprintf("%s:%s", "0.0.0.0", "8080"), e, nil, logging.Logger)
	go apiServer.Run()

	// gracefully shutdown
	// waiting for interrupt signal
	<-stop
	// Stop server with timeout
	apiServer.Stop(ctx)
}
