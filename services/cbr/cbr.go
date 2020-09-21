package cbr

import (
	srv "19u4n4/roebot/services"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	n "golang.org/x/text/number"
)

const url = "https://www.cbr-xml-daily.ru/daily_json.js"

var printer *message.Printer

func init() {
	printer = message.NewPrinter(language.Russian)
	srv.RegisterVariable("cbr_usd", "курс доллара США к рублю")
	srv.RegisterVariable("cbr_eur", "курс евро к рублю")
	srv.RegisterVariable("cbr_cny", "курс китайского юаня к рублю")
	srv.RegisterVariable("cbr_gbp", "курс фунта стерлингов к рублю")
	srv.RegisterService("cbr", "@hourly", SyncCBR)
}

func formatValue(fv float64) string {
	return printer.Sprintf("%v", n.Decimal(
		fv,
		n.Scale(2),
		n.Pad(' '),
		n.FormatWidth(5),
	))
}

func SyncCBR() {
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

	var cbrResp CBRResponse
	if err := json.Unmarshal(body, &cbrResp); err != nil {
		log.Println(err)
		return
	}

	srv.SetValue("cbr_usd", formatValue(cbrResp.Valute["USD"].Value))
	srv.SetValue("cbr_eur", formatValue(cbrResp.Valute["EUR"].Value))
	srv.SetValue("cbr_cny", formatValue(cbrResp.Valute["CNY"].Value))
	srv.SetValue("cbr_gbp", formatValue(cbrResp.Valute["GBP"].Value))
	srv.Commit()
}
