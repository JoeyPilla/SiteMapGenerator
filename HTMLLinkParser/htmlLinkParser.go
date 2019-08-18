package HTMLLinkParser

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

var Links []Link

func ParseLink(htmlByte []byte) []string {
	doc, err := html.Parse(strings.NewReader(string(htmlByte)))
	if err != nil {
		fmt.Println(err)
	}
	return parse(doc)
}

func parse(n *html.Node) []string {
	links := getLinks(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, parse(c)...)
	}
	return links
}

func getLinks(n *html.Node) []string {
	var links []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "a" {
			for _, a := range c.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
	}
	return links
}
