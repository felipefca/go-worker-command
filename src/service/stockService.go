package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"worker/src/config"
	"worker/src/models"
)

func ProcessStockService(stock models.Stock) {
	// Criando WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		stockResponse, err := getNasdaqStockPrice(stock)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(stockResponse)
	}()

	// Aguardando a conclusão das goroutines
	wg.Wait()

}

// Função para obter a última cotação utilizando uma API externa
func getNasdaqStockPrice(stock models.Stock) (models.StockResponse, error) {

	baseURL := fmt.Sprintf("%s%s", config.UrlNasdaq, config.PatchNasdaqQuote)
	params := url.Values{}
	params.Add("symbol", fmt.Sprintf("%s", stock.Stock))

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Finnhub-Token", config.NasdaqToken)
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.StockResponse{}, err
	}

	var stockResponse models.StockResponse
	err = json.Unmarshal(body, &stockResponse)
	if err != nil {
		return models.StockResponse{}, err
	}

	return stockResponse, nil
}
