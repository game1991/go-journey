package gin

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

//DefaultReferer 获取请求的referer信息
var DefaultReferer = func(ctx *gin.Context) string {
	return ctx.Request.Header.Get("referer")
}

//Referer 生成referer验证函数
func Referer(site ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(site) == 0 {
			return
		}
		referer := DefaultReferer(ctx)
		if referer == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		uri, err := url.Parse(referer)
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
		}
		for _, item := range site {
			if strings.HasPrefix(item, "*.") { //站点通配符
				item = strings.Replace(item, "*.", "", 1)
				if strings.HasSuffix(uri.Host, item) {
					return
				}
			} else {
				if item == uri.Host {
					return
				}
			}
		}
		ctx.AbortWithStatus(http.StatusForbidden)
	}
}
