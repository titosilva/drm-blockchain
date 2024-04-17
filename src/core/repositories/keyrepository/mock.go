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

var _ IKeyRepository = new(Mock)
