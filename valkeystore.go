package valkeystore

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/valkey-io/valkey-go"
	"io"
	"net/http"
	"strings"
)

const DefaultKeyPrefix = "session:"

// ValkeyStore stores gorilla sessions in Valkey
type ValkeyStore struct {
	client valkey.Client

	defaultSessionOptions sessions.Options

	keyPrefix string

	keyGen KeyGenFunc

	serializer Serializer
}

// KeyGenFunc is a function that generates a new session key
type KeyGenFunc func() (string, error)

// NewValkeyStore creates a new ValkeyStore with the given client and options.
func NewValkeyStore(client valkey.Client, options ...OptionsFunc) (*ValkeyStore, error) {
	vs := &ValkeyStore{
		defaultSessionOptions: sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		keyPrefix:  DefaultKeyPrefix,
		client:     client,
		keyGen:     generateRandomKey,
		serializer: NewGobSerializer(),
	}

	for _, option := range options {
		option(vs)
	}

	return vs, vs.client.Do(context.Background(), vs.client.B().Ping().Build()).Error()
}

// load loads a session from Valkey
func (s *ValkeyStore) load(ctx context.Context, session *sessions.Session) error {
	resp := s.client.Do(ctx, s.client.B().Get().Key(s.keyPrefix+session.ID).Build())
	if resp.Error() != nil {
		return resp.Error()
	}

	b, err := resp.AsBytes()
	if err != nil {
		return err
	}

	return s.serializer.Deserialize(b, session)
}

// Get returns a session for the given name after adding it to the registry.
func (s *ValkeyStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a new session with the given name without adding it to the registry.
func (s *ValkeyStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = &s.defaultSessionOptions
	session.IsNew = true

	c, err := r.Cookie(name)
	if err != nil {
		return session, nil
	}
	session.ID = c.Value

	err = s.load(r.Context(), session)
	if err != nil {
		return session, err
	} else if errors.Is(err, valkey.Nil) {
		err = nil
	}

	return session, err
}

// Save adds a single session to the response.
//
// If the Options.MaxAge of the session is <= 0 then the session file will be
// deleted from the store. With this process it enforces the properly
// session cookie handling so no need to trust in the cookie management in the
// web browser.
func (s *ValkeyStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Delete if max-age is reached
	if session.Options.MaxAge <= 0 {
		if err := s.delete(r.Context(), session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	}

	if session.ID == "" {
		id, err := s.keyGen()
		if err != nil {
			return errors.New("failed to generate session ID: " + err.Error())
		}
		session.ID = id
	}

	if err := s.save(r.Context(), session); err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), session.ID, session.Options))
	return nil
}

// Options sets the default session options for the store.
func (s *ValkeyStore) Options(options sessions.Options) {
	s.defaultSessionOptions = options
}

// KeyPrefix sets the key prefix for the store.
func (s *ValkeyStore) KeyPrefix(keyPrefix string) {
	s.keyPrefix = keyPrefix
}

// KeyGen sets the key generator for the store.
func (s *ValkeyStore) KeyGen(keyGen KeyGenFunc) {
	s.keyGen = keyGen
}

// Serializer sets the serializer for the store.
func (s *ValkeyStore) Serializer(ss Serializer) {
	s.serializer = ss
}

// Close closes the store
func (s *ValkeyStore) Close() {
	s.client.Close()
}

// delete deletes a session from Valkey
func (s *ValkeyStore) delete(ctx context.Context, session *sessions.Session) error {
	return s.client.Do(ctx, s.client.B().Del().Key(s.keyPrefix+session.ID).Build()).Error()
}

// save saves a session to Valkey
func (s *ValkeyStore) save(ctx context.Context, session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}

	// Save the session
	err = s.client.Do(ctx, s.client.B().Set().Key(s.keyPrefix+session.ID).Value(string(b)).ExSeconds(int64(session.Options.MaxAge)).Build()).Error()
	if err != nil {
		return err
	}
	return nil
}

// generateRandomKey returns a new random key
func generateRandomKey() (string, error) {
	k := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return "", err
	}
	return strings.TrimRight(base32.StdEncoding.EncodeToString(k), "="), nil
}
