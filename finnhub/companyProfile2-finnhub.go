package finnhub

import (
	"stockfyApi/client"
)

func (f *FinnhubApi) CompanyProfile2Finnhub(symbol string) CompanyProfile2 {
	url := "https://finnhub.io/api/v1/stock/profile2?symbol=" + symbol +
		"&token=" + f.Token

	var companyProfile2 CompanyProfile2

	client.RequestAndAssignToBody("GET", url, nil, &companyProfile2)

	return companyProfile2
}
