package config

import (
	"encoding/base64"
)

var tokens = map[string]string{
	"dG1kYg==": "ZTQ0Y2U2YmU1NzVlN2E2OWQ3NDU4YmQ4ZjFmZDllMmY=",
	"dHZkYg==": "OUExRkQ2MTdGMkMyNDgxOQ==",
}

// GetToken returns the token for the given resource.
func GetToken(resource string) string {
	key := base64.StdEncoding.EncodeToString([]byte(resource))
	val, _ := base64.StdEncoding.DecodeString(tokens[key])
	return string(val)
}
