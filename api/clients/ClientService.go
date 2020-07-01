package clients

import (
	"fmt"
	"strings"

	"github.com/adampresley/mytime/api/helpers"
	"github.com/adampresley/simdb"
)

type ClientServicer interface {
	CreateClient(client Client) (int, error)
	ListClients(search ClientSearch) (ClientCollection, error)
	GetClientByCode(code string) (Client, error)
	GetClientByID(id int) (Client, error)
	UpdateClient(client Client) error
}

type ClientService struct {
	DB            *simdb.Driver
	HelperService helpers.HelperServicer
}

type ClientServiceConfig struct {
	DB            *simdb.Driver
	HelperService helpers.HelperServicer
}

type ClientSearch struct {
	Archived bool
	Name     string
}

func NewClientService(config ClientServiceConfig) ClientService {
	return ClientService{
		DB:            config.DB,
		HelperService: config.HelperService,
	}
}

func (s ClientService) CreateClient(client Client) (int, error) {
	var err error

	client.ClientID = s.DB.Open(Client{}).GetNextNumericID()

	if err = s.DB.Insert(client); err != nil {
		return 0, fmt.Errorf("Error creating new client: %w", err)
	}

	return client.ClientID, nil
}

func (s ClientService) ListClients(search ClientSearch) (ClientCollection, error) {
	var err error

	result := make(ClientCollection, 0, 10)
	d := s.DB.Open(Client{})

	if search.Archived {
		d = d.Where("archived", "=", true)
	} else {
		d = d.Where("archived", "=", false)
	}

	if err = d.Get().AsEntity(&result); err != nil {
		return result, fmt.Errorf("Error querying for Clients: %w", err)
	}

	if search.Name != "" {
		result = s.filterClientsByName(result, search.Name)
	}

	return result, nil
}

func (s ClientService) GetClientByCode(code string) (Client, error) {
	var err error
	var client Client

	if err = s.DB.Open(Client{}).Where("code", "=", code).First().AsEntity(&client); err != nil {
		return client, fmt.Errorf("Error finding client with a code '%s': %w", code, err)
	}

	return client, nil
}

func (s ClientService) GetClientByID(id int) (Client, error) {
	var err error
	var client Client

	if err = s.DB.Open(Client{}).Where("clientID", "=", id).First().AsEntity(&client); err != nil {
		return client, fmt.Errorf("Error finding client with an id'%s': %w", id, err)
	}

	return client, nil
}

func (s ClientService) UpdateClient(client Client) error {
	return s.DB.Open(Client{}).Update(client)
}

func (s ClientService) filterClientsByName(clients ClientCollection, clientName string) ClientCollection {
	result := make(ClientCollection, 0, 20)

	for _, c := range clients {
		lowerClientName := strings.ToLower(clientName)

		if strings.Contains(strings.ToLower(c.Name), lowerClientName) || strings.Contains(strings.ToLower(c.Code), lowerClientName) {
			result = append(result, c)
		}
	}

	return result
}
