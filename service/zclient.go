package service

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var debugging = false

func _log(msg string) {
	if debugging {
		fmt.Println(msg)
	}
}

func _createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	Header := map[string][]string{
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Accept-Language": {"en-US,en;q=0.5"},
	}

	for h := range Header {
		req.Header.Set(h, Header[h][0])
	}

	return req, nil
}

func _getResponse(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	if resp.StatusCode == 403 {
		for i := 0; i < 3; i++ {
			backoff := time.Duration(math.Pow(2, float64(i))) * (time.Second / 10)
			time.Sleep(backoff)

			resp, err = client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("error making GET request: %w f", err)
			}

			if resp.StatusCode != 403 {
				break
			}

			_log(fmt.Sprintf("Retry number %d\n", i+1))
		}
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error loading HTML: %s for %s", resp.Status, req.URL.String())
	}

	return resp, nil
}

func _parseBody(resp *http.Response) (*goquery.Document, error) {
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %w", err)
	}

	return doc, nil
}

func debugError(err error) {
	if debugging {
		fmt.Println(err)
	}
}

func get(url string) (*goquery.Document, error) {
	req, err := _createRequest(url)
	if err != nil {
		debugError(err)
		return nil, err
	}

	resp, err := _getResponse(req)
	if err != nil {
		debugError(err)
		return nil, err
	}

	doc, err := _parseBody(resp)
	if err != nil {
		debugError(err)
		return nil, err
	}

	return doc, nil
}
