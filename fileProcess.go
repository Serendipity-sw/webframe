package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/smtc/glog"
	"io"
	"os"
	"path"
	"strings"
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

/**
获取文件名称及后缀名 未防止文件无后缀名,固这里返回值为数组对象
创建人:邵炜
创建时间:2016年9月7日16:40:14
输入参数: 文件路劲或文件名称
输出参数: 文件名 文件后缀名 数组 第一项为文件名称 第二项为文件后缀名
*/
func getMyFileName(filePaths string) *[]string {
	fileName := path.Base(filePaths)
	suffixName := path.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, suffixName)
	files := []string{fileName, suffixName}
	return &files
}

/**
写文件
创建人:邵炜
创建时间:2016年9月7日16:31:39
输入参数:文件内容 写入文件的路劲(包含文件名)
输出参数:错误对象
*/
func fileCreateAndWrite(content *[]byte, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		glog.Error("fileCreateAndWrite os openFile error! fileName: %s err: %s \n", fileName, err.Error())
		return err
	}
	defer f.Close()
	_, err = f.Write(*content)
	if err != nil {
		go glog.Error("fileCreateAndWrite write error! content: %v fileName: %s err: %s \n", *content, fileName, err.Error())
		return err
	}
	glog.Info("fileCreateAndWrite run success! fileName: %s content: %s  \n", fileName, string(*content))
	return nil
}

/**
根据文件路径将文件写入到别的目录下去
创建人:邵炜
创建时间:2016年12月22日11:37:11
输入参数:文件路径 写入文件的路径(包含文件名)
输出参数:错误对象
*/
func fileCreateAndWriteByPath(pathStr string, afterFileName string) error {
	contentLine, err := readFileByLine(pathStr)
	if err != nil {
		glog.Error("fileCreateAndWriteByPath readFileByLine run err! pathStr: %s err: %s \n", pathStr, err.Error())
		return err
	}
	if len(*contentLine) == 0 {
		return errors.New("fileCreateAndWriteByPath  read file mepty! ")
	}
	var (
		contentByte []byte
	)
	for _, value := range *contentLine {
		contentByte = append(contentByte, []byte(fmt.Sprintf("%s\n", value))...)
	}
	err = fileCreateAndWrite(&contentByte, afterFileName)
	if err != nil {
		glog.Error("fileCreateAndWriteByPath fileCreateAndWrite run err! pathStr: %s afterFileName: %s err: %s \n", pathStr, afterFileName, err.Error())
		return err
	}
	return nil
}

/**
文件读取逐行进行读取
创建人:邵炜
创建时间:2016年9月20日10:23:41
输入参数: 文件路劲
输出参数: 字符串数组(数组每一项对应文件的每一行) 错误对象
*/
func readFileByLine(filePath string) (*[]string, error) {
	var (
		readAll     = false
		readByte    []byte
		line        []byte
		err         error
		contentLine []string
	)
	fs, err := os.Open(filePath)
	if err != nil {
		glog.Error("readFileByLine open error! filePath: %s err: %s \n", filePath, err.Error())
		return nil, err
	}
	defer fs.Close()
	buf := bufio.NewReader(fs)
	for err != io.EOF {
		if err != nil {
			glog.Error("readFileByLine read error! err: %s \n", err.Error())
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
			contentLine = append(contentLine, string(line))
			line = line[:0]
		}
	}
	glog.Info("readFileByLine run success! filePath: %s fileContent: %v \n", filePath, contentLine)
	return &contentLine, nil
}
