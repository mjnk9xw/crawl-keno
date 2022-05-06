package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, _ := newChromedp()
	crawlPage(ctx)
}
func newChromedp() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("start-fullscreen", false),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	return ctx, cancel
}

func crawlPage(ctx context.Context) {

	url := "https://www.minhchinh.com/xo-so-dien-toan-keno.html"
	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(5 * time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			i := 2
			for true {
				node, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return err
				}
				res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					return err
				}
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
				if err != nil {
					return err
				}

				doc.Find(".wrapperKQKeno").Each(func(index int, row *goquery.Selection) {
					kyKQKeno := row.Find(".kyKQKeno").Text()
					timeKQ := row.Find(".timeKQ").Text()

					boxKQKeno := make([]string, 20)
					row.Find(".boxKQKeno div").Each(func(i int, row *goquery.Selection) {
						boxKQKeno[i] = strings.TrimSpace(row.Text())
					})
					fmt.Printf("%s - %s - %s \n", strings.TrimSpace(kyKQKeno), strings.TrimSpace(timeKQ), boxKQKeno)
				})

				clickPage := fmt.Sprintf(`*//a[@href="javascript:chosePage(%d)"]`, i)
				chromedp.Click(clickPage, chromedp.BySearch)
				time.Sleep(5 * time.Second)
				i++
			}

			return nil
		}),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}

}
