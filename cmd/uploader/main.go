package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"protocall/infrastructure/storage"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Config struct {
	host string
	port string
	root string
}

type Server struct {
	storage *storage.Storage
	conf    *Config
}

func (s *Server) upload(ctx *fasthttp.RequestCtx) {
	filename := string(ctx.QueryArgs().Peek("filename"))

	f, err := os.Open(s.conf.root + "/" + filename)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if err := s.storage.UploadRecord(context.TODO(), filename, f); err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func main() {
	s3, err := storage.NewStorage(&storage.Config{
		Endpoint: "https://storage.yandexcloud.net/",
		Bucket: "protocall",
	})
	if err != nil {
		fmt.Println("cannot connect to s3")
		os.Exit(1)
	}

	srv := Server{
		storage: s3,
		conf: &Config{
			host: "0.0.0.0",
			port: "8888",
			root: "/Users/murmur/protocall-connector/test_dir",
		},
	}

	r := router.New()
	r.POST("/upload", srv.upload)

	if err := fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", srv.conf.host, srv.conf.port), r.Handler); err != nil {
		fmt.Println(err)
	}
}
