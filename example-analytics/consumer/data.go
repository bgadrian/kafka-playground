package consumer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bgadrian/kafka-playground/example-analytics/sortedAggr"
)

type messageProcessator func([]byte)

var aggrPageViews = sortedAggr.NewAggArray()
var aggrClicks = sortedAggr.NewAggArray()
var aggrAddToCartName = sortedAggr.NewAggArray()
var aggrAddToCartColor = sortedAggr.NewAggArray()
var aggrBuyName = sortedAggr.NewAggArray()
var aggrBuyColor = sortedAggr.NewAggArray()
var aggrErrorModules = sortedAggr.NewAggArray()
var latestErrors = sortedAggr.NewLatest(10)

var aggrBuyAmountName = sortedAggr.NewAggArray()
var aggrBuyAmountColor = sortedAggr.NewAggArray()

func ProcessMessage(msg []byte) {
	var event = make(map[string]string)
	err := json.Unmarshal(msg, &event)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(event)

	//all events are in the aggregator/main.go/buildDispatcher
	//we presume that are not malformed, don't do this in production
	eventType := event["event"]

	switch eventType {
	case "pageView":
		aggrPageViews.Increment(event["page"], 1)
		//logger("%s", aggrPageViews.Top(2))
	case "click":
		aggrClicks.Increment(event["element"], 1)
		//logger("%s", aggrClicks.Top(2))
	case "addToCart":
		aggrAddToCartName.Increment(event["name"], 1)
		aggrAddToCartColor.Increment(event["color"], 1)

		//logger("%s", aggrAddToCartName.Top(2))
		//logger("%s", aggrAddToCartColor.Top(2))
	case "buy":
		aggrBuyName.Increment(event["name"], 1)
		aggrBuyColor.Increment(event["color"], 1)

		//logger("%s", aggrBuyName.Top(2))
		//logger("%s", aggrBuyColor.Top(2))

		price, _ := strconv.ParseFloat(event["price"], 64)
		qnt, _ := strconv.ParseInt(event["amount"], 10, 64)
		total := float32(float32(price) * float32(qnt))

		aggrBuyAmountName.Increment(event["name"], total)
		aggrBuyAmountColor.Increment(event["color"], total)
	case "error":
		//if _, ok := event["location"]; ok {
		latestErrors.Add(event["location"])

		parts := strings.Split(event["location"], ":")
		if len(parts) == 3 {
			module := parts[0]
			aggrErrorModules.Increment(module, 1)
		}

		//logger("%s", latestErrors.Get())
		//logger("%s", aggrErrorModules.Top(2))
		//}

	default:
		//logger("unknown " + eventType)
	}
}
