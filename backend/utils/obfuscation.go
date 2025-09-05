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
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write([]byte(token)); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	compressed := buf.Bytes()

	xored := xorData(compressed)

	return base64.RawURLEncoding.EncodeToString(xored), nil
}

func DecodeToken(encoded string) (string, error) {

	xored, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	compressed := xorData(xored)

	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		return "", err
	}
	return out.String(), nil
}
