package items

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FactoryService interface {
	OrderItem(request OrderItemRequest) (err error)
}

type factoryService struct {
	client *http.Client
	domain string
}

func NewFactoryService(client *http.Client, domain string) FactoryService {
	return &factoryService{
		client: client,
		domain: domain,
	}
}

func (service *factoryService) OrderItem(request OrderItemRequest) (err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", service.domain+"orderItem", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := service.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = StatusHandler(resp.StatusCode)
	if err != nil {
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var dataResp NonDataResponse
	err = json.Unmarshal(responseBody, &dataResp)
	if err != nil {
		return
	}
	err = ResponseHandler(dataResp)
	if err != nil {
		return
	}
	fmt.Println(dataResp)
	return
}

type OrderItemRequest struct {
	ProductId string `json:"productId"`
	Contract  string `json:"contract"`
	Amount    int    `json:"amount"`
}
