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
  
## activeMQ使用
    1.提供activeMQ消息订阅功能,参考mqMessageReceive方法  中 //订阅消息接收处理 注释为消息处理,请编写时自行处理  订阅消息实例名为配置文件中的queue属性
     2.提供activeMQ消息发送功能,方法为mqMessageSend 接收[]byte类型的消息,并将其推送到activeMQ服务消息实例为queueResult中, queueResult为配置文件中的queueResult属性
    
## sql使用

移交至[gutil](https://github.com/swgloomy/gutil)
    
    
## 定时器watchFuncDir

移交至[gutil](https://github.com/swgloomy/gutil)
    
## 监视文件夹目录如发生任何修改,重新载入

移交至[gutil](https://github.com/swgloomy/gutil)
    
## 项目中资源文件

    项目中所有资源文件存放在content文件夹中,该目录读取由go进行,文件资源暂不做缓存处理.在main.go文件中router路由方法g.GET("/assets/*pth", assetsFiles)
    请求路劲为http://域名/assets/资源文件目录  如content文件下有个js文件则为  http://域名/assets/jquery.js
    
    
## 项目添加路由
    
    请在main.go文件中router方法的
        {
        		g.GET("/", func(c *gin.Context) { c.String(200, "ok") })
        
        		g.GET("/assets/*pth", assetsFiles)
        	}
    处添加需要配置的路由项
    
## 文件读取
    
移交至[gutil](https://github.com/swgloomy/gutil)

## 数据写入文件并定时处理

移交至[gutil](https://github.com/swgloomy/gutil)

## 文件上传

   增加文件上传接口,接口 /unitUpLoadFile 需提供fname参数文件名称  上传文件需存放在form表单的file变量中,否则无法正确上传
   文件上传支持文件切片上传,原理类似断点续传
   
## 配置文件解析
      
      "rootPrefix": "",//二级目录地址
      "tempDir": "./template/*",//模版目录位置
      "contentDir": "content",//资源文件目录位置
      "dbuser": "",//数据库账号
      "dbhost": "",//数据库地址
      "dbport": 3306,//数据库端口
      "dbpass": "",//数据库密码
      "dbname": "",//数据库库名
      "redisProto":"tcp",//redis连接方式
      "redisAddr":"127.0.0.1:6379",//redis连接地址
      "redisDatabase":5,//redis  database,
      "port":":8000",//服务监听端口
      "mqAddr":"10.10.188.10:61613",//activeMQ 地址和端口
      "queueResult":"qqwe", //activeMQ发送实例名称
      "queue":"qweqwe",//activeMQ 持续接收实例名称
      "loadFileDir":"./dataFile" //数据写入文件的所在目录
      "logsDir": "./logs" //日志目录文件夹
      
      
## 开源协议

无使用限制