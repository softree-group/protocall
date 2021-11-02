package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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

	if err := s.storage.PutFile(context.Background(), filepath.Join(s.root, from), filepath.Join(to)); err != nil {
		fmt.Println(err)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.Response.SetStatusCode(http.StatusNoContent)
}

var (
	configPath = flag.String("f", "", "path to configuration file")
)

func main() {
	flag.Parse()
	if *configPath == "" {
		fmt.Println("need to specify path to config")
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("cannot read configuration: %v", err)
	}

	config := &struct {
		SrvConf ServerConfig     `yaml:"rising"`
		S3Conf  s3.StorageConfig `yaml:"s3"`
	}{
		SrvConf: ServerConfig{},
		S3Conf:  s3.StorageConfig{},
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		log.Fatalf("cannot parse configuration: %v", err)
	}
	config.S3Conf.AccessKey = os.Getenv("ACCESS_KEY")
	config.S3Conf.SecretKey = os.Getenv("SECRET_KEY")

	storage, err := s3.NewStorage(&config.S3Conf)
	if err != nil {
		log.Fatalf("cannot connect to s3: %v", err)
	}

	srv := NewServer(storage, config.SrvConf.Root, config.S3Conf.Bucket)

	r := router.New()
	r.POST("/upload", srv.upload)

	if err = fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", config.SrvConf.Host, config.SrvConf.Port), r.Handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
