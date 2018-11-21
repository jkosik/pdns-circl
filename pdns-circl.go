package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type PdnsRecord struct {
	Count     int    `json:"count"`
	Origin    string `json:"origin"`
	TimeFirst int64  `json:"time_first"` //int64 needed for unix time conversion
	RRType    string `json:"rrtype"`
	RRName    string `json:"rrname"`
	RData     string `json:"rdata"`
	TimeLast  int64  `json:"time_last"` //int64 needed for unix time conversion
}

var username = flag.String("u", "user", "CIRCL PDNS API Username")
var password = flag.String("p", "pass", "CIRCL PDNS API Password")
var rrName = flag.String("rrname", "www.circl.lu", "Domain to lookup, e.g. www.google.com")
var rrType = flag.String("rrtype", "nil", "RR as subfilter, e.g. A, CNAME, AAAA")
var raw = flag.Bool("raw", false, "Complete raw output for -rrname. Good option for jq processing and filtering. Ignores -rrtype flag.")
var records []PdnsRecord
var url string

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

func listRRType(rrtype string, mystruct []PdnsRecord) {
	fmt.Println("+++++ Listing CIRCL PDNS records for", *rrName, "+++++\n")
	for _, rec := range mystruct {
		if rec.RRType == rrtype {
			s := reflect.ValueOf(&rec).Elem()
			typeOfT := s.Type()
			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)
				if i == 1 {
					continue
				}
				if (i == 2) || (i == 6) { //when running on time fields
					timestamp := f.Interface().(int64) //assertion needed
					fmt.Printf("%s = %v\n", typeOfT.Field(i).Name, time.Unix(timestamp, 0))
				} else {
					fmt.Printf("%s = %v\n", typeOfT.Field(i).Name, f.Interface())
				}
			}
			fmt.Println("------------------------------------------")
		}
	}
}

func listAllRRType(rrtype string, mystruct []PdnsRecord) {
	fmt.Println("+++++ Listing CIRCL PDNS records for", *rrName, "+++++\n")
	for _, rec := range mystruct {
		s := reflect.ValueOf(&rec).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			if i == 1 {
				continue
			}
			if (i == 2) || (i == 6) { //when running on time fields
				timestamp := f.Interface().(int64) //assertion needed
				fmt.Printf("%s = %v\n", typeOfT.Field(i).Name, time.Unix(timestamp, 0))
			} else {
				fmt.Printf("%s = %v\n", typeOfT.Field(i).Name, f.Interface())
			}
		}
		fmt.Println("------------------------------------------")
	}
}

func main() {

	flag.Parse()
	var urlConcat bytes.Buffer
	urlConcat.WriteString("https://www.circl.lu/pdns/query/")
	urlConcat.WriteString(*rrName)
	url = urlConcat.String()

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
	err := json.Unmarshal([]byte(jsonData), &records)
	if err != nil {
		panic(err.Error())
	}

	if *raw == true {
		fmt.Println(jsonData)
	}

	if *raw == false {
		switch *rrType {
		case "A":
			listRRType("A", records)
		case "CNAME":
			listRRType("CNAME", records)
		case "AAAA":
			listRRType("AAAA", records)
		case "PTR":
			listRRType("PTR", records)
		case "SOA":
			listRRType("SOA", records)
		case "NS":
			listRRType("NS", records)
		case "SRV":
			listRRType("SRV", records)
		case "TXT":
			listRRType("TXT", records)
		default:
			listAllRRType("default", records)
		}
	}

}
