package users

import (
	"rory-pearson/pkg/features"
	"sync"
)

var UsersFeatureType = features.FeatureType{ID: "users", Name: "Users"}

type UsersFeature struct {
	mu    sync.Mutex
	Users map[string]*User
}

func (f *UsersFeature) Initialize(c features.Config) error {
	f.Users = make(map[string]*User)

	return nil
}
