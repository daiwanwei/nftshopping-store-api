package items

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ContractManagerService interface {
	DeployCreation(request DeployCreationRequest) (err error)
}

type contractManagerService struct {
	client *http.Client
	domain string
}

func NewContractManagerService(client *http.Client, domain string) ContractManagerService {
	return &contractManagerService{
		client: client,
		domain: domain,
	}
}

func (service *contractManagerService) DeployCreation(request DeployCreationRequest) (err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", service.domain+"deployCreation", bytes.NewBuffer(body))
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

type DeployCreationRequest struct {
	Contract string `json:"contract"`
}
