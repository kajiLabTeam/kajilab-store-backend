package service

import (
	"fmt"
	"os"
	"time"

	"kajilab-store-backend/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AssetService struct{}

// お金を増減させる
func (AssetService) IncreaseMoney(money int64) error {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 現在のお金を取得
	asset := model.Asset{}
	result := db.Last(&asset)
	if result.Error != nil {
		fmt.Printf("財産取得失敗 %v", result.Error)
		return result.Error
	}

	// DBへ更新
	// result = db.Model(&model.Asset{}).Where("id = 1").Update("money", asset.Money + money)
	result = db.Create(&model.Asset{Money: asset.Money + money, Debt: asset.Debt})
	if result.Error != nil {
		fmt.Printf("財産更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 現在の予算を取得
func (AssetService) GetAsset() (model.Asset, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return model.Asset{}, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 現在の予算をDBから取得
	asset := model.Asset{}
	result := db.Last(&asset)
	if result.Error != nil {
		fmt.Printf("財産取得失敗 %v", result.Error)
		return model.Asset{}, result.Error
	}

	return asset, nil
}

// 予算の履歴を取得
func (AssetService) GetAssetHistory(day int64) ([]model.Asset, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	start := time.Now().AddDate(0, 0, -int(day))

	var records []model.Asset
	err = db.Where("created_at >= ?", start).Order("created_at asc").Find(&records).Error
	if err != nil {
		return nil, err
	}

	// 日付ごとの最新データをmapに保存
	assetMap := map[string]model.Asset{}
	for _, r := range records {
		key := r.CreatedAt.Format("2006-01-02")
		assetMap[key] = r
	}

	assets := make([]model.Asset, 0)

	for i := 0; i <= int(day); i++ {
		targetDay := time.Now().AddDate(0, 0, -(int(day) - i))
		key := targetDay.Format("2006-01-02")

		if asset, ok := assetMap[key]; ok {
			assets = append(assets, model.Asset{
				Money: asset.Money,
				Debt:  asset.Debt,
			})
		} else {
			assets = append(assets, model.Asset{
				Money: -1,
				Debt:  -1,
			})
		}
	}

	return assets, nil

	// db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return []model.Asset{}, err
	// }
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return []model.Asset{}, err
	// }
	// defer sqlDB.Close()

	// assets := make([]model.Asset, 0)
	// // 現在の予算をDBから取得
	// for i := 0; i <= int(day); i++ {
	// 	asset := model.Asset{}
	// 	dayAgo := time.Now().AddDate(0, 0, 0-(int(day)-i))
	// 	startOfDay := time.Date(dayAgo.Year(), dayAgo.Month(), dayAgo.Day(), 0, 0, 0, 0, dayAgo.Location())
	// 	endOfDay := time.Date(dayAgo.Year(), dayAgo.Month(), dayAgo.Day(), 23, 59, 59, 999999, dayAgo.Location())
	// 	result := db.Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Last(&asset)
	// 	if result.Error != nil {
	// 		assets = append(assets, model.Asset{
	// 			Money: -1,
	// 			Debt:  -1,
	// 		})
	// 	} else {
	// 		assets = append(assets, model.Asset{
	// 			Money: asset.Money,
	// 			Debt:  asset.Debt,
	// 		})
	// 	}
	// }

	// return assets, nil
}

// 財産を更新する
func (AssetService) UpdateAsset(asset *model.Asset) error {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 財産情報をDBへ登録
	result := db.Create(asset)
	if result.Error != nil {
		fmt.Printf("商品情報更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}
