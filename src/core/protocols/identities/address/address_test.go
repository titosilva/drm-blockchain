package address_test

import (
	"drm-blockchain/src/core/protocols/identities"
	ez "drm-blockchain/src/utils/test"
	"testing"
)

func Test__Address__PublicKey__Conversion(t *testing.T) {
	ez := ez.New(t)
	id := identities.Generate()

	addr := id.GetAddress()
	id2, err := identities.FromAddress(addr)
	ez.AssertNoError(err)

	ez.AssertAreEqual(id.GetAddress(), id2.GetAddress())
}
