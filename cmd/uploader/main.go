package main

import (
	"context"
<<<<<<< HEAD
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
=======
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
>>>>>>> 9dff4d40660aa86a5ea66aa72ef8bfb947260101
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
<<<<<<< HEAD
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
=======
}

func main() {
	s3, err := storage.NewStorage(&storage.Config{
		Endpoint: "https://storage.yandexcloud.net/",
		Bucket:   "protocall",
>>>>>>> 9dff4d40660aa86a5ea66aa72ef8bfb947260101
	})
	if err != nil {
		fmt.Println("cannot connect to s3")
		os.Exit(1)
	}

<<<<<<< HEAD
	srv := NewServer(s3, config)
=======
	srv := Server{
		storage: s3,
		conf: &Config{
			host: "0.0.0.0",
			port: "8888",
			root: "/Users/murmur/protocall-connector/test_dir",
		},
	}
>>>>>>> 9dff4d40660aa86a5ea66aa72ef8bfb947260101

	r := router.New()
	r.POST("/upload", srv.upload)

<<<<<<< HEAD
	if err := fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", config.Host, config.Port), r.Handler); err != nil {
=======
	if err := fasthttp.ListenAndServe(fmt.Sprintf("%v:%v", srv.conf.host, srv.conf.port), r.Handler); err != nil {
>>>>>>> 9dff4d40660aa86a5ea66aa72ef8bfb947260101
		fmt.Println(err)
	}
}
