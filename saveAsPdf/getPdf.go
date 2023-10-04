package saveAsPdf

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GetPageAsPdf(URL string, baseUrl string) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	// get data
	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(URL, &buf)); err != nil {
		log.Fatal(err)
	}
	// create directory
	url, err := url.Parse(baseUrl)
	if err != nil {
		fmt.Println("coudnt parse url")
	}
	hostname := strings.Split(strings.TrimPrefix(url.Hostname(), "www."), ".")[0]
	// file path
	fileName := "../assets/" + hostname
	// splited url
	splitedUrl := strings.Split(func() string {
		if strings.Contains(URL, "docs/") {
			return strings.Split(URL, "docs/")[1]

		}
		return URL + "/docs/"
	}(), "/")
	// last url
	for i, k := range splitedUrl {
		if i == len(splitedUrl)-1 {
			fileName = fileName + k + ".pdf"
		} else {
			_ = os.Mkdir(fileName+k, 0755)
			fileName = fileName + k + "/"
		}
	}

	if err := os.WriteFile(fileName, buf, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote sample.pdf")
}

// print a specific pdf page.
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
