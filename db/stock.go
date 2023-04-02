package db

import "LittleBeeMark/CoolGoPkg/apply_kfk/model"

func InsertStockInfo(param *model.Stock) error {
	return DB.Create(param).Error
}

func GetStock(id int64) (model.Stock, error) {
	stock := model.Stock{}
	err := DB.Model(&model.Stock{}).Find(&stock).Error
	if err != nil {
		return stock, err
	}

	return stock, nil
}

func InsertFianceInfo(param *model.Fiance) error {
	return DB.Create(param).Error
}

func GetFiance(id int64) (model.Fiance, error) {
	f := model.Fiance{}
	err := DB.Model(&model.Fiance{}).Find(&f).Error
	if err != nil {
		return f, err
	}

	return f, nil
}
