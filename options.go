package valkeystore

import (
	"github.com/gorilla/sessions"
)

// OptionsFunc is a function that sets options on a ValkeyStore
type OptionsFunc = func(*ValkeyStore)

// WithSessionOptions sets the default session options for a ValkeyStore
func WithSessionOptions(options sessions.Options) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.Options(options)
	}
}

// WithKeyPrefix sets the key prefix for a ValkeyStore
func WithKeyPrefix(p string) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.KeyPrefix(p)
	}
}

// WithKeyGenFunc sets the key generation function for a ValkeyStore
func WithKeyGenFunc(kg KeyGenFunc) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.KeyGen(kg)
	}
}

// WithSerializer sets the serializer for a ValkeyStore
func WithSerializer(s Serializer) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.Serializer(s)
	}
}
