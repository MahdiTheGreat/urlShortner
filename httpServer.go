package main

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (d data) Post(c echo.Context) (err error) {

	// To avoid security flaws try to avoid passing bound structs directly to other methods
	// if these structs contain fields that should not be bindable.
	client := d.client
	host := d.host
	exp := d.exp
	port := d.port
	url := c.FormValue("url")
	log.Println("received url in post is")
	log.Println(url)
	date := time.Now().Format("01-02-2006 15:04:05")
	code := 0
	flag := 0
	for flag == 0 {
		date = date + "" + strconv.Itoa(code)
		_, err := client.Get(date).Result()
		if err != nil {
			flag = +1
			print("flag updated")
		}
		println("inside loop")
		code += 1
	}
	date = strings.Replace(date, ":", "", 100)
	date = strings.Replace(date, " ", "", 100)
	print("key value is")
	print(date + " " + url)
	client.Set(date, url, time.Duration(exp)*time.Second)
	shortUrl := host + ":" + port + "/" + date
	return c.String(http.StatusOK, shortUrl)
}

func (d data) Get(c echo.Context) error {
	client := d.client
	shortUrl := c.Param("shortUrlVal")
	log.Println("received shortUrl in get is")
	log.Println(shortUrl)

	val, err := client.Get(shortUrl).Result()
	if err == redis.Nil {
		println("expired or no key set")
	}
	// nolint: wrapcheck
	println("val is")
	println(val)
	return c.String(http.StatusOK, val)
}

type data struct {
	client *redis.Client
	host   string
	port   string
	exp    int
}

func main() {
	e := echo.New()
	dbHost := os.Getenv("dbHost")
	//host, err := os.Hostname()
	//if err != nil {
	//	panic(err)
	//	host = "urlShortner"

	//}
	host := os.Getenv("host")
	port := os.Getenv("port")
	dbPort := os.Getenv("dbPort")
	pass := os.Getenv("pass")
	expTemp := os.Getenv("exp")
	println("dbHost is" + " " + dbHost)
	println("host is" + " " + host)
	println("dbPort is" + " " + dbPort)
	println("port is" + " " + port)
	println("pass is" + " " + pass)
	println("exp is" + " " + expTemp)
	exp, _ := strconv.Atoi(expTemp)
	client := redis.NewClient(&redis.Options{
		Addr:     dbHost + ":" + dbPort,
		Password: pass,
		DB:       0,
	})
	println(client.Ping().String())
	d := data{client: client, host: host, port: port, exp: exp}
	e.POST("", d.Post)
	e.GET("/:shortUrlVal", d.Get)
	add := "0.0.0.0:" + port
	println("add is" + add)
	if err := e.Start(add); err != nil {
		log.Fatal(err)
	}
}
