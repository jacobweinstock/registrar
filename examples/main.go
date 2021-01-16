package main

import (
	"context"
	"fmt"

	"github.com/jacobweinstock/registrar"
)

// your specific implementations
type driverOne struct {
	name     string
	protocol string
	metadata string
	features registrar.Features
}

type driverTwo struct {
	name     string
	protocol string
	metadata string
	features registrar.Features
}

func (d *driverOne) Compatible(ctx context.Context) bool {
	// do a compatibility check for driver one
	return true
}

func (d *driverOne) Thing() bool {
	return true
}

func (d *driverOne) Name() string {
	return d.name
}

func (d *driverTwo) Compatible(ctx context.Context) bool {
	// do a compatibility check for driver two
	return true
}

func (d *driverTwo) Thing() bool {
	return true
}

func (d *driverTwo) Name() string {
	return d.name
}

// CoolThinger does the cool thing you want it to do
type CoolThinger interface {
	Name() string
	Thing() bool
}

func main() {
	// create a registry
	reg := registrar.NewRegistry()

	// registry drivers
	one := &driverOne{name: "driverOne", protocol: "tcp", metadata: "this is driver one", features: registrar.Features{registrar.Feature("always double checking")}}
	two := &driverTwo{name: "driverTwo", protocol: "udp", metadata: "this is driver two", features: registrar.Features{registrar.Feature("set and forget")}}
	reg.Register(one.name, one.protocol, one.features, one.metadata, one)
	reg.Register(two.name, two.protocol, two.features, two.metadata, two)

	// do some filtering
	ctx := context.Background()
	reg.Drivers = reg.Using("tcp")
	reg.Drivers = reg.FilterForCompatible(ctx)

	// get the interfaces and run CoolThinger.Thing()
	var didTheThing bool
doTheThing:
	for _, elem := range reg.GetDriverInterfaces() {
		switch d := elem.(type) {
		case CoolThinger:
			didTheThing = d.Thing()
			fmt.Printf("cool thing executed by %v\n", d.Name())
			break doTheThing
		default:
			didTheThing = false
		}
	}
	fmt.Printf("did we do the thing? %v\n", didTheThing)
}
