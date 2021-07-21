package alphaVantage

import "stockfyApi/client"

func CompanyOverviewAlpha(symbol string) CompanyOverview {
	url := "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" +
		symbol + "&apikey=KIUG1ZKFZ13BI08F"

	var companyOverview CompanyOverview

	client.RequestAndAssignToBody(url, &companyOverview)

	return companyOverview
}
