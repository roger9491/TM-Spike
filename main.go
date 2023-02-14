package main

import (
	"TM-Spike/init/configinit"
	"TM-Spike/init/routerinit"
	"TM-Spike/init/sqlinit"
	"TM-Spike/router"
	"TM-Spike/service"
	"log"
	"os"
)

func main() {

	// 設置日誌輸出配置

	f, err := os.OpenFile("mrmaster_log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		// test
		log.Fatal("OpenFile is failed")
	}

	defer f.Close()

	log.SetOutput(f)

	log.Println("test test ")

	configinit.LoadEnv()

	// 初始化資料庫
	service.OrderRepo.Initialize(sqlinit.InitMySQL(configinit.DBUsername, configinit.DBPassword, configinit.DBHost, configinit.DBPort, configinit.DBName))

	// 加載路由
	routerinit.Include(router.OrderApi)
	r := routerinit.InitRouters()

	err = r.Run(configinit.HostIp + ":" + configinit.HostPort) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Println("err ", err.Error())
	}
}
