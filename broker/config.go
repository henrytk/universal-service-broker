package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pivotal-cf/brokerapi"
)

type Config struct {
	API          API
	Catalog      Catalog
	ProviderData ProviderData
}

type API struct {
	BasicAuthUsername string `json:"basic_auth_username"`
	BasicAuthPassword string `json:"basic_auth_password"`
}

type Catalog struct {
	Catalog brokerapi.CatalogResponse `json:"catalog"`
}

type ProviderData struct {
	ProviderCatalog ProviderCatalog `json:"catalog"`
}

type ProviderCatalog struct {
	Services []ProviderService `json:"services"`
}

type ProviderService struct {
	ID             string          `json:"id"`
	ProviderConfig json.RawMessage `json:"provider_config"`
	Plans          []ProviderPlan  `json:"plans"`
}

type ProviderPlan struct {
	ID             string          `json:"id"`
	ProviderConfig json.RawMessage `json:"provider_config"`
}

func NewConfig(source io.Reader) (Config, error) {
	config := Config{}
	bytes, err := ioutil.ReadAll(source)
	if err != nil {
		return config, err
	}

	api := API{}
	if err = json.Unmarshal(bytes, &api); err != nil {
		return config, err
	}

	catalog := Catalog{}
	if err = json.Unmarshal(bytes, &catalog); err != nil {
		return config, err
	}

	providerData := ProviderData{}
	if err = json.Unmarshal(bytes, &providerData); err != nil {
		return config, err
	}

	return Config{
		API:          api,
		Catalog:      catalog,
		ProviderData: providerData,
	}, nil
}

func (c Config) Validate() error {
	if c.API.BasicAuthUsername == "" {
		return fmt.Errorf("Config error: basic auth username required")
	}
	if c.API.BasicAuthPassword == "" {
		return fmt.Errorf("Config error: basic auth password required")
	}
	return nil
}
