package registrar

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	FeaturePowerSet   Feature = "powerset"
	FeatureUserCreate Feature = "usercreate"
)

type driverOne struct {
	name         string
	protocol     string
	features     Features
	isCompatible bool
}

func (do *driverOne) Compatible(ctx context.Context) bool {
	return do.isCompatible
}

func TestInclude(t *testing.T) {
	testCases := []struct {
		name     string
		features Features
		includes Features
		want     bool
	}{
		{name: "feature is not included 1", features: Features{}, includes: Features{FeaturePowerSet}, want: false},
		{name: "feature is not included 2", features: Features{FeatureUserCreate}, includes: Features{FeaturePowerSet}, want: false},
		{name: "feature included", features: Features{FeaturePowerSet}, includes: Features{FeaturePowerSet}, want: true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.features.include(tc.includes...)
			diff := cmp.Diff(tc.want, result)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSupports(t *testing.T) {
	testCases := []struct {
		name       string
		collection *Registry
		supports   Features
		want       *Registry
	}{
		{
			name: "no registry supports UserCreate",
			collection: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}})),
			supports: Features{
				FeatureUserCreate,
			},
			want: NewRegistry(WithDrivers(Drivers{})),
		},
		{
			name: "one registry supports UserCreate",
			collection: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeatureUserCreate},
				}})),
			supports: Features{
				FeatureUserCreate,
			},
			want: NewRegistry(WithDrivers(Drivers{{
				Name:     "ipmitool",
				Protocol: "ipmi",
				Features: []Feature{FeatureUserCreate},
			}})),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.collection.Supports(tc.supports...)
			diff := cmp.Diff(tc.want.Drivers, result)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUsing(t *testing.T) {
	testCases := []struct {
		name       string
		collection *Registry
		proto      string
		want       *Registry
	}{
		{
			name: "proto is not found",
			collection: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}})),
			proto: "web",
			want:  NewRegistry(WithDrivers(Drivers{})),
		},
		{
			name: "proto is found",
			collection: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}})),
			proto: "ipmi",
			want: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}})),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.collection.Using(tc.proto)
			diff := cmp.Diff(tc.want.Drivers, result)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFor(t *testing.T) {
	testCases := []struct {
		name       string
		collection *Registry
		provider   string
		want       *Registry
	}{
		{
			name: "proto is not found",
			collection: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}})),
			provider: "dell",
			want:     NewRegistry(WithDrivers(Drivers{})),
		},
		{
			name: "proto is found",
			collection: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}})),
			provider: "ipmitool",
			want: NewRegistry(WithDrivers(Drivers{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}})),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.collection.For(tc.provider)
			diff := cmp.Diff(tc.want.Drivers, result)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestAll(t *testing.T) {
	testCases := []struct {
		name         string
		addARegistry bool
		want         *Registry
	}{
		{name: "empty collection", want: NewRegistry(WithDrivers(Drivers{}))},
		{name: "single collection", addARegistry: true, want: NewRegistry(WithDrivers(Drivers{
			{Name: "dell", Protocol: "web", Features: []Feature{FeaturePowerSet}},
		})),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", []Feature{FeaturePowerSet}, nil, nil)
			}
			if diff := cmp.Diff(tc.want, rg, cmpopts.IgnoreFields(Registry{}, "Logger")); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSupportFn(t *testing.T) {
	testCases := []struct {
		name         string
		addARegistry bool
		features     []Feature
		want         *Registry
	}{
		{name: "empty collection", want: NewRegistry(WithDrivers(Drivers{}))},
		{name: "single collection", features: []Feature{FeatureUserCreate}, addARegistry: true, want: NewRegistry(WithDrivers(Drivers{
			{Name: "dell", Protocol: "web", Features: []Feature{FeatureUserCreate}},
		})),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", []Feature{FeatureUserCreate}, nil, nil)
			}
			result := rg.Supports(tc.features...)
			if diff := cmp.Diff(tc.want.Drivers, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUsingFn(t *testing.T) {
	testCases := []struct {
		name         string
		addARegistry bool
		proto        string
		want         *Registry
	}{
		{name: "empty collection", want: NewRegistry(WithDrivers(Drivers{}))},
		{name: "single collection", proto: "web", addARegistry: true, want: NewRegistry(WithDrivers(Drivers{{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		}}))},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", []Feature{FeatureUserCreate}, nil, nil)
			}
			result := rg.Using(tc.proto)
			if diff := cmp.Diff(tc.want.Drivers, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestForFn(t *testing.T) {
	testCases := []struct {
		name         string
		addARegistry bool
		provider     string
		want         *Registry
	}{
		{name: "empty collection", want: NewRegistry(WithDrivers(Drivers{}))},
		{name: "single collection", provider: "dell", addARegistry: true, want: NewRegistry(WithDrivers(Drivers{{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		}}))},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", []Feature{FeatureUserCreate}, nil, nil)
			}
			result := rg.For(tc.provider)
			if diff := cmp.Diff(tc.want.Drivers, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPreferProtocol(t *testing.T) {
	collection1 := NewRegistry(WithDrivers(Drivers{
		{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		},
		{
			Name:     "ipmitool",
			Protocol: "ipmi",
			Features: []Feature{FeatureUserCreate},
		},
		{
			Name:     "smc",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		},
	}))
	testCases := map[string]struct {
		baseCollection *Registry
		protocol       []string
		want           *Registry
	}{
		"empty collection": {want: collection1, baseCollection: collection1},
		"collection": {protocol: []string{"web"}, baseCollection: collection1, want: NewRegistry(WithDrivers(Drivers{
			{
				Name:     "dell",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "smc",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "ipmitool",
				Protocol: "ipmi",
				Features: []Feature{FeatureUserCreate},
			},
		}))},
		"collection with duplicate protocols": {protocol: []string{"web", "web"}, baseCollection: collection1, want: NewRegistry(WithDrivers(Drivers{
			{
				Name:     "dell",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "smc",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "ipmitool",
				Protocol: "ipmi",
				Features: []Feature{FeatureUserCreate},
			},
		}))},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.baseCollection.PreferProtocol(tc.protocol...)
			if diff := cmp.Diff(tc.want.Drivers, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPreferDriver(t *testing.T) {
	collection1 := NewRegistry(WithDrivers(Drivers{
		{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		},
		{
			Name:     "ipmitool",
			Protocol: "ipmi",
			Features: []Feature{FeatureUserCreate},
		},
		{
			Name:     "smc",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		},
	}))
	collection2 := NewRegistry(WithDrivers(Drivers{
		{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		},
		{
			Name:     "ipmitool",
			Protocol: "ipmi",
			Features: []Feature{FeatureUserCreate},
		},
		{
			Name:     "smc",
			Protocol: "redfish",
			Features: []Feature{FeatureUserCreate},
		},
		{
			Name:     "smc",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		},
	}))
	testCases := map[string]struct {
		baseCollection *Registry
		driver         []string
		want           *Registry
	}{
		"empty collection": {want: collection1, baseCollection: collection1},
		"normal collection": {driver: []string{"smc"}, baseCollection: collection1, want: NewRegistry(WithDrivers(Drivers{
			{
				Name:     "smc",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "dell",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "ipmitool",
				Protocol: "ipmi",
				Features: []Feature{FeatureUserCreate},
			},
		}))},
		"collection with duplicate drivers": {driver: []string{"smc", "smc"}, baseCollection: collection2, want: NewRegistry(WithDrivers(Drivers{
			{
				Name:     "smc",
				Protocol: "redfish",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "smc",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "dell",
				Protocol: "web",
				Features: []Feature{FeatureUserCreate},
			},
			{
				Name:     "ipmitool",
				Protocol: "ipmi",
				Features: []Feature{FeatureUserCreate},
			},
		}))},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.baseCollection.PreferDriver(tc.driver...)
			if diff := cmp.Diff(tc.want.Drivers, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestGetDriverInterfaces(t *testing.T) {
	rg := NewRegistry()
	do := &driverOne{}
	rg.Register(do.name, do.protocol, do.features, nil, do)
	driverInterfaces := rg.GetDriverInterfaces()
	if diff := cmp.Diff(driverInterfaces, []interface{}{do}, cmp.AllowUnexported(driverOne{})); diff != "" {
		t.Fatal(diff)
	}
}

func TestFilterForCompatible(t *testing.T) {
	compatible := &driverOne{name: "driverOne", protocol: "tcp", features: Features{}, isCompatible: true}
	notCompatible := &driverOne{name: "driverOne", protocol: "tcp", features: Features{}, isCompatible: false}
	testCases := []struct {
		name   string
		driver *driverOne
		want   []interface{}
	}{
		{name: "is compatible", driver: compatible, want: []interface{}{compatible}},
		{name: "is NOT compatible", driver: notCompatible, want: []interface{}{}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rg := NewRegistry(WithLogger(defaultLogger()))
			rg.Register(tc.driver.name, tc.driver.protocol, tc.driver.features, nil, tc.driver)
			rg.Drivers = rg.FilterForCompatible(context.Background())
			driverInterfaces := rg.GetDriverInterfaces()
			if diff := cmp.Diff(driverInterfaces, tc.want, cmp.AllowUnexported(driverOne{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
