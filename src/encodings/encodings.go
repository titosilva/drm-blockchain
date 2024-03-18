package encodings

type Encoding interface {
	Encode(any) ([]byte, error)
	Decode([]byte, any) error
}
