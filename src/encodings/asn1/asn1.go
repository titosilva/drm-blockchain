package asn1

import "encoding/asn1"

type ASN1Encoding struct {
}

func (enc ASN1Encoding) Encode(val any) ([]byte, error) {
	return asn1.Marshal(val)
}

func (enc ASN1Encoding) Decode(bs []byte, val any) error {
	_, err := asn1.Unmarshal(bs, val)
	return err
}
