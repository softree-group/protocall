package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"

	"protocall/pkg/s3"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Root string `yaml:"root"`
}

type Server struct {
	storage *s3.S3
	bucket  string
	root    string
}

func NewServer(s *s3.S3, root, bucket string) *Server {
	return &Server{
		storage: s,
		root:    root,
		bucket:  bucket,
	}
}

func (s *Server) upload(ctx *fasthttp.RequestCtx) {
	from := string(ctx.QueryArgs().Peek("from"))
	if from == "" {
		fmt.Println("empty query parameter: from")
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	to := string(ctx.QueryArgs().Peek("to"))
	if to == "" {
		fmt.Println("empty query parameter: to")
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	if err := s.storage.UploadFile(context.Background(), s.bucket, filepath.Join(s.root, from), filepath.Join(to)); err != nil {
		fmt.Println(err)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.Response.SetStatusCode(http.StatusNoContent)
}

func main() {
	configPath := ""
	flag.StringVar(&configPath, "f", "", "путь до файла конфигурации")

	flag.Parse()
	if configPath == "" {
		fmt.Println("need to specify path to config")
		flag.Usage()
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := &struct {
		SrvConf ServerConfig     `yaml:"uploader"`
		S3Conf  s3.StorageConfig `yaml:"s3"`
	}{
		SrvConf: ServerConfig{},
		S3Conf:  s3.StorageConfig{},
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.S3Conf.AccessKey = os.Getenv("ACCESS_KEY")
	config.S3Conf.SecretKey = os.Getenv("SECRET_KEY")

	storage, err := s3.NewStorage(&config.S3Conf)
	if err != nil {
		fmt.Println("cannot connect to s3", err)
		os.Exit(1)
	}

	srv := NewServer(storage, config.SrvConf.Root, config.S3Conf.Bucket)

	r := router.New()
	r.POST("/upload", srv.upload)

	if err = fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", config.SrvConf.Host, config.SrvConf.Port), r.Handler); err != nil {
		fmt.Println(err)
	}
}
