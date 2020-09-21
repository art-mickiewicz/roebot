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
	srv.RegisterVariable("cbr_usdrub", "курс доллара США к рублю")
	srv.RegisterVariable("cbr_eurrub", "курс евро к рублю")
	srv.RegisterVariable("cbr_cnyrub", "курс китайского юаня к рублю")
	srv.RegisterVariable("cbr_gbprub", "курс фунта стерлингов к рублю")
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

	srv.SetValue("cbr_usdrub", fmt.Sprintf("%f", cbrResp.Valute["USD"].Value))
	srv.SetValue("cbr_eurrub", fmt.Sprintf("%f", cbrResp.Valute["EUR"].Value))
	srv.SetValue("cbr_cnyrub", fmt.Sprintf("%f", cbrResp.Valute["CNY"].Value))
	srv.SetValue("cbr_gbprub", fmt.Sprintf("%f", cbrResp.Valute["GBP"].Value))
	srv.Commit()
}
