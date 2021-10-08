package presenter

import "stockfyApi/entity"

type Brokerage struct {
	Id      string `db:"id" json:",omitempty"`
	Name    string `db:"name" json:",omitempty"`
	Country string `db:"country" json:",omitempty"`
}

func ConvertBrokerageToApiReturn(id string, name string, country string) *Brokerage {
	return &Brokerage{
		Id:      id,
		Name:    name,
		Country: country,
	}
}

func ConvertArrayBrokerageToApiReturn(brokerageFirms []entity.Brokerage) []Brokerage {
	var brokerageFirmsConverted []Brokerage

	for _, brokerage := range brokerageFirms {
		brokerageConverted := ConvertBrokerageToApiReturn(brokerage.Id,
			brokerage.Name, brokerage.Country)
		brokerageFirmsConverted = append(brokerageFirmsConverted,
			*brokerageConverted)
	}

	return brokerageFirmsConverted
}
