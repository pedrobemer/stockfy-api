package alphaVantage

import "stockfyApi/client"

func (a *AlphaApi) CompanyOverviewAlpha(symbol string) CompanyOverview {
	url := "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" +
		symbol + "&apikey=" + a.Token

	var companyOverview CompanyOverview

	client.RequestAndAssignToBody("GET", url, nil, &companyOverview)

	return companyOverview
}
