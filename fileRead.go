package main

import (
	"bufio"
	"github.com/smtc/glog"
	"io"
	"os"
	"strings"
)

/**
打开文件并处理
创建人:邵炜
创建时间:2016年6月1日09:40:03
输入参数:filePath 文件地址
输出参数:文件对象 错误对象
*/
func openFile(filePath string) (*os.File, error) {
	var (
		fs  *os.File
		err error
	)
	fs, err = os.Open(filePath)
	if err != nil {
		glog.Error("open file is error! filePath: %s err: %s \n", filePath, err.Error())
		return nil, err
	}
	glog.Info("file open success! filePath: %s \n", filePath)
	return fs, nil
}

/**
逐行读文件
创建人:邵炜
创建时间:2016年6月1日09:49:45
输入参数:文件地址 赛选条件方法(传入匹配字符串) 回调方法(传入匹配字符串)
*/
func readFile(filepath string, where []func(string) bool, callBack []func(string)) (int, int) {
	var (
		readAll     = false
		readByte    []byte
		line        []byte
		err         error
		contentLine string
		matchNumber = 0
		count       = 0
	)
	read, err := openFile(filepath)
	if err != nil {
		return 0, 0
	}
	defer read.Close()
	buf := bufio.NewReader(read)
	for err != io.EOF {
		if err != nil {
			glog.Error("read error! err: %s \n", err.Error())
		}
		if readAll {
			readByte, readAll, err = buf.ReadLine()
			line = append(line, readByte...)
		} else {
			readByte, readAll, err = buf.ReadLine()
			line = append(line, readByte...)
			if len(strings.TrimSpace(string(line))) == 0 {
				continue
			}
			contentLine = string(line)
			for index, value := range where {
				if value(contentLine) {
					if len(callBack) == (index+1) && callBack[index] != nil {
						callBack[index](contentLine)
					}
					matchNumber++
				}
			}
			line = line[:0]
			count++
		}
	}
	return matchNumber, count
}
