package finnhub

import (
	"stockfyApi/client"
)

func CompanyProfile2Finnhub(symbol string) CompanyProfile2 {
	url := "https://finnhub.io/api/v1/stock/profile2?symbol=" + symbol +
		"&token=c2o3062ad3ie71thpra0"

	var companyProfile2 CompanyProfile2

	client.RequestAndAssignToBody(url, &companyProfile2)

	return companyProfile2
}
