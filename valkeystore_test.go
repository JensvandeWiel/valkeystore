package valkeystore

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/valkey-io/valkey-go"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupValkey(ctx context.Context) (*redis.RedisContainer, error) {
	return redis.Run(ctx, "valkey/valkey:7.2.5")
}

func TestNew(t *testing.T) {

	ctx := context.Background()
	vk, err := setupValkey(ctx)
	if err != nil {
		t.Fatal("failed to setup valkey container", err)
	}

	defer func() {
		if err := vk.Terminate(ctx); err != nil {
			t.Fatal("failed to terminate container", err)
		}
	}()

	endpoint, err := vk.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:" + endpoint.Port()},
	})
	if err != nil {
		t.Fatal("failed to create valkey client", err)
	}

	defer client.Close()

	store, err := NewValkeyStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}
	if session.IsNew == false {
		t.Fatal("session is not new")
	}
}

func TestSessionOptions(t *testing.T) {
	ctx := context.Background()
	vk, err := setupValkey(ctx)
	if err != nil {
		t.Fatal("failed to setup valkey container", err)
	}

	defer func() {
		if err := vk.Terminate(ctx); err != nil {
			t.Fatal("failed to terminate container", err)
		}
	}()

	endpoint, err := vk.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:" + endpoint.Port()},
	})
	if err != nil {
		t.Fatal("failed to create valkey client", err)
	}

	defer client.Close()

	opts := sessions.Options{
		Path:   "/path",
		MaxAge: 99999,
	}

	store, err := NewValkeyStore(client, WithSessionOptions(opts))
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create store", err)
	}
	if session.Options.Path != opts.Path || session.Options.MaxAge != opts.MaxAge {
		t.Fatal("failed to set options")
	}
}

func TestSave(t *testing.T) {
	ctx := context.Background()
	vk, err := setupValkey(ctx)
	if err != nil {
		t.Fatal("failed to setup valkey container", err)
	}

	defer func() {
		if err := vk.Terminate(ctx); err != nil {
			t.Fatal("failed to terminate container", err)
		}
	}()

	endpoint, err := vk.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:" + endpoint.Port()},
	})
	if err != nil {
		t.Fatal("failed to create valkey client", err)
	}

	defer client.Close()

	store, err := NewValkeyStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := httptest.NewRecorder()

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	session.Values["key"] = "value"
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save: ", err)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	vk, err := setupValkey(ctx)
	if err != nil {
		t.Fatal("failed to setup valkey container", err)
	}

	defer func() {
		if err := vk.Terminate(ctx); err != nil {
			t.Fatal("failed to terminate container", err)
		}
	}()

	endpoint, err := vk.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:" + endpoint.Port()},
	})
	if err != nil {
		t.Fatal("failed to create valkey client", err)
	}

	defer client.Close()

	store, err := NewValkeyStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := httptest.NewRecorder()

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	session.Values["key"] = "value"
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save session: ", err)
	}

	session.Options.MaxAge = -1
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to delete session: ", err)
	}
}
