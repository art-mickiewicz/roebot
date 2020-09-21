package cbr

import (
	srv "19u4n4/roebot/services"
)

func init() {
	srv.RegisterVariable("cbr_usdrub", "курс доллара США к рублю")
	srv.RegisterVariable("cbr_eurrub", "курс евро к рублю")
	srv.RegisterVariable("cbr_cnyrub", "курс китайского юаня к рублю")
	srv.RegisterVariable("cbr_gbprub", "курс фунта стерлингов к рублю")
	srv.RegisterService("cbr", SyncCBR)
}

func SyncCBR() {
	//url := "https://www.cbr-xml-daily.ru/daily_json.js"
	srv.SetValue("cbr_usdrub", "70")
	srv.SetValue("cbr_eurrub", "80")
	srv.SetValue("cbr_cnyrub", "11")
	srv.SetValue("cbr_gbprub", "100")
}
