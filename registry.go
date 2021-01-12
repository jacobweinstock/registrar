package registry

import (
	"context"
	"strings"
	"sync"
)

// Features holds the features a driver supports
type Features []Feature

// Feature represents a single feature a driver supports
type Feature string

// Drivers holds a slice of Driver types
type Drivers []*Driver

// Verifier allows implementations to define a method for
// determining whether a driver is compatible for use
type Verifier interface {
	Compatible(context.Context) bool
}

// Driver holds the info about a driver
type Driver struct {
	Name            string
	Protocol        string
	Features        Features
	DriverInterface interface{}
	Verifier
}

// NewRegistry returns a new Driver registry
func NewRegistry() Drivers {
	return make(Drivers, 0)
}

// Register will add a driver a Driver registry
func (d *Drivers) Register(name, protocol string, driverInterface interface{}, compatFn Verifier, features Features) {
	*d = append(*d, &Driver{
		Name:            name,
		Protocol:        protocol,
		Features:        features,
		Verifier:        compatFn,
		DriverInterface: driverInterface,
	})
}

// GetDriverInterfaces returns a slice of just the generic driver interfaces
func (d Drivers) GetDriverInterfaces() []interface{} {
	results := make([]interface{}, 0)
	for _, elem := range d {
		results = append(results, elem.DriverInterface)
	}
	return results
}

// FilterForCompatible updates the driver registry with only compatible implementations
func (d *Drivers) FilterForCompatible(ctx context.Context) {
	var wg sync.WaitGroup
	result := make(Drivers, 0)
	for _, elem := range *d {
		wg.Add(1)
		go func(isCompat Verifier, reg *Driver, wg *sync.WaitGroup) {
			if isCompat.Compatible(ctx) {
				result = append(result, reg)
			}
			wg.Done()
		}(elem.Verifier, elem, &wg)
	}
	wg.Wait()
	*d = result
}

// include does the actual work of filtering for specific features
func (f Features) include(features ...Feature) bool {
	if len(features) > len(f) {
		return false
	}
	fKeys := make(map[Feature]bool)
	for _, v := range f {
		fKeys[v] = true
	}
	for _, f := range features {
		if _, ok := fKeys[f]; !ok {
			return false
		}
	}
	return true
}

// Supports does the actual work of filtering for specific features
func (d Drivers) Supports(features ...Feature) Drivers {
	supportedRegistries := make(Drivers, 0)
	for _, reg := range d {
		if reg.Features.include(features...) {
			supportedRegistries = append(supportedRegistries, reg)
		}
	}
	return supportedRegistries
}

// Using does the actual work of filtering for a specific protocol type
func (d Drivers) Using(proto string) Drivers {
	supportedRegistries := make(Drivers, 0)
	for _, reg := range d {
		if reg.Protocol == proto {
			supportedRegistries = append(supportedRegistries, reg)
		}
	}
	return supportedRegistries
}

// For does the actual work of filtering for a specific driver name
func (d Drivers) For(driver string) Drivers {
	supportedRegistries := make(Drivers, 0)
	for _, reg := range d {
		if reg.Name == driver {
			supportedRegistries = append(supportedRegistries, reg)
		}
	}
	return supportedRegistries
}

// deduplicate returns a new slice with duplicates values removed.
func deduplicate(s []string) []string {
	if len(s) <= 1 {
		return s
	}
	result := []string{}
	seen := make(map[string]struct{})
	for _, val := range s {
		val := strings.ToLower(val)
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = struct{}{}
		}
	}
	return result
}

// PreferProtocol does the actual work of moving preferred protocols to the start of the driver registry
func (d Drivers) PreferProtocol(protocols ...string) Drivers {
	var final Drivers
	var leftOver Drivers
	tracking := make(map[int]Drivers)
	protocols = deduplicate(protocols)
	for _, registry := range d {
		var movedToTracking bool
		for index, pName := range protocols {
			if strings.EqualFold(registry.Protocol, pName) {
				tracking[index] = append(tracking[index], registry)
				movedToTracking = true
			}
		}
		if !movedToTracking {
			leftOver = append(leftOver, registry)
		}
	}
	for x := 0; x <= len(tracking); x++ {
		final = append(final, tracking[x]...)
	}
	final = append(final, leftOver...)
	return final
}
