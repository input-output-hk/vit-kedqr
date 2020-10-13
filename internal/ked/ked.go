package ked

import (
	"crypto/rand"
	"crypto/sha512"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/pbkdf2"
)

/*
	----------------------------------------------------------
	| 0x01 | SALT(16) | NONCE(12) | Encrypted Data | Tag(16) |
	----------------------------------------------------------
*/

const (
	PROTO_SIZE = 1
	SALT_SIZE  = 16
	NONCE_SIZE = 12 // chacha20poly1305.NonceSize
	TAG_SIZE   = 16 // poly1305.TagSize
)

var (
	PASSWORD_DERIVATION_ITERATIONS = 12_983
	PROTO_VERSION                  = byte(0b0000_0001)
)

// Encrypt ...
func Encrypt(password []byte, data []byte) ([]byte, error) {
	salt, err := generate_salt()
	if err != nil {
		return nil, err
	}
	nonce, err := generate_nonce()
	if err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.New(
		passToKey(password, salt),
	)
	if err != nil {
		return nil, err
	}

	dataLen := len(data)
	encFull := make([]byte, 0, PROTO_SIZE+SALT_SIZE+NONCE_SIZE+dataLen+TAG_SIZE)

	encFull = append(encFull, PROTO_VERSION) // | 0x01 |
	encFull = append(encFull, salt...)       // | 0x01 | SALT(16) |
	encFull = append(encFull, nonce...)      // | 0x01 | SALT(16) | NONCE(12) |

	return aead.Seal(
		encFull, // | 0x01 | SALT(16) | NONCE(12) | Encrypted Data | Tag(16) |
		nonce,
		data,
		nil,
	), nil
}

// Decrypt ...
func Decrypt(password []byte, data []byte) ([]byte, error) {
	var decData []byte

	aead, err := chacha20poly1305.New(
		passToKey(password, data[PROTO_SIZE:PROTO_SIZE+SALT_SIZE]),
	)
	if err != nil {
		return nil, err
	}

	return aead.Open(
		decData, // dec data
		data[PROTO_SIZE+SALT_SIZE:PROTO_SIZE+SALT_SIZE+NONCE_SIZE], // nonce
		data[PROTO_SIZE+SALT_SIZE+NONCE_SIZE:],                     // cipher (enc data + tag)
		nil,                                                        // aad
	)
}

func passToKey(password []byte, salt []byte) []byte {
	return pbkdf2.Key(
		password,
		salt,
		PASSWORD_DERIVATION_ITERATIONS,
		chacha20poly1305.KeySize,
		sha512.New,
	)
}

func generate_salt() ([]byte, error) {
	salt := make([]byte, SALT_SIZE)
	err := randBytes(&salt)

	return salt, err
}

func generate_nonce() ([]byte, error) {
	nonce := make([]byte, NONCE_SIZE)
	err := randBytes(&nonce)

	return nonce, err
}

func randBytes(d *[]byte) error {
	_, err := rand.Read(*d)

	return err
}
