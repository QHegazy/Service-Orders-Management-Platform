package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
	"log"
	"os"
)

func xorData(data []byte) []byte {

	key := os.Getenv("XOR_KEY")

	if key == "" {
		log.Fatal("XOR_KEY is not set")
		key = "default_key"
	}

	res := make([]byte, len(data))
	keyLen := len(key)
	for i, b := range data {
		res[i] = b ^ key[i%keyLen]
	}
	return res
}

func EncodeToken(token string) (string, error) {
	log.Printf("Obfuscation - Encoding token (length: %d)", len(token))
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write([]byte(token)); err != nil {
		log.Printf("Obfuscation - Failed to write token for compression: %v", err)
		return "", err
	}
	if err := w.Close(); err != nil {
		log.Printf("Obfuscation - Failed to close compression writer: %v", err)
		return "", err
	}
	compressed := buf.Bytes()
	log.Printf("Obfuscation - Token compressed from %d to %d bytes", len(token), len(compressed))

	xored := xorData(compressed)
	log.Printf("Obfuscation - Token XOR applied, final size: %d bytes", len(xored))

	encoded := base64.RawURLEncoding.EncodeToString(xored)
	log.Printf("Obfuscation - Token encoded successfully (final length: %d)", len(encoded))
	return encoded, nil
}

func DecodeToken(encoded string) (string, error) {
	log.Printf("Obfuscation - Decoding token (encoded length: %d)", len(encoded))

	xored, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		log.Printf("Obfuscation - Failed to decode base64: %v", err)
		return "", err
	}
	log.Printf("Obfuscation - Base64 decoded, size: %d bytes", len(xored))

	compressed := xorData(xored)
	log.Printf("Obfuscation - XOR reversed, compressed size: %d bytes", len(compressed))

	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		log.Printf("Obfuscation - Failed to create zlib reader: %v", err)
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		log.Printf("Obfuscation - Failed to decompress: %v", err)
		return "", err
	}

	decoded := out.String()
	log.Printf("Obfuscation - Token decoded successfully (final length: %d)", len(decoded))
	return decoded, nil
}
