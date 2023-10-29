package main

import (
	"bufio"
	"drm-blockchain/src/crypto/hash/lthash"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func main() {
	hash := lthash.New(500, 110, 2048, nil)

	file, err := os.Open("webpki-setup-64.deb")
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

	hash.ComputeDigest(bs)
	fmt.Println("Hash:")
	fmt.Println(base64.StdEncoding.EncodeToString(hash.GetDigest()))
}
