package main

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/jacobweinstock/registrar"
	"go.uber.org/zap"
)

type driverOne struct {
	name     string
	protocol string
	user     string
	pass     string
	host     string
	port     string
	Features registrar.Features
	Log      logr.Logger
}

func (do *driverOne) Compatible(ctx context.Context) bool {
	do.Log.V(0).Info("compatible method")
	do.Log.V(0).Info("debugging", "driverOne", do)
	return true
}

func (do *driverOne) GetPower(ctx context.Context) (string, error) {
	do.Log.V(0).Info("debugging", "driverOne", do)
	return "on", nil
}

type driverTwo struct {
	name     string
	protocol string
	user     string
	pass     string
	host     string
	port     string
	Features registrar.Features
	Log      logr.Logger
}

func (do *driverTwo) Compatiblea(ctx context.Context) bool {
	do.Log.V(0).Info("compatible method")
	do.Log.V(0).Info("debugging", "driverTwo", do)
	return true
}

func (do *driverTwo) GetPower(ctx context.Context) (string, error) {
	do.Log.V(0).Info("debugging", "driverTwo", do)
	return "on", nil
}

// PowerGetter interface
type PowerGetter interface {
	GetPower(ctx context.Context) (string, error)
}

func defaultLogger() logr.Logger {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return zapr.NewLogger(logger)
}

func main() {
	log := defaultLogger()
	reg := registrar.NewRegistry(registrar.WithLogger(log))
	do := &driverOne{
		name:     "driverOne",
		protocol: "ipmi",
		Features: registrar.Features{"power"},
		user:     "admin",
		pass:     "admin",
		host:     "localhost",
		port:     "623",
		Log:      log,
	}
	reg.Register(do.name, do.protocol, do.Features, nil, do)
	d2 := &driverTwo{
		name:     "driverTwo",
		protocol: "ipmi",
		Features: registrar.Features{"power"},
		user:     "admin",
		pass:     "admin",
		host:     "localhost",
		port:     "623",
		Log:      log,
	}
	reg.Register(d2.name, d2.protocol, d2.Features, nil, d2)
	log.V(0).Info("debugging", "driver interface", reg.Drivers[0].DriverInterface)
	iface := reg.Drivers[0].DriverInterface
	ctx := context.Background()
	reg.FilterForCompatible(ctx)
	log.V(0).Info("debugging", "compatible", reg.Drivers[0].DriverInterface.(registrar.Verifier).Compatible(context.Background()))
	state, err := iface.(PowerGetter).GetPower(ctx)
	log.V(0).Info("debugging", "state", state, "err", err)
	log.V(0).Info("debugging", "driver[0]", reg.Drivers[0])
}
