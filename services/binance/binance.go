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
	srv.RegisterVariable("binance_btc", "BTC / RUB")
	srv.RegisterVariable("binance_eth", "ETH / RUB")
	srv.RegisterVariable("binance_bchusdt", "ВСН / USDT")
	srv.RegisterVariable("binance_usdt", "USDT / RUB")
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
	srv.SetValue("binance_btc", symbolsMap["BTCRUB"])
	srv.SetValue("binance_eth", symbolsMap["ETHRUB"])
	srv.SetValue("binance_bchusdt", symbolsMap["BCHUSDT"])
	srv.SetValue("binance_usdt", symbolsMap["USDTRUB"])
	srv.Commit()
}
