package main

import (
	"bufio"
	"io"
	"os"
	"github.com/smtc/glog"
	"strings"
)

/**
打开文件并处理
创建人:邵炜
创建时间:2016年6月1日09:40:03
输入参数:filePath 文件地址
输出参数:文件对象 错误对象
 */
func openFile(filePath string) (*os.File,error) {
	var (
		fs *os.File
		err error
	)
	fs,err=os.Open(filePath)
	if err != nil {
		glog.Error("open file is error! filePath: %s err: %s \n",filePath,err.Error())
		return nil,err
	}
	glog.Info("file open success! filePath: %s \n",filePath)
	return fs,nil
}

/**
逐行读文件
创建人:邵炜
创建时间:2016年6月1日09:49:45
输入参数:文件地址 赛选条件方法 回调方法
 */
func readFile(filepath string,where func(string) bool,callBack func(string)) {
	var (
		readAll =false
		readByte []byte
		line []byte
		err error
		contentLine string
	)
	read,err:=openFile(filepath)
	if err != nil {
		return
	}
	defer read.Close()
	buf:=bufio.NewReader(read)
	for err!=io.EOF {
		if err!=nil {
			glog.Error("read error! err: %s \n",err.Error())
		}
		if readAll {
			readByte,readAll,err=buf.ReadLine()
			line=append(line,readByte...)
		}else{
			readByte,readAll,err=buf.ReadLine()
			line=append(line,readByte...)
			if len(strings.TrimSpace(string(line)))==0 {
				continue
			}
			contentLine=string(line)
			if where(contentLine) {
				callBack(contentLine)
			}
			line=line[:0]
		}
	}
}