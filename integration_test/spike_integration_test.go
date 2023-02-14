package integration__test

import (
	"TM-Spike/model"
	"TM-Spike/router"
	"TM-Spike/service"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestPutUser(t *testing.T) {

	service.OrderRepo.Initialize(database(t))

	gin.SetMode(gin.TestMode)

	samples := []struct {
		inputJSON  string
		statusCode int
		errMessage string
	}{
		{
			inputJSON:  `{"productname":"apple", "count": 200}`,
			statusCode: 201,
			errMessage: "",
		},
	}
	for _, v := range samples {
		r := gin.Default()
		r.PUT("/product", router.PutProduct)
		req, err := http.NewRequest(http.MethodPut, "/product", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		var product model.Product
		err = json.Unmarshal(rr.Body.Bytes(), &product)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		fmt.Println("this is the response data: ", product)
		fmt.Println("CODE ", rr.Code, v.statusCode)
		assert.Equal(t, v.statusCode, rr.Code)
		// if v.statusCode == 201 {

		// }
	}
}
