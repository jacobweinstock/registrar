package registration

import (
	"context"
	"strings"
)

// Features holds the features a provider supports
type Features []Feature

// Feature represents a single feature available
type Feature string

// Registry holds a slice of Registry types
type Registry []*Driver

// Initializer inits a driver
type Initializer interface {
	Init() interface{}
}

// Compatibler for whether the driver is compatible
type Compatibler interface {
	Compatible(context.Context) bool
}

// Driver holds the info about a driver
type Driver struct {
	Name              string
	Protocol          string
	Features          Features
	RegistryInterface interface{}
	Initializer
	Compatibler
}

// NewRegistry new Collection
func NewRegistry() Registry {
	return make(Registry, 0)
}

// Include does the actual work of filtering for specific features
func (rf Features) Include(features ...Feature) bool {
	if len(features) > len(rf) {
		return false
	}
	fKeys := make(map[Feature]bool)
	for _, v := range rf {
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
func (rc Registry) Supports(features ...Feature) Registry {
	supportedRegistries := make(Registry, 0)
	for _, reg := range rc {
		if reg.Features.Include(features...) {
			supportedRegistries = append(supportedRegistries, reg)
		}
	}
	return supportedRegistries
}

// Using does the actual work of filtering for a specific protocol type
func (rc Registry) Using(proto string) Registry {
	supportedRegistries := make(Registry, 0)
	for _, reg := range rc {
		if reg.Protocol == proto {
			supportedRegistries = append(supportedRegistries, reg)
		}
	}
	return supportedRegistries
}

// For does the actual work of filtering for a specific provider name
func (rc Registry) For(provider string) Registry {
	supportedRegistries := make(Registry, 0)
	for _, reg := range rc {
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
func (rc Registry) PreferProtocol(protocols ...string) Registry {
	var final Registry
	var leftOver Registry
	tracking := make(map[int]Registry)
	protocols = deduplicate(protocols)
	for _, registry := range rc {
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

// Register will add a provider with details to the main registryCollection
func (rc *Registry) Register(provider, protocol string, initfn Initializer, compatFn Compatibler, features Features) {
	*rc = append(*rc, &Driver{
		Name:              provider,
		Protocol:          protocol,
		Initializer:       initfn,
		Features:          features,
		Compatibler:       compatFn,
		RegistryInterface: initfn.Init(),
	})
}
