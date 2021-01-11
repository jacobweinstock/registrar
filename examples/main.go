package main

import (
	"context"
	"fmt"

	"github.com/jacobweinstock/registry"
)

type driverOne struct {
	name     string
	protocol string
	user     string
	pass     string
	host     string
	port     string
	Features registry.Features
}

func (do *driverOne) Compatible(ctx context.Context) bool {
	fmt.Println("compatible method")
	fmt.Printf("%+v\n", do)
	return true
}

func (do *driverOne) GetPower(ctx context.Context) (string, error) {
	fmt.Printf("%+v\n", do)
	return "on", nil
}

type driverTwo struct {
	name     string
	protocol string
	user     string
	pass     string
	host     string
	port     string
	Features registry.Features
}

func (do *driverTwo) Compatible(ctx context.Context) bool {
	fmt.Println("compatible method")
	fmt.Printf("%+v\n", do)
	return true
}

func (do *driverTwo) GetPower(ctx context.Context) (string, error) {
	fmt.Printf("%+v\n", do)
	return "on", nil
}

// PowerGetter interface
type PowerGetter interface {
	GetPower(ctx context.Context) (string, error)
}

func main() {
	reg := registry.NewRegistry()
	do := &driverOne{
		name:     "driverOne",
		protocol: "ipmi",
		Features: registry.Features{"power"},
		user:     "admin",
		pass:     "admin",
		host:     "localhost",
		port:     "623",
	}
	reg.Register(do.name, do.protocol, do, do, do.Features)
	d2 := &driverTwo{
		name:     "driverTwo",
		protocol: "ipmi",
		Features: registry.Features{"power"},
		user:     "admin",
		pass:     "admin",
		host:     "localhost",
		port:     "623",
	}
	reg.Register(d2.name, d2.protocol, d2, d2, d2.Features)

	fmt.Printf("%+v\n", reg[0].DriverInterface)
	iface := reg[0].DriverInterface

	fmt.Printf("compatible: %v\n", reg[0].Compatible(context.Background()))
	state, err := iface.(PowerGetter).GetPower(context.Background())
	fmt.Printf("%v, %v\n", state, err)
	fmt.Printf("%+v\n", reg[0])
}
