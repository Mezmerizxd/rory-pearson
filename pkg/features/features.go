// features.go
package features

import "rory-pearson/pkg/log"

// Config for feature initialization
type Config struct {
	Log log.Log
}

// FeatureType to identify the feature
type FeatureType struct {
	ID   string
	Name string
}

// Features struct to manage features
type Features struct {
	log             log.Log
	FeatureRegistry map[FeatureType]interface{}
}

var instance *Features

// Initialize initializes the Features instance.
func Initialize(cfg Config) *Features {
	feature := &Features{
		FeatureRegistry: make(map[FeatureType]interface{}), // Ensure the registry is initialized
	}

	instance = feature
	return feature
}

func GetInstance() *Features {
	if instance == nil {
		panic("Features instance not initialized")
	}
	return instance
}

// RegisterFeature registers a new feature.
func (f *Features) RegisterFeature(t FeatureType, feature interface{}) {
	f.FeatureRegistry[t] = feature // Assuming this is how you store the feature
	f.log.Info().Str("feature_id", t.ID).Str("feature_name", t.Name).Msg("registered feature")
}

// GetFeature retrieves a feature by its type.
func (f *Features) GetFeature(t FeatureType) (interface{}, bool) {
	feature, exists := f.FeatureRegistry[t]
	return feature, exists
}

// InitializeAll initializes all registered features.
func (f *Features) InitializeAll() error {
	for _, feature := range f.FeatureRegistry {
		if err := feature.(interface{ Initialize(Config) error }).Initialize(Config{}); err != nil {
			return err
		}
	}
	f.log.Info().Msg("all features initialized")
	return nil
}
