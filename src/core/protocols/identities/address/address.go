package address

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"drm-blockchain/src/core/protocols/identities/identitykeys"
	"encoding/base64"
	"errors"
)

func ComputeAddressFromPublicKey(pubKey *ecdsa.PublicKey) string {
	bs := elliptic.MarshalCompressed(pubKey.Curve, pubKey.X, pubKey.Y)
	return base64.RawStdEncoding.EncodeToString(bs)
}

func ComputePublicKeyFromAddress(addr string) (*ecdsa.PublicKey, error) {
	pubKeyBs, err := base64.RawStdEncoding.DecodeString(addr)

	if err != nil {
		return nil, err
	}

	curve := identitykeys.GetCurve()
	x, y := elliptic.UnmarshalCompressed(curve, pubKeyBs)
	if x == nil {
		return nil, errors.New("failed public key unmarshalling")
	}

	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = curve
	pubKey.X = x
	pubKey.Y = y

	return pubKey, nil
}
