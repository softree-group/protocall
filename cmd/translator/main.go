package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type user struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Path     string `json:"path" binding:"required"`
}

type translatorReq struct {
	ConfID    string    `json:"conf_id" binding:"required"`
	StartTime time.Time `json:"start" binding:"required"`
	User      user
}

func test(ctx *fasthttp.RequestCtx) {
	fmt.Println("HERE")
	data := &translatorReq{}
	if err := json.Unmarshal(ctx.PostBody(), data); err != nil {
		fmt.Println(err)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}
	fmt.Println(data)
	ctx.Response.SetStatusCode(http.StatusOK)
}

func main() {
	r := router.New()
	r.POST("/translations", test)

	if err := fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", "127.0.0.1", 8181), r.Handler); err != nil {
		fmt.Println(err)
	}
}
