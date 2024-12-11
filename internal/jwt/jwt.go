package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type GoogleCerts struct {
	Keys []struct {
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

func KeyFunc(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New("invalid signature method")
	}

	keys, err := fetchGooglePublicKey()

	if err != nil {
		return nil, err
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid não encontrado no cabeçalho do token")
	}

	pubKey, ok := keys[kid]
	if !ok {
		return nil, errors.New("chave pública não encontrada para o kid fornecido")
	}

	return pubKey, nil
}

func fetchGooglePublicKey() (map[string]*rsa.PublicKey, error) {
	const googleCertsUrl = "https://www.googleapis.com/oauth2/v3/certs"

	resp, err := http.Get(googleCertsUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching public keys from Google - status %d", resp.StatusCode)
	}

	var certs GoogleCerts
	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		return nil, fmt.Errorf("error decoding public keys: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey)
	for _, key := range certs.Keys {
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)

		if err != nil {
			return nil, err
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)

		if err != nil {
			return nil, err
		}

		n := new(big.Int).SetBytes(nBytes)
		e := int(new(big.Int).SetBytes(eBytes).Int64())

		pubKey := &rsa.PublicKey{
			N: n,
			E: e,
		}

		keys[key.Kid] = pubKey
	}

	return keys, nil
}
