package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os"
	"strconv"

	"github.com/itchyny/base58-go"
)

// CheckDbError handles catching errors
// while trying to connect to the database
func CheckDbErr(errMsg string, err error) {
	if err != nil {
		log.Println(errMsg, err)
	}
}

// IsValidURL checks if the URL is valid
func IsValidURL(urlString string) bool {
	u, err := url.Parse(urlString)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// sha256Of converts a url string into an encrypted byteslice
func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

// base58Encoded converts SHA256 encrypted byteslice to string
// using base58 encoding
func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

// GenerateShortLink generates a short url for the user
func GenerateShortLink(initialLink string, userId string) string {
	urlHashBytes := sha256Of(initialLink + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:7]
}

// UserIDToUint connverts userID gotten after request passes through middle ware
// from interface to uint
func UserIDToUint(any interface{}) uint {
	stringToUint := func(value interface{}) (uint, error) {
		str, ok := value.(string)
		if !ok {
			return 0, fmt.Errorf("value is not a string")
		}

		intValue, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return 0, err
		}

		uintValue := uint(intValue)
		return uintValue, nil
	}

	uintValue, err := stringToUint(any)
	if err != nil {
		log.Fatal("Conversion error:", err)
	}

	return uintValue
}
