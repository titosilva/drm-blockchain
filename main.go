package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"drm-blockchain/src/crypto/hash/lthash"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

func main() {
	file_block_size_bytes := 256
	encrypted_hash := lthash.New(500, 110, file_block_size_bytes, nil)

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
	nonce_hash := lthash.New(500, 110, file_block_size_bytes, nil)
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

	fmt.Println("Hash of nonces:")
	fmt.Println(base64.StdEncoding.EncodeToString(nonce_hash.GetDigest()))

	tampered_hash := lthash.New(500, 110, file_block_size_bytes, nil)
	tampered_hash.ComputeDigest(encrypted)
	fmt.Println("Hash of tampered file:")
	fmt.Println(base64.StdEncoding.EncodeToString(tampered_hash.GetDigest()))

	fmt.Println("Hash of original encrypted file:")
	fmt.Println(base64.StdEncoding.EncodeToString(encrypted_hash.GetDigest()))

	fmt.Println("Computed hash of untampered file:")
	remove := nonce_hash.GetState().Invert()
	tampered_hash.Combine(remove)
	fmt.Println(base64.StdEncoding.EncodeToString(tampered_hash.GetDigest()))
}
