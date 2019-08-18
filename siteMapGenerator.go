package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"./HTMLLinkParser"
)

type Url struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

type Urls struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	Temp    []Url
}

var links = make(map[string]linkStruct)
var baseUrl = ""

type linkStruct struct {
	visited bool
}

func main() {
	filePtr := flag.String("url", "http://calhoun.io", "url to parse")
	flag.Parse()
	baseUrl = *filePtr
	temp(getHTML(baseUrl))
	filename := "newstaffs.xml"
	file, _ := os.Create(filename)
	file.Write([]byte(xml.Header))
	xmlWriter := io.Writer(file)
	for !checkVisited() {

	}
	var keys []string
	for key, _ := range links {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	xmlURL := Urls{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Temp:  []Url{},
	}

	for _, key := range keys {
		xmlURL.Temp = append(xmlURL.Temp, Url{
			Loc: key,
		})
	}
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(xmlURL); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func checkVisited() bool {
	for key, link := range links {
		if !link.visited {
			temp(getHTML(key))
			return false
		}
	}
	return true
}

func getHTML(s string) []byte {
	links[s] = linkStruct{
		visited: true,
	}
	resp, _ := http.Get(s)
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return html
}

func temp(html []byte) {
	for _, link := range HTMLLinkParser.ParseLink(html) {
		if link[0] == string("/")[0] {
			link = baseUrl + link
		}
		if strings.Contains(link, "/calhoun.io") {
			if _, ok := links[link]; !ok {
				links[link] = linkStruct{
					visited: false,
				}
			}
		}
	}
}
