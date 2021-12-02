package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", quote)
	http.HandleFunc("/history/", quoteHistory)
	http.ListenAndServe(":8000", nil)
}

type Quote struct {
	Symbol          string  `json:"symbol"`
	QuoteType       string  `json:"quoteType"`
	DisplayName     string  `json: "displayName"`
	Bid             float64 `json: "bid"`
	Ask             float64 `json: "ask"`
	QuoteSourceName string  `json: "quoteSourceName"`
	Currency        string  `json: "currency"`
	TrailingPE      float64 `json: "trailingPE"`
	ForwardPE       float64 `json: "forwardPE"`
}

type QuoteHistory struct {
	Timestamp []int64
	Symbol    string
	Close     []float64
}

type QuoteResponse struct {
	Quotes []Quote `json:"result"`
}

type ResponseContainer struct {
	QuoteResponse QuoteResponse `json:"quoteResponse"`
}

func quoteHistory(w http.ResponseWriter, r *http.Request) {
	print(r.URL.Query().Encode())
	ticker := strings.Split(r.URL.Path, "/")[2]
	client := &http.Client{}

	request, _ := http.NewRequest(http.MethodGet, "https://yfapi.net/v8/finance/spark", nil)
	query := request.URL.Query()
	query.Add("symbols", ticker)
	query.Add("interval", "1d")
	query.Add("range", "1mo")
	query.Add("region", "US")
	query.Add("lang", "en")
	request.URL.RawQuery = query.Encode()
	request.Header.Add("x-api-key", "qPMkSbpVIy2QelTfQ9HPp8ljVrjgt5mT7UCfIzso")

	res, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var m map[string]QuoteHistory

	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(m[ticker])
	serializedQuoteHistory, _ := json.Marshal(m[ticker])
	fmt.Fprintf(w, string(serializedQuoteHistory))
}

func quote(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Path[1:]
	client := &http.Client{}

	request, _ := http.NewRequest(http.MethodGet, "https://yfapi.net/v6/finance/quote", nil)
	query := request.URL.Query()
	query.Add("symbols", ticker)
	query.Add("region", "US")
	query.Add("lang", "en")
	request.URL.RawQuery = query.Encode()
	request.Header.Add("x-api-key", "qPMkSbpVIy2QelTfQ9HPp8ljVrjgt5mT7UCfIzso")

	res, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(body))
	var responseContainer ResponseContainer
	if err := json.Unmarshal(body, &responseContainer); err != nil {
		log.Fatal("response unmarshalling error", err)
		return
	}
	if len(responseContainer.QuoteResponse.Quotes) > 0 {
		serializedQuote, _ := json.Marshal(responseContainer.QuoteResponse.Quotes[0])
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		fmt.Fprintf(w, string(serializedQuote))
		return
	}

	fmt.Fprintf(w, "No Quote Found")
}
