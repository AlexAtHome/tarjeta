package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
	"strings"
)

// TODO: Move to an environment variable
const SECRET = "very secretive secret to ever secret"

type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type User struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func (u *User) IsExpired() bool {
	return time.Now().Unix() > u.ExpiresAt
}

type Token struct {
	Header    *Header
	Payload   *User
	header    string
	payload   string
	signature string
	encoder *base64.Encoder
}

func (t Token) String() string {
	return t.header + "." + t.payload + "." + t.signature
}

func (t *Token) IsValid() bool {
	sigHashSum, err := t.encoder.DecodeString(t.signature)
	if err != nil {
		return false
	}
	sig := string(sigHashSum)
	// TODO: use hmac sha256 to validate the unpacked signsture
	
}

func EncryptJWT(subjectName string) (*Token, error) {
	header := &Header{"HS256", "JWT"}
	now := time.Now()
	iat := now.Add(time.Hour * 24 * 14)
	user := &User{subjectName, iat.Unix(), now.Unix()}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}
	payloadBytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	encoder := base64.URLEncoding.WithPadding(base64.NoPadding)
	headerBase64 := encoder.EncodeToString(headerBytes)
	payloadBase64 := encoder.EncodeToString(payloadBytes)
	signatureHash := hmac.New(sha256.New, []byte(SECRET))
	signatureHash.Write([]byte(headerBase64 + "." + payloadBase64))

	signature := encoder.EncodeToString(signatureHash.Sum(nil))

	token := &Token{header, user, headerBase64, payloadBase64, signature, &encoder}

	return token, nil
}

func DecryptJWT(token string) (*Token, error) {
	parts := strings.split(token, ".")
	token := &Token{nil, nil, parts[0], parts[1], parts[2], encoder}
	encoder := base64.URLEncoding.WithPadding(base64.NoPadding)
	headerJson := token.encoder.DecodeString(token.header)
	payloadJson := token.encoder.DecodeString(token.payload)
	err := json.Unmarshal(token.header, token.Header)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(token.payload, token.Payload)
	if err != nil {
		return nil, err
	}
	return token, nil
}