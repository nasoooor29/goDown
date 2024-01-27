package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func MatchUrlHosts(url1, url2 string) (bool, error) {
	url1Obj, err := url.Parse(url1)
	if err != nil {
		return false, err
	}
	url2Obj, err := url.Parse(url2)
	if err != nil {
		return false, err
	}
	return url1Obj.Host == url2Obj.Host, nil
}

func EnsureDirExists(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("ensure directory exists (%s): %w", path, err)
	}
	return nil
}

func DecodeAtob(str string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	decodedString := string(decodedBytes)
	return decodedString, nil
}

func ExtractBetweenSingleQuotes(input string) (string, bool) {
	re := regexp.MustCompile(`'([^']*)'`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", false
	}
	return matches[1], true
}

func GenHashBasedOnTime() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	uniqueInfo := rand.Int63()
	id := fmt.Sprintf("%v_%v", timestamp, uniqueInfo)
	hash := sha256.New()
	hash.Write([]byte(id))
	hashed := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashed)
	return hashedString
}

func RemoveNonEnglishLetters(input string) string {
	reg := regexp.MustCompile("[^a-zA-Z]+")
	result := reg.ReplaceAllString(input, "")
	return result
}

func ReturnAsJson(anything any) string {
	jsonString, err := json.Marshal(anything)
	if err != nil {
		return err.Error()
	}
	str := string(jsonString)
	return strings.ReplaceAll(str, "\\u0026", "&")
}


func PadNumber(width, num int) string {
	return fmt.Sprintf("%0*d", width, num)
}