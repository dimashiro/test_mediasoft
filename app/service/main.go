package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimashiro/test_mediasoft/config"
	"github.com/dimashiro/test_mediasoft/internal/handler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	log, err := initLogger("GROUPMANAGE")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("start", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}

}

func run(log *zap.SugaredLogger) error {
	//__________________________________________________________________________
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("loading conf: %w", err)
	}

	//__________________________________________________________________________
	// App start
	log.Infow("start", "version", "develop")

	//__________________________________________________________________________
	// Start service
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	apiRouter, err := handler.NewRouter(ctx, log, cfg)
	if err != nil {
		return fmt.Errorf("can't init router: %s", err.Error())
	}

	apiSrv := http.Server{
		Addr:         cfg.APIHost,
		Handler:      apiRouter,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Infow("start", "status", "start router", "host", apiSrv.Addr)
		serverErrors <- apiSrv.ListenAndServe()
	}()

	// Gracefull shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "start shutdown", "signal", sig)
		defer log.Infow("shutdown", "status", "finish shutdown", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := apiSrv.Shutdown(ctx); err != nil {
			apiSrv.Close()
			return fmt.Errorf("could not stop server: %w", err)
		}
	}

	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
