package service

import (
	"TM-Spike/dao"
	"TM-Spike/model"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	successful = 1
	failed     = 0

	isNotDelete = 0
	isDelete 	= 1
)

var (
	OrderRepo OrderRepoInterface = &orderSQL{}
)

type OrderRepoInterface interface {
	Initialize(*gorm.DB)
	Order(model.Product, *gin.Context) (model.ProductInfo, error)
	Create(string, int64, *gin.Context) (int, error)
}

type orderSQL struct {
	db *gorm.DB
}

func (od *orderSQL) Initialize(db *gorm.DB) {
	od.db = db
}

func (od *orderSQL) Order(product model.Product, c *gin.Context) (produdctInfo model.ProductInfo, err error) {
	tx := od.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println(err)
			return
		}
	}()

	count, err := dao.SelectOrder(product.ProductName, tx, c)
	if count > 0 {
		if err = dao.UpdateOrder(product.ProductName, tx, c); err != nil {
			panic(err)
		}
		produdctInfo.Status = successful
		
		if err = dao.CreateOrder(product.ProductName, 1, tx, c); err != nil{
			panic(err)
		}

		if count == 1 {
			if err = dao.UpdateOrderIsDelete(product.ProductName, tx, c); err != nil{
				panic(err)
			}
		}

	} else {
		produdctInfo.Status = failed
	}
	tx.Commit()

	produdctInfo.Product = product.ProductName
	return
}

func (od *orderSQL) Create(productName string, count int64, c *gin.Context) (status int, err error) {
	tx := od.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	productList, err := dao.SelectProduct(productName, tx, c)
	if err != nil {
		panic(err)
	}

	if len(productList) > 0 {
		status = failed
		return
	}


	if err = dao.CreateProduct(productName, count, isNotDelete, tx, c); err != nil {
		panic(err)
	}

	tx.Commit()

	status = successful

	return
}
