package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"protocall/internal/account"
	"protocall/pkg/logger"

	"github.com/gomodule/redigo/redis"
)

var (
	host     = flag.String("host", "", "set redis host")
	port     = flag.String("port", "", "set redis port")
	accounts = flag.String("file", "", "account file path")
)

func main() {
	flag.Parse()
	if *host == "" || *port == "" || *accounts == "" {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	jsonFile, err := os.Open(*accounts)
	if err != nil {
		logger.L.Fatal("fail to open file ", *accounts, ": ", err)
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var data []account.Account
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		logger.L.Fatal("fail to parse file ", *accounts, ": ", err)
	}

	store := account.NewStorage(conn)
	for idx, account := range data {
		store.SaveAccount(strconv.Itoa(idx), account)
	}
}
