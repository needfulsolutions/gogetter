package main

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var url = "https://www.reddit.com/r/Unexpected/comments/y679qs/uhoh/.json"
var useragent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36"

func main() {

	// This is required because go automatically enables a
	// HTTPS/2 that is broken on some sites, unfortunately including reddit.
	// A workaround that does not require this (meaning the monstrosity
	// in the Transport field) is running the binary as such:
	// $ env GODEBUG=http2client=0 ./gogetter
	client := http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", useragent)
	req.Header.Set("Host", "www.reddit.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "Keep-Alive")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(resp.Header["Content-Encoding"][0])

	content_encoding_exists := (len(resp.Header["Content-Encoding"]) >= 1)

	if !content_encoding_exists || resp.Header["Content-Encoding"][0] != "gzip" {
		fmt.Println("Unknown encoding or none.")
		fmt.Println(string(body))
		return
	}

	buf := bytes.NewBuffer(body)

	gzip_reader, err := gzip.NewReader(buf)
	if err != nil {
		log.Fatal(err)
	}

	defer gzip_reader.Close()

	decomp, err := ioutil.ReadAll(gzip_reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Decompressed data:")
	fmt.Println(string(decomp))
}
