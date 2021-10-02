package presenter

type AssetBody struct {
	AssetType string `json:"assetType"`
	// Sector    string `json:"sector"`
	Symbol   string `json:"symbol"`
	Fullname string `json:"fullname"`
	Country  string `json:"country"`
}

type AssetApiReturn struct {
	Id         string `json:"id"`
	Preference string `json:"preference"`
	Fullname   string `json:"fullname"`
	Symbol     string `json:"symbol"`
}

func ConvertAssetToApiReturn(id string, preference string, fullname string,
	symbol string) AssetApiReturn {
	return AssetApiReturn{
		Id:         id,
		Preference: preference,
		Fullname:   fullname,
		Symbol:     symbol,
	}
}
