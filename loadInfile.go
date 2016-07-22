package main

import (
	"fmt"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxLoadBytes = 1000000 // 1048575
)

type loadInfile struct {
	sync.Mutex
	f        *os.File
	dir      string
	fn       string
	bytes    int
	seq      int
	interval int // how many second to load infile
}

var (
	loadFile     *loadInfile
	loadInterval = 15 // 1s
)

func init() {
	deferinit.AddInit(dataFileInit, nil, 989)

	deferinit.AddRoutine(dataProcessFunc)
}

func dataProcessFunc(ch chan struct{}, wg *sync.WaitGroup) {
	f := loadFile
	tmr := time.NewTicker(time.Second * time.Duration(f.interval))
	// 1天
	oneDay := time.Duration(time.Hour * time.Duration(24))

	go func() {
		<-ch
		tmr.Stop()
		loadFile.exit()
		wg.Done()
	}()

	for {
		<-tmr.C

		loadFile.rotate()
		filepath.Walk(loadFileDir, func(path string, info os.FileInfo, err error) error {
			now := time.Now()
			// 跳过目录
			if info.IsDir() {
				return nil
			}
			if err != nil {
				glog.Error("Walk file %s occurs error: %v\n", path, err)
				return nil
			}

			// 正在写入的文件, 不处理
			if ext := filepath.Ext(path); ext == ".tmp" {
				// 超过1天，删除
				if now.Sub(info.ModTime()) > oneDay {
					glog.Error("Found tmp file %s, modTime: %v now: %v, size: %v, delete it\n",
						path, info.ModTime(), now, info.Size())
					os.Remove(path)
				}
				return nil
			}
			// 空文件, 不处理
			if info.Size() == 0 {
				//glog.Info("File %s size is zero, skip it.\n", path)
				os.Remove(path)
				return nil
			}
			//数据入库操作
			//if ierr := fileLoadDataBase(path); ierr == nil {
			//	os.Remove(path)
			//}
			return nil
		})
	}
}

/**
数据写入文件构造函数
创建人:邵炜
创建时间:2016年7月22日10:59:10
*/
func dataFileInit() {
	createDir()
	loadFile = openLoadFile(loadInterval, loadFileDir)
	glog.Info("init load data successfully.\n")
}

/**
打开数据写入
创建人:邵炜
创建时间:2016年7月22日11:00:19
*/
func openLoadFile(second int, dir string) *loadInfile {
	lf := &loadInfile{
		dir:      dir,
		interval: second,
	}
	lf.rotate()
	return lf
}

func (f *loadInfile) exit() {
	f.Lock()
	f.close()
	f.f = nil
	f.Unlock()
}

/**
文件切换
创建人:郭铁
*/
func (f *loadInfile) rotate() {
	f.Lock()
	f.seq = 0
	f.close()
	f.createNewFile()
	f.Unlock()
}

/**
文件关闭 并切换写入地址
创建人:郭铁
*/
func (f *loadInfile) close() {
	if f.f != nil {
		f.f.Close()
		os.Rename(f.fn, f.fn[0:len(f.fn)-4]) // 去掉末尾的.tmp
	}
}

/**
新建文件
创建人:郭铁
*/
func (f *loadInfile) createNewFile() (err error) {
	tm := time.Now()

	f.bytes = 0 // 重置

	name := fmt.Sprintf("link-%04d%02d%02d-%02d%02d%02d-%d.tmp",
		tm.Year(), int(tm.Month()), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second(),
		f.seq)
	f.seq++
	f.fn = filepath.Join(loadFileDir, name)

	// linux下必须加O_WRONLY 选项，否则报错：bad file descriptor
	f.f, err = os.OpenFile(f.fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		glog.Error("create file %s failed: %s\n", f.fn, err.Error())
		return
	}
	//glog.Info("create new load file %s\n", f.fn)
	return
}

/**
判断文件夹是否存在 如文件夹不存在则创建该文件夹
创建人:邵炜
创建时间:2016年1月28日16:09:33
输入参数:无
输出参数:无
*/
func createDir() {
	loadFileDir = strings.TrimSpace(loadFileDir)

	if path.IsAbs(loadFileDir) == false {
		// 强制使用绝对路径
		root, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}
		//println(root)
		loadFileDir = path.Join(root, loadFileDir)
	}

	s, err := os.Stat(loadFileDir)
	if os.IsExist(err) {
		return
	}
	if err == nil {
		if s.IsDir() == false {
			panic(loadFileDir + " exist but not dir!")
		}
		return
	}
	if os.IsNotExist(err) == false {
		panic("createDir: " + err.Error())
	}

	err = os.MkdirAll(loadFileDir, 0666)
	if err != nil {
		panic(err)
	}

	glog.Info("loadfile Dir %s create success\n", loadFileDir)
}

/**
数据写入文件
创建人:邵炜
创建时间:2016年1月28日11:44:44
输入参数:separator 分隔符   content 数据参数
输出参数:无
*/
func logFile(separator string, content ...string) {
	loadFile.Lock()
	defer loadFile.Unlock()

	if loadFile.f == nil {
		glog.Warn("load File is nil\n")
		return
	}
	s := strings.Join(content, separator)
	// 控制load文件的大小
	if loadFile.bytes+len(s) > maxLoadBytes {
		loadFile.close()
		if err := loadFile.createNewFile(); err != nil {
			return
		}
	}

	loadFile.bytes += len(s)
	_, err := loadFile.f.WriteString(s)
	if err != nil {
		glog.Error("logFile: Write data file failed: %s\n", err.Error())
	}
}
