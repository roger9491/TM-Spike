package dao

import (
	"TM-Spike/model"
	"log"

	"gorm.io/gorm"
)






func SelectOrder(product string, tx *gorm.DB) (count int64, err error) {
	sql := "select count from product where productname = ? AND isdelete = 0 for update;"

	if err = tx.Raw(sql, product).Scan(&count).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func SelectProduct(productName string, tx *gorm.DB) (products []model.Product, err error){
	sql := "select * from product where productname = ? AND isdelete = 0"
	
	if err = tx.Raw(sql, productName).Scan(&products).Error; err != nil {
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

func UpdateOrderIsDelete(product string, tx *gorm.DB) (err error) {
	sql := "update product set isdelete = 1 where productname = ?"
	if err = tx.Exec(sql, product).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func CreateProduct(productName string, count int64, isDelete int, tx *gorm.DB) (err error) {
	sql := "INSERT INTO product (productname, count, isdelete) VALUES ( ?, ?, ? );"

	if err = tx.Exec(sql, productName, count, isDelete).Error; err != nil {
		log.Println("err ", err)
	}
	return
}
