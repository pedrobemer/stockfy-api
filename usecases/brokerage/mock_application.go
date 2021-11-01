package brokerage

import (
	"errors"
	"stockfyApi/entity"
)

type MockApplication struct {
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (a *MockApplication) SearchBrokerage(searchType string, name string,
	country string) ([]entity.Brokerage, error) {
	var brokerageInfo []entity.Brokerage
	var err error

	if searchType == "COUNTRY" && country != "BR" && country != "US" &&
		country != "ERROR_BROKERAGE_SEARCH" {
		return nil, entity.ErrInvalidCountryCode
	}

	if searchType == "SINGLE" && name == "" {
		return nil, entity.ErrInvalidBrokerageNameSearchBlank
	}

	if name == "ERROR_BROKERAGE_SEARCH" || country == "ERROR_BROKERAGE_SEARCH" {
		err = errors.New("Unknown error in the brokerage repository")
	} else {
		switch searchType {
		case "ALL":
			brokerageInfo = []entity.Brokerage{
				{
					Id:       "TestBrokerageID1",
					Name:     "Test US 1",
					Fullname: "Test US 1",
					Country:  "US",
				},
				{
					Id:       "TestBrokerageID2",
					Name:     "Test US 2",
					Fullname: "Test US 2",
					Country:  "US",
				},
				{
					Id:       "TestBrokerageID3",
					Name:     "Test BR 1",
					Fullname: "Test BR 1",
					Country:  "BR",
				},
			}
			break
		case "SINGLE":
			if name != "UNKNOWN_BROKERAGE" {
				brokerageInfo = []entity.Brokerage{
					{
						Id:       "TestBrokerageID1",
						Name:     name,
						Fullname: "Test US 1",
						Country:  "US",
					},
				}
			}
			break
		case "COUNTRY":
			if country == "BR" || country == "US" {
				brokerageInfo = []entity.Brokerage{
					{
						Id:       "TestBrokerageID1",
						Name:     "Test " + country + " 1",
						Fullname: "Test " + country + " 1",
						Country:  country,
					},
					{
						Id:       "TestBrokerageID1",
						Name:     "Test " + country + " 2",
						Fullname: "Test " + country + " 2",
						Country:  country,
					},
				}
			}
			break
		default:
			return nil, entity.ErrInvalidBrokerageSearchType
		}

	}

	if err != nil {
		return nil, err
	}

	if brokerageInfo == nil {
		return nil, entity.ErrInvalidBrokerageNameSearch
	}

	return brokerageInfo, nil
}
