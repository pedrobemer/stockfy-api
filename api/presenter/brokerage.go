package presenter

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
