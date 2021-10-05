package brokerage

import "stockfyApi/entity"

type Application struct {
	repo Repository
}

//NewApplication create new use case
func NewApplication(r Repository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) SearchBrokerage(searchType string, name string,
	country string) (*[]entity.Brokerage, error) {
	var brokerageInfo []entity.Brokerage
	var err error

	if searchType == "COUNTRY" && country != "BR" && country != "US" {
		return nil, entity.ErrInvalidCountryCode
	}

	switch searchType {
	case "ALL":
		brokerageInfo, err = a.repo.Search(searchType)
		break
	case "SINGLE":
		brokerageInfo, err = a.repo.Search(searchType, name)
		break
	case "COUNTRY":
		brokerageInfo, err = a.repo.Search(searchType, country)
		break
	default:
		return nil, entity.ErrInvalidBrokerageSearchType
	}

	if err != nil {
		return nil, err
	}

	return &brokerageInfo, err

}
