package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Logger logs messages.
type Logger interface {
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
}

// Shutdown shuts down the given HTTP server gracefully when receiving an os.Interrupt or syscall.SIGTERM signal.
// It will wait for the specified timeout to stop hanging HTTP handlers.
func Shutdown(hs *http.Server, timeout time.Duration, logger Logger) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Infof("shutting down server with %s timeout", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Errorf("error while shutting down server: %v", err)
	} else {
		logger.Info("server was shut down gracefully")
	}
}
