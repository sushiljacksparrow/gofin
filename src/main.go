package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexPage)
	http.ListenAndServe(":8000", nil)
}

type Quote struct {
	Symbol string `json:"symbol"`
	QuoteType string `json:"quoteType"`
	DisplayName string `json: "displayName"`
	Bid float64 `json: "bid"`
	Ask float64 `json: "ask"`
}

type QuoteResponse struct {
	Quotes[] Quote `json:"result"`
}

type ResponseContainer struct {
	QuoteResponse QuoteResponse	`json:"quoteResponse"`
}

func indexPage(w http.ResponseWriter, r *http.Request) {
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
	defer  res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(body))
	var responseContainer ResponseContainer
	if err := json.Unmarshal(body, &responseContainer); err != nil {
		log.Fatal("response unmarshalling error", err)
		return
	}
	log.Print(len(responseContainer.QuoteResponse.Quotes))
	if len(responseContainer.QuoteResponse.Quotes) > 0 {
		serializedQuote, _ := json.Marshal(responseContainer.QuoteResponse.Quotes[0])
		fmt.Fprintf(w, string(serializedQuote))
		return
	}

	fmt.Fprintf(w,"No Quote Found")
}