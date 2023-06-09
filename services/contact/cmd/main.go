package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"advanced.microservices/pkg/jsonlog"
	"advanced.microservices/pkg/store"
	"advanced.microservices/pkg/store/postgres"
	"advanced.microservices/services/contact/internal/delivery"
	"advanced.microservices/services/contact/internal/repository"
	"advanced.microservices/services/contact/internal/useCase"
	"github.com/julienschmidt/httprouter"
)

type config struct {
	port int
	env  string
	db   store.DbConfig
}

type service struct {
	config config
	logger *jsonlog.Logger
	wg     *sync.WaitGroup
	db     *sql.DB
	router *httprouter.Router
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.Dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	db, err := postgres.OpenDB(cfg.db)
	if err != nil {
		fmt.Println("Error while opening DB " + err.Error())
		return
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	router := httprouter.New()
	contactRepository := repository.NewContactRepository(db)
	contactUseCase := useCase.NewContactUsecase(contactRepository, 6*time.Second)
	delivery.NewContactHandler(router, logger, contactUseCase)

	groupRepository := repository.NewGroupRepository(db)
	groupUseCase := useCase.NewGroupUsecase(groupRepository, 6*time.Second)
	delivery.NewGroupHandler(router, logger, groupUseCase)

	service := &service{
		config: cfg,
		db:     db,
		logger: logger,
		router: router,
		// mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = service.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
