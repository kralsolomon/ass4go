package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (service *service) serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", service.config.port),
		Handler:      service.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		service.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		service.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": server.Addr,
		})
		service.wg.Wait()
		shutdownError <- nil
	}()

	service.logger.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  service.config.env,
	})

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	service.logger.PrintInfo("stopped server", map[string]string{
		"addr": server.Addr,
	})
	return nil
}
