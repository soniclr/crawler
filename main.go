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
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/chromedp/chromedp"
	"github.com/soniclr/crawler/collect"
)

func main() {
	if err := ChromeDB(); err != nil {
		fmt.Println(err)
	}
}

func ChromeDB() error {
	// 1、创建谷歌浏览器实例
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	// 2、设置context超时时间
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 3、爬取页面，等待某一个元素出现,接着模拟鼠标点击，最后获取数据
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),
		chromedp.WaitVisible(`body > footer`),
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		chromedp.Value(`#example-After textarea`, &example),
	)

	if err != nil {
		fmt.Printf("chromedp.Run error:%v", err)
		return err
	}

	log.Printf("Go's time.After example:\\n%s", example)

	return nil
}

func FetcherCount() error {
	var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

	// url := "https://book.douban.com/subject/1007305/"
	url := "https://www.thepaper.cn/"

	var f collect.Fetcher = &collect.BrowserFetcher{
		Timeout: 300 * time.Millisecond,
	}
	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("read content failed:%v\n", err)
		return err
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
		return err
	}

	nodes := htmlquery.Find(doc, `//div[@class="small_cardcontent__BTALp"]//a[@target="_blank"]`)
	for _, node := range nodes {
		fmt.Println("fetch card ", node.FirstChild.Data)
	}

	fmt.Println("==== use goquery ==== ")
	docs, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("goquery.NewDocumentFromReader failed:%v\n", err)
		return err
	}

	docs.Find("div.small_cardcontent__BTALp a[target=_blank]").Each(func(i int, s *goquery.Selection) {
		// 获取匹配标签中的文本
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})

	return nil
}
