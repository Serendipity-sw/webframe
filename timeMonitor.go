package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/howeyc/fsnotify"
	"github.com/smtc/glog"
	"html/template"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	jsTmr    *time.Timer
	funcName = template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeurl": func(s string) template.URL {
			return template.URL(s)
		},
	}
	autoMatedLock sync.RWMutex
	autoMatedFile map[string]int = make(map[string]int)
)

func init() {
	deferinit.AddInit(func() {
		tempDir = config.GetStringDefault("tempDir", "template/")
		if !strings.HasSuffix(tempDir, "/") {
			tempDir += "/"
		}
	}, nil, 40)
	deferinit.AddRoutine(notifyTemplates)
	deferinit.AddRoutine(watchFuncDir)
	deferinit.AddRoutine(watchAutoMate)
}

/**
定时运行程序
创建人:邵炜
创建时间:2016年3月7日09:51:42
输入参数: 终止命令  计数器对象
*/
func watchFuncDir(ch chan struct{}, wg *sync.WaitGroup) {
	go func() {
		<-ch

		jsTmr.Stop()
		wg.Done()
	}()

	jsTmr = time.NewTimer(time.Minute)
	for {
		//需要定时执行的方法
		jsTmr.Reset(time.Minute)
		<-jsTmr.C
	}
}

/**
加载模版
创建人:邵炜
创建时间:2016年2月26日11:34:12
输入参数: gin对象
*/
func loadTemplates(e *gin.Engine) {
	t, err := template.New("tmpls").Funcs(funcName).ParseGlob(tempDir + "*")

	if err != nil {
		glog.Error("loadTemplates failed: %s %s \n", tempDir, err.Error())
		return
	}

	e.SetHTMLTemplate(t)
}

/**
监视文件夹目录如发生任何修改,重新载入
创建人:邵炜
创建时间:2016年3月7日09:47:50
输入参数: 终止命令 计数器对象
*/
func notifyTemplates(ch chan struct{}, wg *sync.WaitGroup) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error("notifyTemplates: create new watcher failed: %v\n", err)
		return
	}

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				glog.Debug("notifyTemplates: event: %v\n", ev)
				loadTemplates(rt)
			case err := <-watcher.Error:
				glog.Error("notifyTemplates: error: %v\n", err)
			}
		}
	}()

	err = watcher.Watch(tempDir)
	if err != nil {
		glog.Error("notifyTemplates: watch dir %s failed: %v \n", tempDir, err)
	}

	// Hang so program doesn't exit
	<-ch

	/* ... do clean stuff ... */
	watcher.Close()
	wg.Done()
}

/**
自动化文件夹监控
创建人:邵炜
创建时间:2016年12月22日11:27:45
*/
func watchAutoMate(ch chan struct{}, wg *sync.WaitGroup) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error("watchAutoMate: fsnotify newWatcher is error! err: %s \n", err.Error())
		return
	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev == nil {
					continue
				}
				glog.Info("watchAutoMate: fsnotify watcher fileName: %s is change!  ev: %v \n", ev.Name, ev)
				if ev.IsDelete() {
					continue
				}
				autoMatedLock.RLock()
				_, ok := autoMatedFile[ev.Name]
				autoMatedLock.RUnlock()
				if !ok {
					autoMatedLock.Lock()
					autoMatedFile[ev.Name] = 0
					autoMatedLock.Unlock()
					go watchFileAutoMated(ev.Name, fileProcess)
				}
			case err := <-watcher.Error:
				if err == nil {
					continue
				}
				glog.Error("watchAutoMate: fsnotify watcher is error! err: %s \n", err.Error())
			}
		}
		done <- true
	}()
	err = watcher.WatchFlags(autoMatedDir, fsnotify.FSN_MODIFY)
	if err != nil {
		glog.Error("watchAutoMate watch error. userListLoadDir: %s  err: %s \n", autoMatedDir, err.Error())
	}

	// Hang so program doesn't exit
	<-ch

	/* ... do stuff ... */
	watcher.Close()
	wg.Done()
}

/**
自动化创建任务 需要监控的文件,判断文件是否上传完毕
创建人:邵炜
创建时间:2016年9月5日14:58:13
输入参数: 文件路劲
*/
func watchFileAutoMated(filePath string, callBack func(string)) {
	defer func() {
		autoMatedLock.Lock()
		delete(autoMatedFile, filePath)
		autoMatedLock.Unlock()
	}()
	tmrIntal := 10 * time.Second
	fileSaveTmr := time.NewTimer(tmrIntal)
	fileState, err := os.Stat(filePath)
	if err != nil {
		glog.Error("watchFileAutoMatedTask can't load file! path: %s err: %s \n", filePath, err.Error())
		return
	}
	var (
		size   = fileState.Size()
		number int64
	)
	<-fileSaveTmr.C
	for {
		fileState, err = os.Stat(filePath)
		if err != nil {
			glog.Error("watchFileAutoMatedTask can't load file! path: %s err: %s \n", filePath, err.Error())
			return
		}
		number = fileState.Size()
		if size == number {
			go callBack(filePath)
			return
		}
		size = number
		fileSaveTmr.Reset(tmrIntal)
		<-fileSaveTmr.C
	}
}

/**
文件处理
创建人:邵炜
创建时间:2016年9月5日15:04:18
输入参数:文件路劲
*/
func fileProcess(filePathStr string) {
	fileNameIn := getMyFileName(filePathStr)
	fileName := *fileNameIn
	if strings.HasPrefix(fileName[0], "picture-") {
		err := fileCreateAndWriteByPath(filePathStr, fmt.Sprintf("./%s", strings.Join(fileName, "")))
		if err != nil {
			glog.Error("fileProcess run err! \n")
			return
		}
		glog.Info("picture file move run success! \n")
	}
}
