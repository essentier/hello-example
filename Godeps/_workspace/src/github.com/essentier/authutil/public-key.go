package authutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	publicKeyString = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzMyDkeRr37kCM8bD0cHu
SZcqczEA3X8yQvzk0Pxg1NO1qDwUXPt2aOQZx0x2jFJHHTSgJkg9bwVOFLmSqBTN
UH/Fz+lMQw+2UWNGmbe07n/+y58o5jC+jZNgwCQ9EdHasg6sZfEJJY35lY0eMhn2
YvVQwJZOq83gALBtclnKZ7wJ6Ue5BqjPP0ucY4WNk7cd04ySY8z3AV92R0K9TxwQ
3s7IKjMJrlcbNB7PTfj+yIzbrIT230omQLAh+oKlKf8KxAYVdyB9pI72DV+ljuEe
O+T2Qbp0uIH117A0tFe40hRI5rOEN0o18yHejTGOx2a8/+naP4nJ0FGI5hpafGSq
8wIDAQAB
-----END PUBLIC KEY-----
`
)

var publicKey *rsa.PublicKey = parsePublicKey()

func parsePublicKey() *rsa.PublicKey {
	publicKeyData, _ := pem.Decode([]byte(publicKeyString))
	publicKeyImported, err := x509.ParsePKIXPublicKey(publicKeyData.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)
	if !ok {
		panic(err)
	}
	return rsaPub
}

