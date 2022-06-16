package golibs

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

type WebScraping struct {
	driver *agouti.WebDriver
	page   *agouti.Page
}

type optionsWebScraping struct {
	IsHeadless   bool
	IsSecretMode bool
	EnablePopup  bool
}

type optionWebScraping func(optionsWebScraping *optionsWebScraping)

func IsHeadless(isHeadless bool) optionWebScraping {
	return func(opts *optionsWebScraping) {
		opts.IsHeadless = isHeadless
	}
}

func IsSecretMode(isSecretMode bool) optionWebScraping {
	return func(opts *optionsWebScraping) {
		opts.IsSecretMode = isSecretMode
	}
}

func EnablePopup(enablePopup bool) optionWebScraping {
	return func(opts *optionsWebScraping) {
		opts.EnablePopup = enablePopup
	}
}

func NewWebScraping(options ...optionWebScraping) (*WebScraping, error) {

	opts := &optionsWebScraping{
		IsHeadless:   true,
		IsSecretMode: true,
		EnablePopup:  false,
	}
	for _, option := range options {
		option(opts)
	}

	chromeOpts := []string{
		"--disable-gpu", // ref: https://developers.google.com/web/updates/2017/04/headless-chrome#cli
		"no-sandbox",    // ref: https://github.com/theintern/intern/issues/878
		"disable-dev-shm-usage",
	}
	// ヘッドレスブラウザ
	if opts.IsHeadless {
		chromeOpts = append(chromeOpts, "--headless")
	}
	// シークレットモード
	if opts.IsSecretMode {
		chromeOpts = append(chromeOpts, "–incognito")
	}
	chromeOptions := agouti.ChromeOptions(
		"args",
		chromeOpts,
	)

	ws := new(WebScraping)
	client := new(http.Client)
	ws.driver = agouti.ChromeDriver(
		// httpクライアントを設定する(設定していない場合は、内部的にhttp.DefaultClientが使われる)
		agouti.HTTPClient(client),
		// // なんのタイムアウトなのか不明
		// agouti.Timeout(190),
		// SSL証明書が無効な場合は拒否する(デフォルトではすべて受け入れる)
		agouti.RejectInvalidSSL,
		chromeOptions,
	)

	err := ws.driver.Start()
	if err != nil {
		log.Fatal("NewWebScraping error:", err)
		return nil, err
	}

	ws.page, err = ws.driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatal("NewWebScraping error:", err)
		return nil, err
	}

	// ポップアップの無効化
	if opts.EnablePopup {
		err = ws.page.CancelPopup()
		if err != nil {
			log.Fatal("NewWebScraping error:", err)
			return nil, err
		}
	}

	return ws, nil
}

func (ws *WebScraping) Close() error {
	return ws.driver.Stop()
}

// ページ操作

func (ws *WebScraping) GetPage() *agouti.Page {
	return ws.page
}

func (ws *WebScraping) NavigatePage(requestUrl string) error {
	return ws.page.Navigate(requestUrl)
}

// 画面操作

// By ID

func (ws *WebScraping) GetStringByID(id string) string {
	return ws.page.FindByID(id).String()
}

func (ws *WebScraping) SetStringByID(id string, str string) error {
	element := ws.page.FindByID(id)
	err := element.Clear()
	if err != nil {
		return err
	}
	return element.Fill(str)
}

func (ws *WebScraping) GetValueByID(id string) *agouti.Selection {
	return ws.page.FindByID(id)
}

func (ws *WebScraping) ClickButtonByID(id string) error {
	return ws.page.FindByID(id).Click()
}

func (ws *WebScraping) SubmitButtonByID(id string) error {
	return ws.page.FindByID(id).Submit()
}

// By Class

func (ws *WebScraping) GetStringByClass(class string) string {
	return ws.page.FindByClass(class).String()
}

func (ws *WebScraping) SetStringByClass(class string, str string) error {
	element := ws.page.FindByClass(class)
	err := element.Clear()
	if err != nil {
		return err
	}
	return element.Fill(str)
}

func (ws *WebScraping) GetValueByClass(class string) *agouti.Selection {
	return ws.page.FindByClass(class)
}

func (ws *WebScraping) ClickButtonByClass(class string) error {
	return ws.page.FindByClass(class).Click()
}

func (ws *WebScraping) SubmitButtonByClass(class string) error {
	return ws.page.FindByClass(class).Submit()
}

// By Name

func (ws *WebScraping) GetStringByName(name string) string {
	return ws.page.FindByName(name).String()
}

func (ws *WebScraping) SetStringByName(name string, str string) error {
	element := ws.page.FindByName(name)
	err := element.Clear()
	if err != nil {
		return err
	}
	return element.Fill(str)
}

func (ws *WebScraping) GetValueByName(name string) *agouti.Selection {
	return ws.page.FindByName(name)
}

func (ws *WebScraping) ClickButtonByName(name string) error {
	return ws.page.FindByName(name).Click()
}

func (ws *WebScraping) SubmitButtonByName(name string) error {
	return ws.page.FindByName(name).Submit()
}

// By Text

func (ws *WebScraping) ClickButtonByText(text string) error {
	return ws.page.FindByButton(text).Click()
}

func (ws *WebScraping) SubmitButtonByText(text string) error {
	return ws.page.FindByButton(text).Submit()
}

// By Selector

func (ws *WebScraping) GetAllElementsBySelector(selector string) *agouti.MultiSelection {
	return ws.page.All(selector)
}

func (ws *WebScraping) ClickButtonBySelector(selector string) error {
	return ws.page.Find(selector).Click()
}

func (ws *WebScraping) SubmitButtonBySelector(selector string) error {
	return ws.page.Find(selector).Submit()
}

// Execute JavaScript

func (ws *WebScraping) ExecJavaScript(script string, variable map[string]interface{}) error {
	return ws.page.RunScript(script, variable, nil)
}

// Others

func (ws *WebScraping) GetFileByWeb(fileName string, url string) error {

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("GetFileByWeb error:", err)
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal("GetFileByWeb error:", err)
		return err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("GetFileByWeb error:", err)
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		log.Fatal("GetFileByWeb error:", err)
		return err
	}

	return nil
}

func (ws *WebScraping) ReadFile(filePath string) (*goquery.Document, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("GetFileByWeb error:", err)
		return nil, err
	}

	stringReader := strings.NewReader(string(data))

	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		log.Fatal("GetFileByWeb error:", err)
		return nil, err
	}

	return doc, nil
}
