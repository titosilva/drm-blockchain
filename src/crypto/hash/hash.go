package hash

type HashAlgorithm interface {
	ComputeDigest(bytes []byte)
	GetDigest() []byte
}

type HomHashAlgorithm interface {
	HashAlgorithm
	Add(bytes []byte)
	Remove(bytes []byte)
}
