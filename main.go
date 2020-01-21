package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"github.com/swgloomy/gutil"
	"html/template"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

var (
	configFn      = flag.String("config", "./config.json", "config file path") //配置文件地址
	debugFlag     = flag.Bool("d", false, "debug mode")                        //是否为调试模式
	rootPrefix    string                                                       //二级目录地址
	tempDir       string                                                       //模版目录
	contentDir    string                                                       //脚本目录
	rt            *gin.Engine
	mqAddr        string //activeMQ 地址和端口
	queueResult   string //activeMQ发送实例名称
	queue         string //activeMQ 持续接收实例名称
	loadFileDir   string //数据写入文件的所在目录
	upLoadFileDir string //文件上传目录
	autoMatedDir  string //自动化文件夹监控目录地址
	logsDir       string //日志文件夹
)

/**
主函数入口
创建人:邵炜
创建时间:2016年2月26日11:22:03
*/
func main() {

	if gutil.CheckPid(gutil.ProgramServicePIDPath) { //判断程序是否启动
		return
	}

	flag.Parse()

	serverRun(*configFn, *debugFlag)

	c := make(chan os.Signal, 1)
	gutil.WritePid(gutil.ProgramServicePIDPath)
	// 信号处理
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// 等待信号
	<-c

	serverExit()
	gutil.RmPidFile(gutil.ProgramServicePIDPath)
	os.Exit(0)
}

/**
服务启动
创建人:邵炜
创建时间:2016年2月26日11:22:16
输入参数: cfn(配置文件地址) debug(是否为调试模式)
*/
func serverRun(cfn string, debug bool) {
	config.ReadCfg(cfn)

	rootPrefix = strings.TrimSpace(config.GetStringMust("rootPrefix"))
	logsDir = strings.TrimSpace(config.GetStringMust("logsDir"))
	tempDir = strings.TrimSpace(config.GetStringMust("tempDir"))
	err := gutil.CreateFileProcess(tempDir)
	if err != nil {
		fmt.Printf("serverRun upLoadFile exists! path: %s err: %s \n", tempDir, err.Error())
		os.Exit(1)
		return
	}
	contentDir = strings.TrimSpace(config.GetStringMust("contentDir"))
	err = gutil.CreateFileProcess(contentDir)
	if err != nil {
		fmt.Printf("serverRun upLoadFile exists! path: %s err: %s \n", contentDir, err.Error())
		os.Exit(1)
		return
	}
	port := strings.TrimSpace(config.GetStringMust("port"))
	mqAddr = strings.TrimSpace(config.GetStringMust("mqAddr"))
	queueResult = strings.TrimSpace(config.GetStringMust("queueResult"))
	queue = strings.TrimSpace(config.GetStringMust("queue"))
	loadFileDir = strings.TrimSpace(config.GetStringMust("loadFileDir"))
	err = gutil.CreateFileProcess(loadFileDir)
	if err != nil {
		fmt.Printf("serverRun upLoadFile exists! path: %s err: %s \n", loadFileDir, err.Error())
		os.Exit(1)
		return
	}
	upLoadFileDir = strings.TrimSpace(config.GetStringMust("upLoadFileDir"))
	err = gutil.CreateFileProcess(upLoadFileDir)
	if err != nil {
		fmt.Printf("serverRun upLoadFile exists! path: %s err: %s \n", upLoadFileDir, err.Error())
		os.Exit(1)
		return
	}
	autoMatedDir = strings.TrimSpace(config.GetStringMust("autoMatedDir"))
	err = gutil.CreateFileProcess(autoMatedDir)
	if err != nil {
		fmt.Printf("serverRun upLoadFile exists! path: %s err: %s \n", autoMatedDir, err.Error())
		os.Exit(1)
		return
	}

	if len(rootPrefix) != 0 {
		if !strings.HasPrefix(rootPrefix, "/") {
			rootPrefix = "/" + rootPrefix
		}
		if strings.HasSuffix(rootPrefix, "/") {
			rootPrefix = rootPrefix[0 : len(rootPrefix)-1]
		}
	}

	gutil.LogInit(debug, logsDir)

	//初始化所有go文件中的init内方法
	deferinit.InitAll()
	glog.Info("init all module successfully \n")

	//设置多CPU运行
	runtime.GOMAXPROCS(runtime.NumCPU())
	glog.Info("set many cpu successfully \n")

	//启动所有go文件中的init方法
	deferinit.RunRoutines()
	glog.Info("init all run successfully \n")

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	rt = gin.Default()
	loadTemplates(rt)
	router(rt)

	go rt.Run(port)
}

/**
加载模版
创建人:邵炜
创建时间:2016年2月26日11:34:12
输入参数: gin对象
*/
func loadTemplates(e *gin.Engine) {
	t, err := template.New("tmpls").Funcs(template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeurl": func(s string) template.URL {
			return template.URL(s)
		},
	}).ParseGlob(tempDir + "*")

	if err != nil {
		glog.Error("loadTemplates failed: %s %s \n", tempDir, err.Error())
		return
	}

	e.SetHTMLTemplate(t)
}

/**
结束进程
创建人:邵炜
创建时间:2016年3月7日14:21:24
*/
func serverExit() {
	// 结束所有go routine
	deferinit.StopRoutines()
	glog.Info("stop routine successfully.\n")

	deferinit.FiniAll()
	glog.Info("finish all modules successfully.\n")

	glog.Close()
}

/**
路由配置
创建人:邵炜
创建时间:2016年3月7日10:15:20
输入参数: gin对象
*/
func router(r *gin.Engine) {
	g := &r.RouterGroup
	if rootPrefix != "" {
		g = r.Group(rootPrefix)
	}

	{
		g.GET("/", func(c *gin.Context) { c.String(200, "ok") })
		g.Static("/assets", contentDir)
		g.GET("/template/:name", unitGetHtml)
		g.POST("/unitUpLoadFile", unitUploadFile) //文件上传
	}
}
