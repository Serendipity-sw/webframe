#web框架

## 程序须知
    
    1.如发现程序有不可用引用,请使用go get进行添加
    
    2.目前程序使用数据库为mysql数据库
    
    3.所有配置请参考config.json配置文件,配置文件中统一各配置注释名称
    
    
## GO版本
    1.6
    
    
## redis使用
    
    1.提供setRedisCachePs设置缓存,输入参数: uuid键值  ps存入字符串  输出参数: 错误对象
    
    2.提供getRedisCachePs获取缓存,输入参数: key 键值    输出参数: redis获取值 错误对象
    
## sql使用

    1.提供sqlSelect数据库查询方法 输入参数: sqlstr 要执行的sql语句 param执行SQL的语句参数化传递   输出参数: 查询返回条数  错误对象输出
    
    2.提供sqlExec数据库增删改方法 输入参数: sqlstr 要执行的sql语句  param执行SQL的语句参数化传递  输出参数: 执行结果对象  错误对象输出
    
    
## 定时器watchFuncDir

    watchFuncDir方法寄存在timeMonitor.go文件中,方法中注释处可添加所有需要定时执行的方法
    
## 监视文件夹目录如发生任何修改,重新载入

    增加notifyTemplates 方法,检测文件夹是否发生变化,如发生变化,则重新将模版文件载入到gin对象中
    
## 项目添加路由
    
    请在main.go文件中router方法的
        {
        		g.GET("/", func(c *gin.Context) { c.String(200, "ok") })
        
        		g.GET("/assets/*pth", assetsFiles)
        	}
    处添加需要配置的路由项