package golibs

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func RequestPost(requestUrl string, params map[string]string, headers map[string]string) (string, error) {

	// パラメータを作る
	values := url.Values{}
	for key, val := range params {
		values.Add(key, val)
	}

	// リクエストを作る
	// 例えばPOSTでパラメータつけてUser-Agentを指定する
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal("Error:", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return "", err
	}

	return requestCommon(req, headers)
}

func RequestGet(requestUrl string, params map[string]string, headers map[string]string) (string, error) {

	// パラメータを作る
	paramArray := make([]string, len(params))
	for key, val := range params {
		paramArray = append(paramArray, fmt.Sprintf("%s=%s", key, url.QueryEscape(val)))
	}
	paramString := strings.Join(paramArray, "&")

	// リクエストを作る
	req, err := http.NewRequest("GET", requestUrl+"?"+paramString, nil)
	if err != nil {
		log.Fatal("Error:", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return "", err
	}

	return requestCommon(req, headers)
}

func requestCommon(req *http.Request, headers map[string]string) (string, error) {

	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	// リスエスト投げてレスポンスを得る
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Request error:", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return "", err
	}

	// レスポンスをNewDocumentFromResponseに渡してドキュメントを得る
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal("Request error:", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return "", err
	}

	ret, err := doc.Html()
	if err != nil {
		log.Fatal("Request error:", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return "", err
	}

	return ret, nil
}
