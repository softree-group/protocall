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

type Config struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	LocalRoot  string `yaml:"root"`
	RemoteRoot string `yaml:"s3root"`
	Bucket     string `yaml:"bucket"`
	Endpoint   string `yaml:"endpoint"`
	DisableSSL string `yaml:"disableSSL"`
}

type Server struct {
	storage    repository.VoiceStorage
	bucket     string
	localRoot  string
	remoteRoot string
}

func NewServer(s3 repository.VoiceStorage, c *Config) *Server {
	return &Server{
		storage:    s3,
		bucket:     c.Bucket,
		localRoot:  c.LocalRoot,
		remoteRoot: c.RemoteRoot,
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

	if err := s.storage.UploadFile(
		context.Background(),
		s.bucket,
		filepath.Join(s.localRoot, from),
		filepath.Join(s.remoteRoot, to),
	); err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	ctx.Response.SetStatusCode(http.StatusOK)
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

	config := &Config{}
	if err := yaml.Unmarshal(data, struct {
		Config Config `yaml:"server"`
	}{Config: *config}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s3, err := storage.NewStorage(&storage.Config{
		Endpoint:   config.Endpoint,
		AccessKey:  os.Getenv("ACCESS_KEY"),
		SecretKey:  os.Getenv("SECRET_KEY"),
		DisableSSL: true,
	})
	if err != nil {
		fmt.Println("cannot connect to s3")
		os.Exit(1)
	}

	srv := NewServer(s3, config)

	r := router.New()
	r.POST("/upload", srv.upload)

	if err := fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", config.Host, config.Port), r.Handler); err != nil {
		fmt.Println(err)
	}
}
