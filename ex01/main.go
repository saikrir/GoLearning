package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

var rawHtml = `
	<!DOCTYPE html>
	<html>
	<body>
		<h1>My First Heading</h1>
		<p>My first paragraph.</p>
		<p>HTML 
			<a href="https://www.w3schools.com/html/html_images.asp">images</a> 
		are defined with the img tag:</p>
		<img src="xxx.jpg" width="104" height="142">
	</body>
	</html>
`

func visit(node *html.Node, wordCnt *int, imgCnt *int) {

	for el := node; el != nil; el = el.FirstChild {

		if el.Type == html.TextNode {
			*wordCnt += len(strings.Fields(el.Data))
		} else if el.Type == html.ElementNode && el.Data == "img" {
			*imgCnt += 1
		}

		nextSibling := el.NextSibling
		if nextSibling != nil {
			visit(nextSibling, wordCnt, imgCnt)
		}
	}

}

func countWordsAndImages(node *html.Node) (int, int) {
	var wordCnt, imgCnt int
	visit(node, &wordCnt, &imgCnt)
	return wordCnt, imgCnt
}

func main() {
	htmlReader := strings.NewReader(rawHtml)
	nodePtr, err := html.Parse(htmlReader)

	if err != nil {
		log.Fatalln("Failed to parse html")
	}

	fmt.Println(countWordsAndImages((nodePtr)))
}
