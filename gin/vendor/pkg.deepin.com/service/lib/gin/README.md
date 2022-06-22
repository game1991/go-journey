# gin

一些基础的功能函数集合 包括referer等

# referer.go
```
	eng:=gin.Default()
	eng.Use(Referer("*.chinauos.com"))
	eng.Any("",func(ctx *gin.Context){})
```


