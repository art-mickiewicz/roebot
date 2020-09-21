package binance

import (
	srv "19u4n4/roebot/services"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://api.binance.com/api/v3/ticker/price"

func init() {
	srv.RegisterVariable("binance_btcusdt", "BTC / USDT")
	srv.RegisterVariable("binance_ethusdt", "ETH / USDT")
	srv.RegisterVariable("binance_bchusdt", "ВСН / USDT")
	srv.RegisterVariable("binance_usdtrub", "USDT / RUB")
	srv.RegisterService("binance", "@hourly", SyncBinance)
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
	// srv.SetValue("binance_btcusdt", fmt.Sprintf("%f", symbolsMap["BTCUSDT"]))
	// srv.SetValue("binance_ethusdt", fmt.Sprintf("%f", symbolsMap["ETHUSDT"]))
	// srv.SetValue("binance_bchusdt", fmt.Sprintf("%f", symbolsMap["BCHUSDT"]))
	// srv.SetValue("binance_usdtrub", fmt.Sprintf("%f", symbolsMap["USDTRUB"]))
	srv.SetValue("binance_btcusdt", symbolsMap["BTCUSDT"])
	srv.SetValue("binance_ethusdt", symbolsMap["ETHUSDT"])
	srv.SetValue("binance_bchusdt", symbolsMap["BCHUSDT"])
	srv.SetValue("binance_usdtrub", symbolsMap["USDTRUB"])
}
