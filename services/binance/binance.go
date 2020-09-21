package binance

import (
	srv "19u4n4/roebot/services"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	n "golang.org/x/text/number"
)

const url = "https://api.binance.com/api/v3/ticker/price"

var printer *message.Printer

func init() {
	printer = message.NewPrinter(language.Russian)
	srv.RegisterVariable("binance_btc", "BTC / RUB")
	srv.RegisterVariable("binance_eth", "ETH / RUB")
	srv.RegisterVariable("binance_bchusdt", "ВСН / USDT")
	srv.RegisterVariable("binance_usdt", "USDT / RUB")
	srv.RegisterService("binance", "@hourly", SyncBinance)
}

func formatValue(v string) string {
	fv, _ := strconv.ParseFloat(v, 64)
	return printer.Sprintf("%v", n.Decimal(
		fv,
		n.Scale(2),
		n.Pad(' '),
		n.FormatWidth(10),
	))
}

func SyncBinance() {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var symbols []Symbol
	if err := json.Unmarshal(body, &symbols); err != nil {
		log.Println(err)
		return
	}

	symbolsMap := make(map[string]string)
	for _, sym := range symbols {
		symbolsMap[sym.Name] = sym.Price
	}
	srv.SetValue("binance_btc", formatValue(symbolsMap["BTCRUB"]))
	srv.SetValue("binance_eth", formatValue(symbolsMap["ETHRUB"]))
	srv.SetValue("binance_bchusdt", formatValue(symbolsMap["BCHUSDT"]))
	srv.SetValue("binance_usdt", formatValue(symbolsMap["USDTRUB"]))
	srv.Commit()
}
