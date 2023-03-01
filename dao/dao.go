package dao

import (
	"TM-Spike/model"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)






func SelectOrder(product string, tx *gorm.DB, c *gin.Context) (count int64, err error) {
	sql := "select count from product where productname = ? AND isdelete = 0 for update;"

	if err = tx.WithContext(c.Request.Context()).Raw(sql, product).Scan(&count).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func SelectProduct(productName string, tx *gorm.DB, c *gin.Context) (products []model.Product, err error){
	sql := "select * from product where productname = ? AND isdelete = 0"
	
	if err = tx.WithContext(c.Request.Context()).Raw(sql, productName).Scan(&products).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func UpdateOrder(product string, tx *gorm.DB, c *gin.Context) (err error) {
	sql := "update product set count = count - 1 where productname = ?"
	if err = tx.WithContext(c.Request.Context()).Exec(sql, product).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func UpdateOrderIsDelete(product string, tx *gorm.DB, c *gin.Context) (err error) {
	sql := "update product set isdelete = 1 where productname = ?"
	if err = tx.WithContext(c.Request.Context()).Exec(sql, product).Error; err != nil {
		log.Println("err ", err)
	}

	return
}

func CreateProduct(productName string, count int64, isDelete int, tx *gorm.DB, c *gin.Context) (err error) {
	sql := "INSERT INTO `product` (productname, count, isdelete) VALUES ( ?, ?, ? );"

	if err = tx.WithContext(c.Request.Context()).Exec(sql, productName, count, isDelete).Error; err != nil {
		log.Println("err ", err)
	}
	return
}
// 建立訂單
func CreateOrder(productName string, count int64, tx *gorm.DB, c *gin.Context) (err error) {
	sql := "INSERT INTO `order` (productname, count) VALUES ( ?, ?);"

	if err = tx.WithContext(c.Request.Context()).Exec(sql, productName, count).Error; err != nil {
		log.Println("err ", err)
	}
	return
}
