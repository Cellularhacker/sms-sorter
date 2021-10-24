package renewSpamDB

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sms-sorter/model/finefssCategory"
	"sms-sorter/model/thecall"
	"sms-sorter/util"
	"strings"
	"time"
)

const finefssURL = "http://fine.fss.or.kr/main/fin_comp/fincomp_inqui/comsearch01list.jsp"

const (
	businessOptionSelector = "div#container > section.content > div.fixed_width > div.srchCont_box.innerTab > form > fieldset > dl > dd > span.barDesign > select.select_default > option"
)

func FineFss() error {
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

func getBusinessList() ([]finefssCategory.FineFssCategory, error) {
	req, err := http.NewRequest("GET", finefssURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(util.MinifyBody(resp.Body))
	if err != nil {
		return nil, err
	}

	list := make([]finefssCategory.FineFssCategory, 0)

	doc.Find(businessOptionSelector).Each(func(i int, sel *goquery.Selection) {
		fc := finefssCategory.New()
		fc.Value, _ = sel.Attr("value")
		fc.Text = sel.Text()

		list = append(list, *fc)
	})

	return list, nil
}

func repeatFineFss(boTable string, page int) (bool, error) {
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
