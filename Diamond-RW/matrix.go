package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	ivSize       = aes.BlockSize
	iterationNum = 4096
	chunkSize    = 10 * 1024 * 1024 // 10 MB chunk size
)

type SpaceEncryptor struct {
	password []byte
	Salt     []byte
	keySize  int
}

func reverseBytes(data []byte) []byte {
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-1-i] = data[len(data)-1-i], data[i]
	}

	return data
}

func NewSpaceEncryptor(password string, salt string) *SpaceEncryptor {
	return &SpaceEncryptor{[]byte(password), []byte(salt), len(password)}
}

func (e *SpaceEncryptor) EncryptMessageV2(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	outputFile, err := os.OpenFile(filename+".diamond", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	key := pbkdf2.Key(e.password, e.Salt, iterationNum, e.keySize, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, ivSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	// Write Salt and IV to the beginning of the output file
	if _, err := writer.Write(e.Salt); err != nil {
		return err
	}
	if _, err := writer.Write(iv); err != nil {
		return err
	}

	// Create a new CTR stream for each chunk
	stream := cipher.NewCTR(block, iv)

	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		ciphertext := make([]byte, n)
		stream.XORKeyStream(ciphertext, buffer[:n])

		if _, err := writer.Write(ciphertext); err != nil {
			return err
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
