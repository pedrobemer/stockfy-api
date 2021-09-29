package alphaVantage

import "stockfyApi/client"

func (a *AlphaApi) CompanyOverview(symbol string) CompanyOverview {
	url := "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" +
		symbol + "&apikey=" + a.Token

	var companyOverview CompanyOverview

	client.RequestAndAssignToBody("GET", url, nil, &companyOverview)

	return companyOverview
}
