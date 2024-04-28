package Providers

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const BASEURL_PROXYSCRAPE = "https://api.proxyscrape.com/v2/?request=displayproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all"

type httpsServers []string

type httpServers []string

type eliteServers []string

type anonymousServers []string

type httpEliteServers []string

type httpAnonymousServers []string

type httpsEliteServers []string

type httpsAnonymousServers []string

func fetchHttpsServers(client *http.Client) httpsServers {
	serverList := requestServerList("all", "yes", client)
	serv := append(httpsServers{}, serverList...)
	return serv
}

func fetchHttpServers(client *http.Client) httpServers {
	serverList := requestServerList("all", "no", client)
	serv := append(httpServers{}, serverList...)
	return serv
}

func fetchEliteServers(client *http.Client) eliteServers {
	serverList := requestServerList("elite", "all", client)
	serv := append(eliteServers{}, serverList...)
	return serv
}

func fetchAnonymousServers(client *http.Client) anonymousServers {
	serverList := requestServerList("anonymous", "all", client)
	serv := append(anonymousServers{}, serverList...)
	return serv
}

func fetchHttpEliteServers(client *http.Client) httpEliteServers {
	serverList := requestServerList("elite", "no", client)
	serv := append(httpEliteServers{}, serverList...)
	return serv
}

func fetchHttpsEliteServers(client *http.Client) httpsEliteServers {
	serverList := requestServerList("elite", "yes", client)
	serv := append(httpsEliteServers{}, serverList...)
	return serv
}

func fetchHttpAnonymousServers(client *http.Client) httpAnonymousServers {
	serverList := requestServerList("anonymous", "no", client)
	serv := append(httpAnonymousServers{}, serverList...)
	return serv
}

func fetchHttpsAnonymousServers(client *http.Client) httpsAnonymousServers {
	serverList := requestServerList("anonymous", "yes", client)
	serv := append(httpsAnonymousServers{}, serverList...)
	return serv
}

func requestServerList(anonLevel string, ssl string, client *http.Client) []string {
	urlA, err := url.Parse(BASEURL_PROXYSCRAPE)
	if err != nil {
		log.Fatal(err)
	}
	values := urlA.Query()

	values.Set("anonymity", anonLevel)
	values.Set("ssl", ssl)
	urlA.RawQuery = values.Encode()

	return doRequest(urlA.String(), client)
}

func doRequest(url string, client *http.Client) []string {
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
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
	bodyString := string(body)
	split := strings.Split(bodyString, "\n")

	return split
}
