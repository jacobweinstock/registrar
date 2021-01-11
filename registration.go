package registry

import (
	"context"
	"strings"
	"sync"
)

// Features holds the features a provider supports
type Features []Feature

// Feature represents a single feature available
type Feature string

// Registry holds a slice of Registry types
type Registry []*Driver

// Verifier for whether the driver is compatible
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

// NewRegistry new Collection
func NewRegistry() Registry {
	return make(Registry, 0)
}

// Register will add a provider with details to the main registryCollection
func (r *Registry) Register(name, protocol string, driverInterface interface{}, compatFn Verifier, features Features) {
	*r = append(*r, &Driver{
		Name:            name,
		Protocol:        protocol,
		Features:        features,
		Verifier:        compatFn,
		DriverInterface: driverInterface,
	})
}

// GetDriverInterfaces returns a slice of just the driver interfaces
func (r Registry) GetDriverInterfaces() []interface{} {
	results := make([]interface{}, 0)
	for _, elem := range r {
		results = append(results, elem.DriverInterface)
	}
	return results
}

// FilterForCompatible updates the registry with only compatible implementations
func (r *Registry) FilterForCompatible(ctx context.Context) {
	var wg sync.WaitGroup
	result := make(Registry, 0)
	for _, elem := range *r {
		wg.Add(1)
		go func(isCompat Verifier, reg *Driver, wg *sync.WaitGroup) {
			if isCompat.Compatible(ctx) {
				result = append(result, reg)
			}
			wg.Done()
		}(elem.Verifier, elem, &wg)
	}
	wg.Wait()
	*r = result
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
func (r Registry) Supports(features ...Feature) Registry {
	supportedRegistries := make(Registry, 0)
	for _, reg := range r {
		if reg.Features.include(features...) {
			supportedRegistries = append(supportedRegistries, reg)
		}
	}
	return supportedRegistries
}

// Using does the actual work of filtering for a specific protocol type
func (r Registry) Using(proto string) Registry {
	supportedRegistries := make(Registry, 0)
	for _, reg := range r {
		if reg.Protocol == proto {
			supportedRegistries = append(supportedRegistries, reg)
		}
	}
	return supportedRegistries
}

// For does the actual work of filtering for a specific provider name
func (r Registry) For(provider string) Registry {
	supportedRegistries := make(Registry, 0)
	for _, reg := range r {
		if reg.Name == provider {
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

// PreferProtocol does the actual work of moving preferred protocols to the start of the collection
func (r Registry) PreferProtocol(protocols ...string) Registry {
	var final Registry
	var leftOver Registry
	tracking := make(map[int]Registry)
	protocols = deduplicate(protocols)
	for _, registry := range r {
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
