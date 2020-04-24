# 配置文件使用  
该项目的目的是方便后续使用到初始化 配置文件的问题上面的使用 
方便配置文件的管理和使用   

如何使用呢？  
项目下面有 resource 目录   
app.conf 指定了使用哪个环境变量   
```
resource 目录下包含 app.conf 和 app-xxx.conf
app.conf  下面的内容 确保使用了哪个配置环境  
[app]  
app=test  
```
// how to use
// parse := NewConfigParse()  
// iniconf = parse.GetConfig()  