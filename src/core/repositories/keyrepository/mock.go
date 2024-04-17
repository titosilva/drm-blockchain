package keyrepository

import "drm-blockchain/src/core/protocols/identities"

type Mock struct {
	id *identities.Identity
}

// Initialize implements IKeyRepository.
func (*Mock) Initialize() error {
	panic("unimplemented")
}

// GetSelfIdentity implements IKeyRepository.
func (m *Mock) GetSelfIdentity() *identities.Identity {
	if m.id == nil {
		m.id = identities.Generate()
	}

	return m.id
}

// Sign implements IKeyRepository.
func (m *Mock) Sign(data []byte) ([]byte, error) {
	return m.id.Sign(data)
}

var _ IKeyRepository = new(Mock)
