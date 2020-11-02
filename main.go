package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	var buf []byte

	// get product availability
	isAvailable, err := getProductAvailability(os.Getenv("URL"), os.Getenv("STORE"), &buf)
	if err != nil {
		log.Fatal(err)
	}

	if isAvailable {
		// initialize Telegram bot
		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
		if err != nil {
			log.Fatal(err)
		}

		// parse chatID
		chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHATID"), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		// create file to be sent
		file := tgbotapi.FileBytes{
			Name:  "screenshot.png",
			Bytes: buf,
		}

		// send availability screenshot via Telegram message
		photo := tgbotapi.NewPhotoUpload(chatID, file)
		photo.Caption = "Product is available"
		if _, err := bot.Send(photo); err != nil {
			log.Fatal(err)
		}

	}

	// log the execution result
	log.Println(isAvailable)

}

func getProductAvailability(url, store string, buf *[]byte) (bool, error) {
	var isAvailable bool

	// create context
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// run
	var availability string
	if err := chromedp.Run(ctx,
		// open url
		chromedp.Navigate(url),
		// wait to page load
		chromedp.WaitVisible(`#onetrust-banner-sdk > div`),
		// accept cookies
		chromedp.Click(`#onetrust-accept-btn-handler`, chromedp.ByID),
		// check store stock
		chromedp.Click(`#content > div > div > div > div.range-revamp-product__subgrid.product-pip.js-product-pip > div.range-revamp-product__buy-module-container.range-revamp-product__grid-gap > div > div.range-revamp-product-availability > div.js-stockcheck-section > div > span > a`, chromedp.NodeVisible),
		// select store
		chromedp.SendKeys(`#change-store-input`, store),
		chromedp.Click(`#radio`, chromedp.NodeVisible),
		chromedp.Click(`#range-modal-mount-node > div > div:nth-child(3) > div > div.range-revamp-modal__buttons > button`, chromedp.NodeVisible),
		// check product availability
		chromedp.OuterHTML(`#range-modal-mount-node > div > div:nth-child(3) > div > div.range-revamp-modal__content > div:nth-child(3) > div > div.range-revamp-store-info__header > div > span`, &availability),
		// take screenshot
		chromedp.CaptureScreenshot(buf),
	); err != nil {
		return isAvailable, err
	}

	// parse availability
	isAvailable = strings.Contains(availability, "success")

	return isAvailable, nil
}
