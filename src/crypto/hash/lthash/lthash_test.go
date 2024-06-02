package lthash_test

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"drm-blockchain/src/crypto/hash/lthash"
	"drm-blockchain/src/math/uintp"
	ez "drm-blockchain/src/utils/test"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"
)

func EncryptMessage(key []byte, message []byte) ([]byte, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return cipherText, nil
}

func Test__LtHash__Should__EnableFileRecoveryEasily(t *testing.T) {
	file_block_size_bytes := 256
	encrypted_hash := lthash.NewDirect(500, 128, file_block_size_bytes, nil)

	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	aes_key := [32]byte{
		123, 23, 243, 123, 32, 89, 11, 33,
		123, 23, 243, 123, 32, 89, 11, 33,
		123, 23, 243, 123, 32, 89, 11, 33,
		123, 23, 243, 123, 32, 89, 11, 33,
	}

	encrypted, err_encryption := EncryptMessage(aes_key[:], bs)

	if err_encryption != nil {
		fmt.Println(err_encryption)
		return
	}

	encrypted_hash.ComputeDigest(encrypted)

	blocks_to_insert := 250
	nonce_hash := lthash.NewDirect(500, 128, file_block_size_bytes, nil)
	for i := 0; i < blocks_to_insert; i++ {
		nonces_and_position := make([]byte, file_block_size_bytes+8)
		_, err_nonces := rand.Read(nonces_and_position)
		if err_nonces != nil {
			fmt.Println(err)
			return
		}

		nonces := nonces_and_position[:file_block_size_bytes]
		position_bs := binary.BigEndian.Uint64(nonces_and_position[file_block_size_bytes:])
		count_of_blocks := uint64(len(encrypted)) / uint64(file_block_size_bytes)
		position := uint64(0)

		if count_of_blocks > 0 {
			position = (position_bs % uint64(count_of_blocks)) * uint64(file_block_size_bytes)
		}

		nonce_hash.Add(nonces)

		encrypted = append(encrypted[:position], append(nonces, encrypted[position:]...)...)
	}

	tampered_hash := lthash.NewDirect(500, 128, file_block_size_bytes, nil)
	tampered_hash.ComputeDigest(encrypted)

	original_hash_b64 := base64.StdEncoding.EncodeToString(encrypted_hash.GetDigest())

	remove := nonce_hash.GetState()
	tampered_hash.CombineInverse(remove)
	tampered_hash_b64 := base64.StdEncoding.EncodeToString(tampered_hash.GetDigest())
	if original_hash_b64 != tampered_hash_b64 {
		t.Error("Expected LtHash did not match")
	}
}

var test_vectors = []struct {
	chunk_count     uint
	chunk_size_bits uint
	m1              string
	m2              string
	bytes_to_add    []byte
}{
	{1, 64, "ffffffffffffffff", "cafe", []byte{0x01}},
	{1, 64, "ffffffffffffffff", "cafe", []byte{0x01, 0x02, 0xff, 0xdd}},
	{128, 64, "ffffffffffffffff", "cafe", []byte{0x01, 0x02, 0xff, 0xdd}},
	{512, 64, "ffffffffffffffff", "cafe", []byte{0x01, 0x02, 0xff, 0xdd, 0xfe, 0x45}},
	{1, 128,
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffffffffffffffffffffff",
		[]byte{0x01},
	},
	{512, 128,
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffffffffffffffffffffff",
		[]byte{0x01, 0x02, 0xff, 0xdd, 0xfe, 0x45},
	},
}

func Test__AddMul__Should__BeHomomorphic(t *testing.T) {
	for _, vec := range test_vectors {
		m1 := uintp.FromHex(uint64(vec.chunk_size_bits), vec.m1)
		m2 := uintp.FromHex(uint64(vec.chunk_size_bits), vec.m2)

		hash_mul := lthash.NewDirect(vec.chunk_count, vec.chunk_size_bits, 256, nil)
		hash_mul.AddMul(m1, vec.bytes_to_add)
		hash_mul.AddMul(m2, vec.bytes_to_add)

		m1.Add(m2)

		hash := lthash.NewDirect(vec.chunk_count, vec.chunk_size_bits, 256, nil)
		hash.AddMul(m1, vec.bytes_to_add)

		e := ez.New(t)
		e.AssertAreEqual(hash.GetDigest(), hash_mul.GetDigest())
	}
}

func Test__RemoveMul__Should__BeHomomorphic(t *testing.T) {
	for _, vec := range test_vectors {
		m1 := uintp.FromHex(uint64(vec.chunk_size_bits), vec.m1)
		m2 := uintp.FromHex(uint64(vec.chunk_size_bits), vec.m2)

		hash_mul := lthash.NewDirect(vec.chunk_count, vec.chunk_size_bits, 256, nil)
		hash_mul.AddMul(m1, vec.bytes_to_add)

		m1.Add(m2)

		hash := lthash.NewDirect(vec.chunk_count, vec.chunk_size_bits, 256, nil)
		hash.AddMul(m1, vec.bytes_to_add)

		hash.RemoveMul(m2, vec.bytes_to_add)

		e := ez.New(t)
		e.AssertAreEqual(hash.GetDigest(), hash_mul.GetDigest())
	}
}
