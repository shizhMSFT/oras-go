package storage

import (
	"bytes"
	"context"
	"io"
	"sync"

	"oras.land/oras-go/dag"
	"oras.land/oras-go/errdef"
	"oras.land/oras-go/internal/ioutil"
)

// Memory is a memory based CAS.
type Memory struct {
	content sync.Map
}

// NewMemory creates a new Memory CAS.
func NewMemory() *Memory {
	return &Memory{}
}

// Fetch fetches the content identified by the descriptor.
func (m *Memory) Fetch(_ context.Context, target dag.Descriptor) (io.ReadCloser, error) {
	content, exists := m.content.Load(target.Digest)
	if !exists {
		return nil, errdef.ErrNotFound
	}
	return ioutil.NopCloser(bytes.NewReader(content.([]byte))), nil
}

// Push pushes the content, matching the expected descriptor.
func (m *Memory) Push(_ context.Context, expected dag.Descriptor, content io.Reader) error {
	value, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	if _, exists := m.content.LoadOrStore(expected.Digest, value); exists {
		return errdef.ErrAlreadyExists
	}
	return nil
}

// Exists returns true if the described content exists.
func (m *Memory) Exists(_ context.Context, target dag.Descriptor) (bool, error) {
	_, exists := m.content.Load(target.Digest)
	return exists, nil
}
