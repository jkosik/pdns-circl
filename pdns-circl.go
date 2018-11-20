package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type PdnsRecord struct {
	Count     int    `json:"count"`
	Origin    string `json:"origin"`
	TimeFirst int    `json:"time_first"`
	RRType    string `json:"rrtype"`
	RRName    string `json:"rrname"`
	RData     string `json:"rdata"`
	TimeLast  int    `json:"time_last"`
}

var username = flag.String("u", "foo", "Username")
var password = flag.String("p", "bar", "Password")
var rrName = flag.String("rrname", "www.circl.lu", "Domain to lookup, e.g. www.google.com")
var interactive = flag.Bool("i", false, "Interactive Mode")

func callAPI(url string) []byte {
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

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		panic(readErr.Error())
	}
	return body
}

func main() {

	flag.Parse()
	var urlConcat bytes.Buffer
	urlConcat.WriteString("https://www.circl.lu/pdns/query/")
	urlConcat.WriteString(*rrName)
	url := urlConcat.String()

	body := callAPI(url)

	//API response is a list of JSONs not one JSON object. Prior unmarshaling it is needed to fix JSON format
	s := string(body[:]) //byte array to string

	dataRaw := strings.Split(s, "\n") //slice from multiline string
	data := dataRaw[:len(dataRaw)-1]  //delete last empty line

	//Adding JSON delimiters and comma separation
	var buffer bytes.Buffer
	buffer.WriteString("[") //starting JSON char
	for i, rec := range data {
		buffer.WriteString(rec) //trailing JSON char
		if i < len(data)-1 {    //add comma delimiter except last iteration
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")     //trailing JSON char
	jsonData := buffer.String() //final valid JSON string ready for unmarshaling
	//fmt.Println("jsonData:", jsonData)
	var records []PdnsRecord
	err := json.Unmarshal([]byte(jsonData), &records)
	if err != nil {
		panic(err.Error())
	}

	if *interactive == false {
		fmt.Println(jsonData)
	} else {
		fmt.Println("Running interactive mode for", url)
       		fmt.Print("Enter First String: ")   //Print function is used to display output in same line
		var rrtype string    
	        fmt.Scanln(&rrtype) 
		fmt.Println("rrtype:",rrtype)
	}

	//	fmt.Println(records[0].Count)

}
