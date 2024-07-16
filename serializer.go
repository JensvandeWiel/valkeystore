package valkeystore

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
)

// Serializer is an interface for serializing and deserializing sessions for Valkey
type Serializer interface {
	Serialize(s *sessions.Session) ([]byte, error)
	Deserialize(b []byte, s *sessions.Session) error
}

// GobSerializer is a Serializer that uses gob encoding
type GobSerializer struct{}

// Serialize serializes a session using gob encoding
func (s *GobSerializer) Serialize(session *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(session.Values)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize deserializes a session using gob encoding
func (s *GobSerializer) Deserialize(b []byte, session *sessions.Session) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&session.Values)
}

// NewGobSerializer creates a new GobSerializer
func NewGobSerializer() *GobSerializer {
	return &GobSerializer{}
}

// JSONSerializer is a Serializer that uses JSON encoding
type JSONSerializer struct{}

// NewJSONSerializer creates a new JSONSerializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

// Serialize serializes a session using JSON encoding
func (s *JSONSerializer) Serialize(session *sessions.Session) ([]byte, error) {
	m := make(map[string]interface{}, len(session.Values))
	for k, v := range session.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("non-string key value, cannot serialize session to JSON: %v", k)
			fmt.Printf("redistore.JSONSerializer.serialize() Error: %v", err)
			return nil, err
		}
		m[ks] = v
	}
	return json.Marshal(m)
}

// Deserialize deserializes a session using JSON encoding
func (s *JSONSerializer) Deserialize(b []byte, session *sessions.Session) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Printf("redistore.JSONSerializer.deserialize() Error: %v", err)
		return err
	}
	for k, v := range m {
		session.Values[k] = v
	}
	return nil
}
