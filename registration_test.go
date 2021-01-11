package registration

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	// FeaturePowerState represents the powerstate functionality
	// an implementation will use these when they have implemented
	// corresponding interface method.
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
			result := tc.features.Include(tc.includes...)
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
		collection Registry
		supports   Features
		want       Registry
	}{
		{
			name: "no registry supports UserCreate",
			collection: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}},
			supports: Features{
				FeatureUserCreate,
			},
			want: Registry{},
		},
		{
			name: "one registry supports UserCreate",
			collection: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeatureUserCreate},
				}},
			supports: Features{
				FeatureUserCreate,
			},
			want: Registry{{
				Name:     "ipmitool",
				Protocol: "ipmi",
				Features: []Feature{FeatureUserCreate},
			}},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.collection.Supports(tc.supports...)
			diff := cmp.Diff(tc.want, result)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUsing(t *testing.T) {
	testCases := []struct {
		name       string
		collection Registry
		proto      string
		want       Registry
	}{
		{
			name: "proto is not found",
			collection: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}},
			proto: "web",
			want:  Registry{},
		},
		{
			name: "proto is found",
			collection: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}},
			proto: "ipmi",
			want: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.collection.Using(tc.proto)
			diff := cmp.Diff(tc.want, result)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFor(t *testing.T) {
	testCases := []struct {
		name       string
		collection Registry
		provider   string
		want       Registry
	}{
		{
			name: "proto is not found",
			collection: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}},
			provider: "dell",
			want:     Registry{},
		},
		{
			name: "proto is found",
			collection: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}},
			provider: "ipmitool",
			want: Registry{
				{
					Name:     "ipmitool",
					Protocol: "ipmi",
					Features: []Feature{FeaturePowerSet},
				}},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.collection.For(tc.provider)
			diff := cmp.Diff(tc.want, result)
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
		want         Registry
	}{
		{name: "empty collection", want: Registry{}},
		{name: "single collection", addARegistry: true, want: Registry{{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeaturePowerSet},
		}}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// register
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", nil, nil, []Feature{FeaturePowerSet})
				t.Log(rg)
			}
			if diff := cmp.Diff(tc.want, rg); diff != "" {
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
		want         Registry
	}{
		{name: "empty collection", want: Registry{}},
		{name: "single collection", features: []Feature{FeatureUserCreate}, addARegistry: true, want: Registry{{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		}}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// register
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", nil, nil, []Feature{FeatureUserCreate})
			}
			result := rg.Supports(tc.features...)
			if diff := cmp.Diff(tc.want, result); diff != "" {
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
		want         Registry
	}{
		{name: "empty collection", want: Registry{}},
		{name: "single collection", proto: "web", addARegistry: true, want: Registry{{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		}}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// register
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", nil, nil, []Feature{FeatureUserCreate})
			}
			result := rg.Using(tc.proto)
			if diff := cmp.Diff(tc.want, result); diff != "" {
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
		want         Registry
	}{
		{name: "empty collection", want: Registry{}},
		{name: "single collection", provider: "dell", addARegistry: true, want: Registry{{
			Name:     "dell",
			Protocol: "web",
			Features: []Feature{FeatureUserCreate},
		}}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// register
			rg := NewRegistry()
			if tc.addARegistry {
				rg.Register("dell", "web", nil, nil, []Feature{FeatureUserCreate})
			}
			result := rg.For(tc.provider)
			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPrefer(t *testing.T) {
	unorderedCollection := Registry{
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
	}
	testCases := []struct {
		name         string
		addARegistry bool
		protocol     []string
		want         Registry
	}{
		{name: "empty collection", want: unorderedCollection},
		{name: "collection", protocol: []string{"web"}, addARegistry: true, want: Registry{
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
		}},
		{name: "collection with duplicate protocols", protocol: []string{"web", "web"}, addARegistry: true, want: Registry{
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
		}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// register
			registries := unorderedCollection
			result := registries.PreferProtocol(tc.protocol...)
			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestGetDriverInterfaces(t *testing.T) {
	rg := NewRegistry()
	do := &driverOne{}
	rg.Register(do.name, do.protocol, do, do, do.features)
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
			rg := NewRegistry()
			rg.Register(tc.driver.name, tc.driver.protocol, tc.driver, tc.driver, tc.driver.features)
			rg.FilterForCompatible(context.Background())
			driverInterfaces := rg.GetDriverInterfaces()
			if diff := cmp.Diff(driverInterfaces, tc.want, cmp.AllowUnexported(driverOne{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
