package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func initChrome(ctx context.Context) error {
	var err error

	// create chrome instance
	c, err := chromedp.New(ctx, chromedp.WithLog(log.Printf))
	if err != nil {
		return err
	}

	go func() {
		var lastURL string
		for {
			u := currentURL()
			if lastURL != u {
				lastURL = u
				err = c.Run(ctx, chromedp.Tasks{chromedp.Navigate(u)})
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}
