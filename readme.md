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

    增加notifyTemplates 方法,检测template文件夹是否发生变化,如发生变化,则重新将模版文件载入到gin对象中
    
## 项目中资源文件

    项目中所有资源文件存放在content文件夹中,该目录读取由go进行,文件资源暂不做缓存处理.在main.go文件中router路由方法g.GET("/assets/*pth", assetsFiles)
    请求路劲为http://域名/assets/资源文件目录  如content文件下有个js文件则为  http://域名/assets/jquery.js
    
## 增加图片流字符串解析保存到给定的目录中
    
    提供imageFileSave方法,该方法寄存在imageFile.go文件中,方法接收两个参数,其中对图片流格式要求为: data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGIAAABiCAYAAACrpQYOAAAAGX.........
    path为接收的需要保存路径,如 ./temp/
    方法会将图片流转化为实际图片存放在给定目录中.并返回路径
    
## 项目添加路由
    
    请在main.go文件中router方法的
        {
        		g.GET("/", func(c *gin.Context) { c.String(200, "ok") })
        
        		g.GET("/assets/*pth", assetsFiles)
        	}
    处添加需要配置的路由项
    
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
      "port":":8000"//服务监听端口