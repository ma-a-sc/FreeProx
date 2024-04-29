package FreeProx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const BASEURL = "https://proxylist.geonode.com/api/proxy-list?limit=500&page=1&sort_by=lastChecked&sort_type=desc"
const LIMIT = "limit"
const PAGE = "page"

type Server struct {
	ID                 string   `json:"_id"`
	IP                 string   `json:"ip"`
	AnonLevel          string   `json:"anonymityLevel"`
	Asn                string   `json:"asn"`
	City               string   `json:"city"`
	Country            string   `json:"country"`
	CreatedAt          string   `json:"created_at"`
	Google             bool     `json:"google"`
	Isp                string   `json:"isp"`
	LastCheck          int64    `json:"lastChecked"`
	Latency            float64  `json:"latency"`
	Org                string   `json:"org"`
	Port               string   `json:"port"`
	Protocols          []string `json:"protocols"`
	Region             string   `json:"region"`
	ResponseTime       uint     `json:"responseTime"`
	Speed              uint     `json:"speed"`
	UpdatedAt          string   `json:"updated_at"`
	WorkingPercent     string   `json:"workingPercent"`
	UpTime             float64  `json:"upTime"`
	UpTimeSuccessCount uint     `json:"upTimeSuccessCount"`
	UpTimeTryCount     uint     `json:"upTimeTryCount"`
}

type ResponseData struct {
	Data  []Server `json:"data"`
	Total uint     `json:"total"`
	Page  uint     `json:"page"`
	Limit uint     `json:"limit"`
}

type AllServer struct {
	Servers   []Server
	CreatedAt time.Time
}

func (a *AllServer) Len() int { return len(a.Servers) }

func (a *AllServer) OrderByAnonLevel(top string) {

}

func (a *AllServer) FilterByAnonLevel(anonLevel string) {

}

func (a *AllServer) OrderByLatency(decending bool) {

}

func (a *AllServer) FilterByLatency(keepUntil float64) {

}

func (a *AllServer) OrderByUpdateDate(decending bool) {

}

func (a *AllServer) OrderByUptime(keepUntil float64) {}

func GetAllProxies() AllServer {
	client := &http.Client{}

	initialRequestData := makeRequest(1, client)
	numberOfRemainingPages := int(math.Abs(float64(initialRequestData.Total / 500)))

	a := AllServer{
		Servers:   initialRequestData.Data,
		CreatedAt: time.Now(),
	}

	for page := 2; page <= numberOfRemainingPages; page++ {
		time.Sleep(500 * time.Millisecond)
		requestData := makeRequest(page, client)

		a.Servers = append(a.Servers, requestData.Data...)
	}

	return a
}

func makeRequest(page int, client *http.Client) ResponseData {
	geoNodeUrl := getUrlPage(page)
	method := "GET"

	req, err := http.NewRequest(method, geoNodeUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, clientErr := client.Do(req)
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	defer res.Body.Close()

	body, ioErr := io.ReadAll(res.Body)
	if ioErr != nil {
		log.Fatal(ioErr)
	}
	var r ResponseData

	jsonErr := json.Unmarshal([]byte(string(body)), &r)
	if jsonErr != nil {
		fmt.Println(geoNodeUrl)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}

		log.Fatal(jsonErr)
	}
	return r
}

func getUrlPage(page int) string {
	urlA, err := url.Parse(BASEURL)
	if err != nil {
		print(5)
		log.Fatal(err)
	}
	values := urlA.Query()

	values.Set(PAGE, strconv.Itoa(page))
	urlA.RawQuery = values.Encode()

	return urlA.String()
}
