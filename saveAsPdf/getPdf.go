package saveAsPdf

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func getPageAsPdf(URL string, docsPath string) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(URL, &buf)); err != nil {
		log.Fatal(err)
	}

	fileName := "../assets/" + strings.ReplaceAll(strings.Split(URL, docsPath)[1], "/", "_") + ".pdf"

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
