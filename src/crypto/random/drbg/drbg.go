package drbg

type DRBG interface {
	Seed(seed []byte) error
	GenerateByte() (byte, error)
}
