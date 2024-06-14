package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// TODO: Move to an environment variable
const SECRET = "very secretive secret to ever secret"

type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type UserPayload struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func (u *UserPayload) IsExpired() bool {
	return time.Now().Unix() > u.ExpiresAt
}

type Token struct {
	Header    *Header
	Payload   *UserPayload
	header    string
	payload   string
	signature string
	encoder   *base64.Encoding
}

func (t Token) String() string {
	return t.header + "." + t.payload + "." + t.signature
}

func (t *Token) IsExpired() bool {
	return t.Payload.IsExpired()
}

func getBase64Encoder() *base64.Encoding {
	return base64.URLEncoding.WithPadding(base64.NoPadding)
}

func EncryptJWT(subjectName string) (*Token, error) {
	header := &Header{"HS256", "JWT"}
	now := time.Now()
	iat := now.Add(time.Hour * 24 * 14)
	user := &UserPayload{subjectName, iat.Unix(), now.Unix()}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}
	payloadBytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	encoder := getBase64Encoder()
	headerBase64 := encoder.EncodeToString(headerBytes)
	payloadBase64 := encoder.EncodeToString(payloadBytes)
	signatureHash := hmac.New(sha256.New, []byte(SECRET))
	signatureHash.Write([]byte(headerBase64 + "." + payloadBase64))

	signature := encoder.EncodeToString(signatureHash.Sum(nil))

	jwt := headerBase64 + "." + payloadBase64 + "." + signature

	token := &Token{header, user, headerBase64, payloadBase64, jwt, encoder}

	return token, nil
}

func DecryptJWT(key string) (token *Token, err error) {
	parts := strings.Split(key, ".")

	var header Header
	var payload UserPayload

	token = &Token{&header, &payload, parts[0], parts[1], key, getBase64Encoder()}

	headerJson, err := token.encoder.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(headerJson, token.Header)
	if err != nil {
		return nil, err
	}

	payloadJson, err := token.encoder.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(payloadJson, token.Payload)
	if err != nil {
		return nil, err
	}

	return token, nil
}
