package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"worker/src/config"
	"worker/src/models"
)

func ProcessFiatService(fiat models.Fiat) {
	// Criando WaitGroup
	var wg sync.WaitGroup
	wg.Add(2)

	// Chamada das duas APIs externas em paralelo
	go func() {
		defer wg.Done()
		fiatResponse, err := getFiatCotation(fiat)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fiatResponse)
	}()

	go func() {
		defer wg.Done()
		lastFiatResponse, err := getLastFiatCotation(fiat)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(lastFiatResponse)
	}()

	// Aguardando a conclusão das goroutines
	wg.Wait()

}

// Função para obter a última cotação utilizando uma API externa
func getFiatCotation(fiat models.Fiat) (models.FiatResponse, error) {

	url := fmt.Sprintf("%s%s-%s", config.UrlLastContation, fiat.Code, fiat.CodeIn)

	response, err := http.Get(url)
	if err != nil {
		return models.FiatResponse{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.FiatResponse{}, err
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		return models.FiatResponse{}, err
	}

	dataFieldName := fmt.Sprintf("%s%s", fiat.Code, fiat.CodeIn)
	dataFieldValue, ok := objmap[dataFieldName]
	if !ok {
		return models.FiatResponse{}, fmt.Errorf("campo de resposta não encontrado")
	}

	var data models.FiatResponse
	err = json.Unmarshal(*dataFieldValue, &data)
	if err != nil {
		return models.FiatResponse{}, err
	}

	return data, nil
}

// Função para obter a cotação dos últimos dias utilizando uma API externa
func getLastFiatCotation(fiat models.Fiat) (models.FiatResponse, error) {

	url := fmt.Sprintf("%s%s-%s/%d", config.UrlLastDaysContation, fiat.Code, fiat.CodeIn, config.LastDaysContation)

	response, err := http.Get(url)
	if err != nil {
		return models.FiatResponse{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.FiatResponse{}, err
	}

	var fiatResponse []models.FiatResponse
	err = json.Unmarshal(body, &fiatResponse)
	if err != nil {
		return models.FiatResponse{}, err
	}

	lastFiat := fiatResponse[len(fiatResponse)-1]

	return lastFiat, nil
}
