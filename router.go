package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smtc/glog"
	"io/ioutil"
	"net/http"
	"os"
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

/**
获取网页页面
创建人:邵炜
创建时间:2017年2月2日21:12:58
输入参数:gin对象
输出参数:无
数据反馈由gin进行
*/
func unitGetHtml(c *gin.Context) {
	htmlName := c.Param("name")
	c.HTML(http.StatusOK, htmlName, gin.H{
		"webFrameRoot": fmt.Sprintf("http://%s%s", c.Request.Host, rootPrefix),
	})
}

/**
公共文件上传主键
创建人:邵炜
创建时间:2016年12月21日17:27:33
输入参数:gin对象
数据反馈由gin进行
*/
func unitUploadFile(c *gin.Context) {
	fileName := c.Query("fname")
	c.Request.ParseMultipartForm(32 << 20)
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		glog.Error("unitUpLoadFile formFile err! err: %s \n", err.Error())
		jsonPRequest(c, true, "您提交的表单文件有误,或参数属性不为file!")
		return
	}
	defer file.Close()
	fmt.Fprintf(c.Writer, "%v", handler.Header)
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", upLoadFileDir, fileName), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		glog.Error("unitUpLoadFile openFile err! file: %s err: %s \n", fileName, err.Error())
		return
	}
	defer f.Close()
	n, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		glog.Error("unitUpLoadFile file seek err! fileName: %s err: %s \n", fileName, err.Error())
		return
	}
	fileContentByte, err := ioutil.ReadAll(file)
	_, err = f.WriteAt(fileContentByte, n)
	if err != nil {
		glog.Error("unitUpLoadFile file write err! fileName: %s err: %s \n", fileName, err.Error())
	}
}
