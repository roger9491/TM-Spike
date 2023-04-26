package service

import (
	"TM-Spike/dao"
	"TM-Spike/init/configinit"
	"TM-Spike/model"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

const (
	successful = 1
	failed     = 0

	isNotDelete = 0
	isDelete    = 1
)

var (
	OrderRepo OrderRepoInterface = &orderSQL{}
)

type OrderRepoInterface interface {
	Initialize(*gorm.DB, *redis.Client)
	Order(model.Product, *gin.Context) (model.ProductInfo, error)
	Create(string, int64, *gin.Context) (int, error)
}

type orderSQL struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (od *orderSQL) Initialize(db *gorm.DB, rdb *redis.Client) {
	od.db = db
	od.rdb = rdb
}

func (od *orderSQL) Order(product model.Product, c *gin.Context) (produdctInfo model.ProductInfo, err error) {
	// redis
	// 判斷有沒有緩存
	if hasStock := loadDataFromDBIfCacheMiss(od.db, od.rdb, c, product.ProductName); !hasStock {
		produdctInfo.Status = failed
		return
	}
	// 緩存有庫存
	// lua腳本
	var luaScript = redis.NewScript(`
		local value = redis.call("Get", KEYS[1])
		if( value - KEYS[2] >= 0 ) then
			local leftStock = redis.call("DecrBy", KEYS[1], KEYS[2])
			return leftStock
		else
			return -1
		end
	`)

	// 執行腳本
	n, err := luaScript.Run(c, od.rdb, []string{product.ProductName, "1"}).Int()
	if err != nil {
		panic(err)
	}
	if n == -1 {
		produdctInfo.Status = failed
		return
	}

	// kafka 解偶mysql
	kafkaURL := configinit.KafkaIP + ":" + configinit.KafkaPort
	// create a topic if not exist
	_, err = kafka.DialLeader(context.Background(), "tcp", kafkaURL, configinit.KafkaPort, 0)
	if err != nil {
		panic(err.Error())
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    configinit.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	// 形成 message
	msg := kafka.Message{
		Value: []byte(product.ProductName),
	}
	
	// 送出 message
	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Println("送出失敗",err)
		return
	}

	// 修改數據庫
	// tx := od.db.Begin()
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		tx.Rollback()
	// 		log.Println("12", err)
	// 		return
	// 	}
	// }()

	// if err = dao.UpdateOrder(product.ProductName, tx, c); err != nil {
	// 	panic(err)
	// }
	// produdctInfo.Status = successful
	// if err = dao.CreateOrder(product.ProductName, 1, tx, c); err != nil {
	// 	panic(err)
	// }

	// tx.Commit()

	// 成功下訂單
	log.Println("下訂單成功")
	produdctInfo.Product = product.ProductName
	produdctInfo.Status = successful
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

// 防止緩存擊穿
var mu sync.Mutex

// 查詢緩存是否已載入，如果沒有則從資料庫獲取
func loadDataFromDBIfCacheMiss(db *gorm.DB, rdb *redis.Client, c *gin.Context, productNmae string) (hasStock bool) {
	fmt.Println("asdadsa")
	if _, err := rdb.Get(c, productNmae).Result(); err == nil {
		return true
	}
	fmt.Println("3333")
	mu.Lock()
	defer mu.Unlock()
	// 從資料庫獲取庫存
	count, err := dao.SelectCountFromProduct(productNmae, db, c)
	if err != nil {
		panic(err)
	}
	if count <= 0 {
		log.Println("--商品賣完--")
		return false
	}

	countStr := strconv.Itoa(count)
	/*
		Values == false 代表key不存在
			更新
	*/
	// LUA 更新緩存
	var luaScript = redis.NewScript(`
		local value = redis.call("Get", KEYS[1])
		if( value == false ) then
			redis.call("set" , KEYS[1],KEYS[2])
			return KEYS[2]
		end
		return -1
	`)
	n, err := luaScript.Run(c, rdb, []string{productNmae, countStr}).Result()
	if err != nil {
		log.Println(err)
	}
	if n == count {
		log.Println("key 載入緩存 ")
	}

	return true
}
