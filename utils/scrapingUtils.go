package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func NewClient() *http.Client {
	return http.DefaultClient
}
func SaveToHTML(url, path string) error {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func HtmlToDoc(htmlFilePath string) (*goquery.Document, error) {
	reader, err := os.Open(htmlFilePath)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GetDocFromUrl(logger *zap.SugaredLogger, method, url string) (*goquery.Document, error) {
	cl := NewClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorw("could not create new request", "error", err, "requestUrl", url, "requestMethod", method)
		return nil, err
	}
	res, err := cl.Do(req)
	if err != nil {
		logger.Errorw("could not send the request", "error", err, "requestUrl", url, "requestMethod", method)
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Errorw("could not parse the response body", "error", err, "requestUrl", url, "requestMethod", method)
		return nil, err
	}
	return doc, nil
}
