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

	"protocall/domain/repository"
	"protocall/infrastructure/storage"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Root string `yaml:"root"`
}

type Server struct {
	storage repository.VoiceStorage
	bucket  string
	root    string
}

func NewServer(s3 repository.VoiceStorage, root, bucket string) *Server {
	return &Server{
		storage: s3,
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
	localFile := filepath.Join(s.root, from)

	to := string(ctx.QueryArgs().Peek("to"))
	if to == "" {
		fmt.Println("empty query parameter: to")
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	if err := s.storage.UploadFile(context.Background(), s.bucket, localFile, filepath.Join(to)); err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if err := os.Remove(localFile); err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		fmt.Println(err)
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
		SrvConf ServerConfig          `yaml:"uploader"`
		S3Conf  storage.StorageConfig `yaml:"s3"`
	}{
		SrvConf: ServerConfig{},
		S3Conf:  storage.StorageConfig{},
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.S3Conf.AccessKey = os.Getenv("ACCESS_KEY")
	config.S3Conf.SecretKey = os.Getenv("SECRET_KEY")

	s3, err := storage.NewStorage(&config.S3Conf)
	if err != nil {
		fmt.Println("cannot connect to s3", err)
		os.Exit(1)
	}

	srv := NewServer(s3, config.SrvConf.Root, config.S3Conf.Bucket)

	r := router.New()
	r.POST("/upload", srv.upload)

	if err := fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", config.SrvConf.Host, config.SrvConf.Port), r.Handler); err != nil {
		fmt.Println(err)
	}
}
