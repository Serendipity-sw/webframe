package main

import (
	"fmt"
	"github.com/smtc/glog"
	"os"
)

/**
判断文件或文件夹是否存在
创建人:邵炜
创建时间:2016年12月21日17:07:42
输入参数:需要查询的文件或文件夹路径
输出参数:返回值true存在 否则不存在  错误对象
*/
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return false, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
根据文件夹路径创建文件,如文件存在则不做任何操作
创建人:邵炜
创建时间:2016年12月21日17:23:54
输入参数:文件夹路径
输出参数:错误对象
*/
func createFileProcess(path string) error {
	fileExists, err := pathExists(path)
	if err != nil {
		glog.Error("serverRun upLoadFile exists! path: %s err: %s \n", path, err.Error())
		return err
	}
	if !fileExists {
		err = os.MkdirAll(path, 0666)
		if err != nil {
			glog.Error("serverRun mkdir uploadFile err! path: %s err: %s \n", path, err.Error())
			return err
		}
	}
	return nil
}
