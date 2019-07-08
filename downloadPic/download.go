package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

var url = `https://cn.bing.com/images/search?q=%E6%BC%AB%E5%A8%81+%E7%94%B5%E5%BD%B1%E6%B5%B7%E6%8A%A5`

func forAllNode(node *html.Node, f func(*html.Node)) {
	if node == nil {
		return
	}

	f(node)
	forAllNode(node.FirstChild, f)
	forAllNode(node.NextSibling, f)
}

func main() {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	n, _ := html.Parse(resp.Body)
	forAllNode(n, func(node *html.Node) {
		if node.Data == "img" {
			fmt.Println(node.Attr)
			ok := false
			src := ""
			for _, v := range node.Attr {
				if v.Key == "class" && v.Val == "mimg" {
					ok = true
				} else if v.Key == "src" {
					src = v.Val
				}
			}
			if ok {
				fmt.Println(src)
			}
		}
	})
	return
}
