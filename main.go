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
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn/"
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("http get error", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("http get error :%v, status code : %v\n", err, resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("read body error", err)
		return
	}

	fmt.Println("body: ", string(body))
}
