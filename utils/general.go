package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"regexp"
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
