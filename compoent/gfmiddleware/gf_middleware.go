package gfmiddleware

import (
	"github.com/gogf/gf/net/ghttp"
)

// MiddlewareCORS 跨域中间件
func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}


