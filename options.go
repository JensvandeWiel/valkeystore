package valkeystore

import (
	"github.com/gorilla/sessions"
)

type OptionsFunc = func(*ValkeyStore)

func WithSessionOptions(options sessions.Options) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.Options(options)
	}
}

func WithKeyPrefix(p string) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.KeyPrefix(p)
	}
}

func WithKeyGenFunc(kg KeyGenFunc) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.KeyGen(kg)
	}
}

func WithSerializer(s Serializer) OptionsFunc {
	return func(vs *ValkeyStore) {
		vs.Serializer(s)
	}
}
