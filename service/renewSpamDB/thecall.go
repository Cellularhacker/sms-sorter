package renewSpamDB

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sms-sorter/service/proxy2"
	"sms-sorter/util"
	"sms-sorter/util/logger"
	"strings"
	"time"

	"sms-sorter/model/thecall"
)

const url = "http://www.thecall.co.kr/bbs/board.php"
const testURL = "http://ch002.cafe24.com/thecall_whitelist.html"

const (
	articleSelector     = "div.phone-list > form > article"
	phoneNumberSelector = "h2 > a"
	subjectSelector     = "p"
	canonicalSelector   = "head > link[rel=\"canonical\"]"

	paramBoTable          = "bo_table"
	paramBoTablePhone     = "phone"
	paramBoTableWhitelist = "whitelist"
	paramPage             = "page"
)

func TheCall() error {
	logger.L.Infof("[TheCall] Parsing %s", paramBoTablePhone)
	i := 1
	for {
		fmt.Printf("%d ", i)
		more, err := repeat(paramBoTablePhone, i)
		if err != nil {
			return err
		}
		if !more {
			break
		}
		i++
		<-time.NewTicker(3 * time.Second).C
	}
	logger.L.Infof("[TheCall] Parsing %s", paramBoTableWhitelist)
	i = 1
	for {
		fmt.Printf("%d ", i)
		more, err := repeat(paramBoTableWhitelist, i)
		if err != nil {
			return err
		}
		if !more {
			break
		}
		i++
		<-time.NewTicker(3 * time.Second).C
	}

	return nil
}

func repeat(boTable string, page int) (bool, error) {
	uri := fmt.Sprintf("%s?%s=%s&%s=%d", url, paramBoTable, boTable, paramPage, page)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("User-Agent", userAgent)

	client := proxy2.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(util.MinifyBody(resp.Body))
	if err != nil {
		return false, err
	}

	path, _ := doc.Find(canonicalSelector).Attr("href")

	amount := doc.Find(articleSelector).Size()
	if amount <= 0 {
		return false, nil
	}

	doc.Find(articleSelector).Each(func(i int, selection *goquery.Selection) {
		c := thecall.New()
		c.Subject = selection.Find(subjectSelector).Text()
		c.PhoneNumber = selection.Find(phoneNumberSelector).Text()
		c.IsWhiteList = strings.Contains(path, "whitelist")
		go c.Upsert()
	})

	return true, nil
}
