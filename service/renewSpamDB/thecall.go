package renewSpamDB

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sms-sorter/util"
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
	log.Printf("[TheCall] Parsing %s\n", paramBoTablePhone)
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
	fmt.Println()
	log.Printf("[TheCall] Parsing %s\n", paramBoTableWhitelist)
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
	url := fmt.Sprintf("%s?%s=%s&%s=%d", url, paramBoTable, boTable, paramPage, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("User-Agent", userAgent)

	client := &http.Client{}
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
