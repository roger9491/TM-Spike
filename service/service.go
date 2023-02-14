package service

import (
	"TM-Spike/dao"
	"TM-Spike/model"
	"log"

	"gorm.io/gorm"
)

const (
	successful = 1
	failed     = 0
)

var (
	OrderRepo OrderRepoInterface = &orderSQL{}
)

type OrderRepoInterface interface {
	Initialize(*gorm.DB)
	Order(model.Product) (model.ProductInfo, error)
	Create(string, int64) error
}

type orderSQL struct {
	db *gorm.DB
}

func (od *orderSQL) Initialize(db *gorm.DB) {
	od.db = db
}

func (od *orderSQL) Order(product model.Product) (produdctInfo model.ProductInfo, err error) {
	tx := od.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	count, err := dao.SelectOrder(product.ProductName, tx)
	if count > 0 {
		if err = dao.UpdateOrder(product.ProductName, tx); err != nil {
			panic(err)
		}
		produdctInfo.Status = successful
	} else {
		produdctInfo.Status = failed
	}
	tx.Commit()

	produdctInfo.Product = product.ProductName
	return
}

func (od *orderSQL) Create(productName string, count int64) (err error) {
	tx := od.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	if err = dao.CreateProduct(productName, count, tx); err != nil {
		panic(err)
	}

	tx.Commit()

	return
}
