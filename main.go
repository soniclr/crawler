// Package crawler
/**
 * @Author: Soren.Liang
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2023/9/19 11:12
 * @Software : GoLand
 */
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func main() {
	var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)
	if err != nil {
		fmt.Printf("read content failed:%v\n", err)
		return
	}

	matches := headerRe.FindAllSubmatch(body, -1)
	for _, m := range matches {
		fmt.Println("fetch card news:", string(m[1]))
	}

	numLinks := strings.Count(string(body), "<a")
	fmt.Println("numLinks: ", numLinks)

	fmt.Println("==== use htmlquery ==== ")

	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("htmlquery.Parse failed:%v\n", err)
		return
	}

	nodes := htmlquery.Find(doc, `//div[@class="small_cardcontent__BTALp"]//a[@target="_blank"]`)
	for _, node := range nodes {
		fmt.Println("fetch card ", node.FirstChild.Data)
	}

	fmt.Println("==== use goquery ==== ")
	docs, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("goquery.NewDocumentFromReader failed:%v\n", err)
		return
	}

	docs.Find("div.small_cardcontent__BTALp a[target=_blank]").Each(func(i int, s *goquery.Selection) {
		// 获取匹配标签中的文本
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
	return
}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get url error : %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get url error : status code %v", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	byteSrc, err := r.Peek(1024)
	if err != nil {
		log.Printf("fetcher error : %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(byteSrc, "")
	return e
}
