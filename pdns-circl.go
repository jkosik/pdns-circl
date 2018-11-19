package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Pdns struct {
	Count      int    `json:"count"`
	Origin     string `json:"origin"`
	TimeFirst int    `json:"time_first"`
	RRType     string `json:"rrtype"`
	RRName     string `json:"rrname"`
	RData      string `json:"rdata"`
	TimeLast  int    `json:"time_last"`
}

type Record struct {
	Record	string
}


===========================
var input = `
{"foo": "bar", "aaa":123}
{"foo": "baz", "aaa":222}
`

type Doc struct {
	Foo string
	Aaa int
}

func main() {
	dec := json.NewDecoder(strings.NewReader(input))
	for {
		var doc Doc

		err := dec.Decode(&doc)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", doc)
	}
}


================




func main() {

	username := flag.String("u", "foo", "Username")
	password := flag.String("p", "bar", "Password")
	rrName := flag.String("rrname", "www.circl.lu", "Domain to lookup, e.g. www.google.com")
	//	rData := flag.String("rdata", "cpa.circl.lu", "Data in response")
	flag.Parse()

	var urlConcat bytes.Buffer
	urlConcat.WriteString("https://www.circl.lu/pdns/query/")
	urlConcat.WriteString(*rrName)
	url := urlConcat.String()
	fmt.Println(url)
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		panic(reqErr.Error())
	}

	req.SetBasicAuth(*username, *password)
	req.Header.Set("User-Agent", "pdns-circl-golang-client")

	res, resErr := netClient.Do(req)
	if resErr != nil {
		panic(resErr.Error())
	}

	body, readErr := io.ioutil.ReadAll(res.Body)
	if readErr != nil {
		panic(readErr.Error())
	}
	fmt.Println(body)

	//content := Pdns{}
	//content := map[string]Pdns{}
	content := Record{}
	fmt.Println("before")
	fmt.Println(content)
	unmarshalErr := json.Unmarshal(body, &content)
	if unmarshalErr != nil {
		panic(unmarshalErr.Error())
	}
	fmt.Println("after")
	fmt.Println(content)

	marshal, marshalErr := json.MarshalIndent(content, "", "  ")
	if marshalErr != nil {
		panic(marshalErr.Error())
	}
	fmt.Println(string(marshal))
}

/*
//http://blog.josephmisiti.com/parsing-json-responses-in-golang
func getCapec(body []byte) (*Capec, error) {
	var s = new(Capec)
	err := json.Unmarshal(body, &s)
	if(err != nil){
	        fmt.Println("whoops:", err)
	}
	return s, err
}

*/
