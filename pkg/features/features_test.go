// features_test.go
package features_test

import (
	"rory-pearson/pkg/features"
	"testing"
)

// TestInitialize verifies that the Features instance initializes correctly.
func TestInitialize(t *testing.T) {
	cfg := features.Config{}

	f := features.Initialize(cfg)
	if f == nil {
		t.Fatal("expected Features instance to be initialized")
	}
}

// FeatureType to identify the feature
var Type = features.FeatureType{ID: "feature_one", Name: "Feature One"}

// FeatureOne represents the first feature
type FeatureOne struct{}

// Initialize initializes the feature with the provided configuration.
func (f *FeatureOne) Initialize(c features.Config) error {
	// Perform any setup tasks here
	return nil // Return nil to indicate success
}

func (f *FeatureOne) DoSomething() {
	// Perform some action
}

// TestRegisterAndInitializeFeature verifies that features can be registered and initialized correctly.
func TestRegisterAndInitializeFeature(t *testing.T) {
	cfg := features.Config{}
	f := features.Initialize(cfg)

	if f == nil {
		t.Fatal("expected Features instance to be initialized")
	}

	f.RegisterFeature(Type, &FeatureOne{})

	if len(f.FeatureRegistry) != 1 {
		t.Fatalf("expected 1 feature to be registered, got %d", len(f.FeatureRegistry))
	}

	if err := f.InitializeAll(); err != nil {
		t.Fatalf("expected no error during feature initialization, got %v", err)
	}

	feature, exists := f.GetFeature(Type)
	if !exists {
		t.Fatal("expected FeatureOne to be registered")
	}

	featureOne, ok := feature.(*FeatureOne)
	if !ok {
		t.Fatal("expected feature to be of type *FeatureOne")
	}

	featureOne.DoSomething()

	if feature == nil {
		t.Fatal("expected FeatureOne instance to be initialized")
	}
}
