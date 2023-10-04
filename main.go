package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	p "github.com/AbdelilahOu/DocsPdf-go/saveAsPdf"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

var (
	visitedPathset = make(map[string]bool)
	docsBase       string
	baseUrl        string
	stop           int
)

func main() {
	stop = 3
	// get command line args
	firstPageInDocs := os.Args[1]
	docsBase = os.Args[2]
	// get base url
	baseUrl = strings.Split(firstPageInDocs, docsBase)[0]
	// get all links
	getLinksRecursively(firstPageInDocs)
	// save link as pdf
	for k, _ := range visitedPathset {
		p.GetPageAsPdf(baseUrl+k, baseUrl)
	}

}

func getLinksRecursively(url string) {
	if stop == 0 {
		return
	}

	// check if link is alreadu visisted
	if visitedPathset[url] {
		return
	}
	stop = stop - 1
	// add tto visisted
	visitedPathset[url] = true
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
					getLinksRecursively(strings.Split(href, "#")[0])
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
