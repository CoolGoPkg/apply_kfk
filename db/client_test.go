package db

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/conf"
	"LittleBeeMark/CoolGoPkg/apply_kfk/model"
	"testing"
)

func TestInitDB(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	InitDB()
}

func TestInsertStockInfo(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	InitDB()
	err := InsertStockInfo(&model.Stock{
		ProdCode:    "000001.SS",
		ProdName:    "上证指数",
		TradeStatus: "Trade",
		Last:        3000.00,
	})
	if err != nil {
		t.Log("err : ", err)
		return
	}
}

func TestInsertFianceInfo(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	InitDB()
	err := InsertFianceInfo(&model.Fiance{
		Eps:         1.02,
		Bps:         2.01,
		TotalShares: 1000.00,
	})
	if err != nil {
		t.Log("err : ", err)
		return
	}
}
