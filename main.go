package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

var (
	visitedPathset []string
	docsBase       string
	baseUrl        string
)

func main() {
	// get command line args
	firstPageInDocs := os.Args[1]
	docsBase = os.Args[2]
	// get base url
	baseUrl = strings.Split(firstPageInDocs, docsBase)[0]
	// get all links
	getLinksRecursively(firstPageInDocs)
}

func getLinksRecursively(url string) {
	// check if link is alreadu visisted
	if isAvailable(visitedPathset, url) {
		return
	}
	// add tto visisted
	visitedPathset = append(visitedPathset, url)
	fmt.Println("visit url :", url, baseUrl)
	// get full url
	updatedUrl := func() string {
		if strings.Split(url, docsBase)[0] == "/" {
			return baseUrl + url
		} else {
			return url
		}
	}()
	// parse full url
	doc, err := ParseWebApp(updatedUrl)
	if err != nil {
		fmt.Println("error parsing the webapp")
	}
	// get all the links in page
	allLinks := doc.Find("a")
	// loop through
	allLinks.Each(func(i int, s *goquery.Selection) {
		// get href
		href, exist := s.Attr("href")
		if exist {
			// is under docs
			if isAvailable(strings.Split(href, "/"), docsBase) {
				// this href is for current page
				if strings.Split(href, docsBase)[0] == "/" {
					getLinksRecursively(href)
				}
			} else {
				return
			}
		}
	})
}

func isAvailable(alpha []string, str string) bool {
	// iterate using the for loop
	for i := 0; i < len(alpha); i++ {
		// check
		if alpha[i] == str {
			// return true
			return true
		}
	}
	return false
}

// parse dynamic webapp
func ParseWebApp(url string) (*goquery.Document, error) {
	// where to store generated html
	var outterHTML string
	// create ctx
	ctx, cancel := chromedp.NewContext(context.Background())
	// cancel whene we done
	defer cancel()
	//
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		// js rendering happens asynchronously and this call seems to be enough to account for that
		chromedp.WaitReady(":root"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			// get html
			outterHTML, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	}); err != nil {
		return nil, fmt.Errorf("ParseWebApp(): ActionFunc(): %w", err)
	}
	// parse html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(outterHTML))
	if err != nil {
		return nil, fmt.Errorf("ParseWebApp(): goquery.NewDocumentFromReader(): %w", err)
	}

	return doc, nil
}
