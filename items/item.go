package items

import (
	"net/http"
)

type Client struct {
	ContractManager ContractManagerService
	Factory         FactoryService
}

func NewClient(domain string) *Client {
	client := &http.Client{}
	return &Client{
		ContractManager: NewContractManagerService(client, domain+"contractManager/"),
		Factory:         NewFactoryService(client, domain+"factory/"),
	}
}
