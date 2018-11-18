package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Property struct {
	LastUpdate string `json:"last_update"`
	Size       int    `json:"size"`
}

type DBInfo struct {
	Capec    Property `json:"capec"`
	Cpe      Property `json:"cpe"`
	CpeOther Property `json:"cpeOther"`
	Cves     Property `json:"cves"`
	Cwe      Property `json:"cwe"`
	Via4     Property `json:"via4"`
}

func main() {
	url := "https://cve.circl.lu/api/dbInfo"

//	var username string = "foo"
//	var passwd string = "bar"

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		panic(reqErr.Error())
	}

//	req.SetBasicAuth(username, passwd)
	req.Header.Set("User-Agent", "pdns-circl-golang-client")

	res, resErr := netClient.Do(req)
	if resErr != nil {
		panic(resErr.Error())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		panic(readErr.Error())
	}
//	fmt.Println(body)

//	content := DBInfo{}
	content := map[string]Property{}
	fmt.Println(content)
	unmarshalErr := json.Unmarshal(body, &content)
	if unmarshalErr != nil {
		panic(unmarshalErr.Error())
	}
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
