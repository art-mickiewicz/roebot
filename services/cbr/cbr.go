package cbr

import (
	srv "19u4n4/roebot/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://www.cbr-xml-daily.ru/daily_json.js"

func init() {
	srv.RegisterVariable("cbr_usd", "курс доллара США к рублю")
	srv.RegisterVariable("cbr_eur", "курс евро к рублю")
	srv.RegisterVariable("cbr_cny", "курс китайского юаня к рублю")
	srv.RegisterVariable("cbr_gbp", "курс фунта стерлингов к рублю")
	srv.RegisterService("cbr", "@hourly", SyncCBR)
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

	srv.SetValue("cbr_usd", fmt.Sprintf("%f", cbrResp.Valute["USD"].Value))
	srv.SetValue("cbr_eur", fmt.Sprintf("%f", cbrResp.Valute["EUR"].Value))
	srv.SetValue("cbr_cny", fmt.Sprintf("%f", cbrResp.Valute["CNY"].Value))
	srv.SetValue("cbr_gbp", fmt.Sprintf("%f", cbrResp.Valute["GBP"].Value))
	srv.Commit()
}
