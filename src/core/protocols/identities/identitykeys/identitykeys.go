package identitykeys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	errorutils "drm-blockchain/src/utils/error"
	"errors"
)

func GetCurve() elliptic.Curve {
	return elliptic.P256()
}

func GeneratePrivateKey() *ecdsa.PrivateKey {
	priv, _ := ecdsa.GenerateKey(GetCurve(), rand.Reader)
	return priv
}

func DecodeIdentityPrivateKey(bs []byte) (*ecdsa.PrivateKey, error) {
	privKeyGen, err := x509.ParsePKCS8PrivateKey(bs)
	if err != nil {
		return nil, err
	}

	privKey, ok := privKeyGen.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("could not convert received key to ecdsa private key")
	}

	if privKey.Curve != GetCurve() {
		return nil, errorutils.Newf("expected key with curve %s, but got curve %s", GetCurve().Params().Name, privKey.Params().Name)
	}

	return privKey, nil
}

func EncodeIdentityPrivateKey(privKey *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalPKCS8PrivateKey(privKey)
}
