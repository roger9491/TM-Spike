package dao

import (
	"log"

	"gorm.io/gorm"
)

func SelectOrder(product string, tx *gorm.DB) (count int64, err error) {
	sql := "select count from product where productname = ? for update;"

	if err = tx.Raw(sql, product).Scan(&count).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func UpdateOrder(product string, tx *gorm.DB) (err error) {
	sql := "update product set count = count - 1 where productname = ?"
	if err = tx.Exec(sql, product).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func CreateProduct(productName string, count int64, tx *gorm.DB) (err error) {
	sql := "INSERT INTO product (productname, count) VALUES ( ?, ? );"

	if err = tx.Exec(sql, productName, count).Error; err != nil {
		log.Println("err ", err)
	}
	return
}
