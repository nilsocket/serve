package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var port = flag.Int("p", 0, "port")

func main() {
	flag.Parse()

	addr := ":"
	if *port != 0 {
		addr += strconv.Itoa(*port)
	}

	ip, err := getPublicIP()
	if err != nil {
		log.Println(err)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
	}
	defer ln.Close()

	if na, ok := ln.Addr().(*net.TCPAddr); ok {
		fmt.Println("http://" + ip + ":" + strconv.Itoa(na.Port))
	}

	err = http.Serve(ln, http.FileServerFS(os.DirFS(".")))
	if err != nil {
		log.Println(err)
	}
}

var publicIPURL = "https://cloudflare.com/cdn-cgi/trace"
var ipRegex = regexp.MustCompile("ip=(.*)")

func getPublicIP() (string, error) {
	resp, err := http.Get(publicIPURL)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := ipRegex.FindStringSubmatch(string(data))[1]

	if isIPV6(ip) {
		ip = "[" + ip + "]"
	}

	return ip, nil
}

func isIPV6(ip string) bool {
	return strings.Contains(ip, ":")
}
