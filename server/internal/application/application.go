package application

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dimasmir03/web-calculator-server/internal/transport/grpc/api"
	"google.golang.org/grpc"

	_ "github.com/dimasmir03/web-calculator-server/docs"
	"github.com/dimasmir03/web-calculator-server/internal/auth"
	"github.com/dimasmir03/web-calculator-server/internal/calculator/cmd/calculator"
	"github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	mygrpc "github.com/dimasmir03/web-calculator-server/internal/transport/grpc"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http"
	"github.com/dimasmir03/web-calculator-server/internal/transport/http/router"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Config struct {
	Addr      string
	DB_URL    string `env:"DB_URL" env-default:"sqlite.db"`
	JWTSecret string `env:"JWT_SECRET" env-default:"secret"`
}

func ConfigFromEnv() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("ошибка загрузки переменных из файла .env: %v", err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка загрузки переменных из окружения: %v", err)
	}

	return &cfg, nil
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

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logger.Debug("Config: ", cfg)
	logger.Debug(strconv.Atoi(os.Getenv("TIME_MULTIPLICATION_MS")))

	return &Application{
		config: cfg,
		log:    logger,
		r:      echo.New(),
	}
}

func (a *Application) Run() {
	// Create context
	ctx := context.Background()

	// Create channel for gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	db, err := sqlite.NewStorage(a.config.DB_URL)
	if err != nil {
		a.log.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		a.log.Fatal(err)
	}
	authService := auth.NewService(db, a.config.JWTSecret)
	authHandler := auth.Handler{Service: authService}
	fmt.Println(authService)

	calc := calculator.NewCalculator(db)
	if err = calc.RestoreExpressions(); err != nil {
		a.log.Fatal("ошибка восстановления выражений: ", err)
	}

	// Create router
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_unix}, host=${host}, method=${method}, uri=${uri}, status=${status}, error=${error} user_agent=${user_agent} latency_human=${latency_human}\n",
	}))
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Static("/docs", "docs")
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/api/v1/register", authHandler.Register)
	//	@Summary		Register new user
	//	@Description	Register new user
	//	@Tags			auth
	//	@Accept			json
	//	@Produce		json
	//	@Param			data	body		auth.RegisterRequest	true	"User data"
	//	@Success		200		{object}	model.User
	//	@Failure		400		{object}	echo.HTTPError
	//	@Failure		500		{object}	echo.HTTPError
	//	@Router			/api/v1/register [post]

	e.POST("/api/v1/login", authHandler.Login)
	//	@Summary		Login user
	//	@Description	Login user
	//	@Tags			auth
	//	@Accept			json
	//	@Produce		json
	//	@Param			data	body		auth.LoginRequest	true	"User data"
	//	@Success		200		{string}	token
	//	@Failure		400		{object}	echo.HTTPError
	//	@Failure		401		{object}	echo.HTTPError
	//	@Failure		500		{object}	echo.HTTPError
	//	@Router			/api/v1/login [post]
	authorized := e.Group("/api/v1/")
	authorized.Use(authService.JWTMiddleware())

	router.InternalRouter(e, calc)
	router.ApiRouter(authorized, db, calc)
	router.SwaggerRouter(e)
	a.log.Infof("API routes: %v", e.Routers())

	grpcServer := grpc.NewServer()
	api.RegisterCalculatorServer(grpcServer, &mygrpc.Server{Calc: calc})

	go func() {
		// Запуск gRPC
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			a.log.Fatalf("failed to serve: %v", err)
		}
		a.log.Infof("Starting gRPC server on %s", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			a.log.Fatalf("failed to serve: %v", err)
		}
	}()

	// API subroute
	// apiRouter := router.PathPrefix("/api").Subrouter()
	// apiRouter.Use(logging.LoggingMiddleware)

	// Create new server and start
	apiServer := http.NewServer(fmt.Sprintf("%s:%s", "0.0.0.0", "8080"), e, nil, a.log)
	go apiServer.Run()

	// gracefully shutdown
	// waiting for interrupt signal
	<-stop
	// Stop server with timeout
	apiServer.Stop(ctx)
	grpcServer.Stop()
	a.log.Infof("Stopping API server")
}
