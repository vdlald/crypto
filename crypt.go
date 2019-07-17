package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/argon2"
)

type Cryptor struct {
	Aead  cipher.AEAD
	Nonce []byte
}

func (c *Cryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return c.Aead.Open(nil, c.Nonce, ciphertext, nil)
}

func (c *Cryptor) Encrypt(plaintext []byte) []byte {
	return c.Aead.Seal(nil, c.Nonce, plaintext, nil)
}

func NewCryptor(passphrase string, salt string) (*Cryptor, error) {
	c := new(Cryptor)

	kdf := argon2.Key([]byte(passphrase), []byte(salt), 4, 32*1024, 4, 44)
	c.Nonce = kdf[32:]

	block, err := aes.NewCipher(kdf[:32])
	if err != nil {
		return c, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return c, err
	}

	c.Aead = aead

	return c, nil
}

func check(e error) {
	if e != nil {
		log.Printf("Error: %s\n", e.Error())
		os.Exit(0)
	}
}

func Crypt(action string, ifile string, password string, salt string) []byte {
	var idata []byte
	var odata []byte

	if salt == "" {
		salt = "AKatmtgdkMKq5SFYLt8tBlUxuwLccdCjFfFNi2b3o9A"
	}

	cryptor, err := NewCryptor(password, salt)
	check(err)

	idata, err = ioutil.ReadFile(ifile)
	check(err)

	switch action {
	case "encrypt":
		odata = cryptor.Encrypt(idata)
	case "decrypt":
		odata, err = cryptor.Decrypt(idata)
		check(err)
	}
	return odata
}
