package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"

	"github.com/slok/tracing-example/pkg/log"
	"github.com/slok/tracing-example/pkg/server"
)

// Main is the main app.
type Main struct {
	Flags *Flags
}

// Run runs the main program.
func (m *Main) Run() error {
	logger := log.Base().With("app", m.Flags.ServiceName)
	mux := http.NewServeMux()

	// Init the tracer
	tracer, closer, err := m.createTracer(m.Flags.ServiceName)
	if err != nil {
		return err
	}
	defer closer.Close()

	// Start server.
	s := server.New(m.Flags.Endpoints, tracer, mux, logger)
	errC := make(chan error)
	go func() {
		logger.Infof("listenig on: %s", m.Flags.listenAddress)
		errC <- http.ListenAndServe(m.Flags.listenAddress, s)
	}()

	// Capture signals.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errC:
		if err != nil {
			logger.Errorf("app finished with error: %s", err)
			return err
		}
		logger.Infof("app finished successfuly")
	case s := <-sigC:
		logger.Infof("signal %s received", s)
	}

	return nil
}

func (m *Main) createTracer(service string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot init Jaeger: %s", err)
	}
	return tracer, closer, nil
}

// Start the app.
func main() {
	m := &Main{
		Flags: NewFlags(),
	}
	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error during exection: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
