package router

import (
	"TM-Spike/model"
	"TM-Spike/service"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OrderApi(e *gin.Engine) {

	e.DELETE("/order", Order)
	e.PUT("/product", PutProduct)
}

/* 

{
	"productname": "apple"
}


*/
func Order(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	var product model.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	productInfo, err := service.OrderRepo.Order(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, productInfo)
}

/* 

{
	"productname": "apple"
	"count": 100	
}
*/
// 新增商品
func PutProduct(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}

	var product model.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var productInfo model.ProductInfo
	productInfo.Status, err = service.OrderRepo.Create(product.ProductName, product.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, productInfo)
}
