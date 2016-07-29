package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
接口反馈对象
创建人:邵炜
创建时间:2016年7月29日19:56:32
*/
type requestData struct {
	ResultCode string //返回码
	Message    string //返回消息信息
}

/**
JSON请求数据返回
创建人:邵炜
创建时间:2016年1月4日20:23:36
输入参数: gin指针 bo判断是否使用错误返回对象 param泛型参数
输出参数: 无
数据反馈由gin进行
*/
func jsonPRequest(c *gin.Context, bo bool, param interface{}) {

	var cb string

	if c.Request.Method == "GET" {
		cb = c.Query("callback")
	} else {
		cb = c.PostForm("callback")
	}

	jsonResP := &requestData{
		ResultCode: "00000",
		Message:    "",
	}

	switch paramType := param.(type) {
	case string:
		if bo {
			jsonResP.ResultCode = "00001"
		}
		jsonResP.Message = param.(string)
		if cb != "" {
			b, _ := json.Marshal(jsonResP)
			c.Data(http.StatusOK, "application/javascript", []byte(fmt.Sprintf("%s(%s)", cb, b)))
		} else {
			c.JSON(http.StatusOK, jsonResP)
		}
	case int32:
		jsonResP.Message = strconv.Itoa(int(paramType))
		if cb != "" {
			b, _ := json.Marshal(jsonResP)
			c.Data(http.StatusOK, "application/javascript", []byte(fmt.Sprintf("%s(%s)", cb, b)))
		} else {
			c.JSON(http.StatusOK, jsonResP)
		}
	case int64:
		jsonResP.Message = strconv.Itoa(int(paramType))
		if cb != "" {
			b, _ := json.Marshal(jsonResP)
			c.Data(http.StatusOK, "application/javascript", []byte(fmt.Sprintf("%s(%s)", cb, b)))
		} else {
			c.JSON(http.StatusOK, jsonResP)
		}
	default:
		if cb != "" {
			b, _ := json.Marshal(paramType)
			c.Data(http.StatusOK, "application/javascript", []byte(fmt.Sprintf("%s(%s)", cb, b)))
		} else {
			c.JSON(http.StatusOK, param)
		}
	}
}
