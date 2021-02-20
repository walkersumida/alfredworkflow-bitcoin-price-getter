package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dustin/go-humanize"
)

// Item is Alfred's item struct.
type Item struct {
	Type     string `json:"type"`
	Icon     string `json:"icon"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
}

// Menu is Alfred's menu struct.
type Menu struct {
	Items []Item `json:"items"`
}

// Ticker is bitFlyer's API response.
// https://lightning.bitflyer.com/docs
type Ticker struct { // Generated by JSON-to-Go(https://mholt.github.io/json-to-go/)
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func outputPrice(ticker Ticker, currencyCode string) {
	var item Item
	item.Icon = "./icon.png"
	item.Title = fmt.Sprintf("Bitcoin Ask on bitFlyer: %s %s", humanize.Commaf(ticker.BestAsk), currencyCode)
	item.Subtitle = fmt.Sprintf("Bitcoin Bid on bitFlyer: %s %s", humanize.Commaf(ticker.BestBid), currencyCode)
	item.Arg = fmt.Sprintf("Bitcoin Ask on bitFlyer: %s %s", humanize.Commaf(ticker.BestAsk), currencyCode)

	outputFormat(item)
}

func outputError(msg string) {
	var item Item
	item.Icon = "./icon.png"
	item.Title = "Error: " + msg

	outputFormat(item)
}

func outputInfo(msg string) {
	var item Item
	item.Icon = "./icon.png"
	item.Title = msg

	outputFormat(item)
}

func outputFormat(item Item) {
	var menu Menu
	menu.Items = append(menu.Items, item)

	menuJSON, _ := json.Marshal(menu)
	fmt.Println(string(menuJSON))
}

func main() {
	url := "https://api.bitflyer.com"
	path := "/v1/ticker?product_code=BTC_"
	var ticker Ticker

	flag.Parse()
	currencyCode := flag.Arg(0)
	currencyCode = strings.ToUpper(currencyCode)

	if currencyCode != "JPY" && currencyCode != "USD" {
		outputInfo("Please enter `USD` or `JPY`")
		return
	}

	path = path + currencyCode

	response, err := http.Get(url + path)

	if err != nil {
		outputError(err.Error())
		return
	}

	if response.StatusCode != 200 {
		outputError(
			fmt.Sprintf("Response status error(%d): %s", response.StatusCode, response.Status))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)

	if err := json.Unmarshal(body, &ticker); err != nil {
		outputError(err.Error())
		return
	}

	outputPrice(ticker, currencyCode)
}
