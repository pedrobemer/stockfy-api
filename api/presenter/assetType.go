package presenter

type AssetType struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
}

func ConvertAssetTypeToApiReturn(id string, assetType string, name string,
	country string) *AssetType {

	if id == "" && assetType == "" && name == "" && country == "" {
		return nil
	}

	return &AssetType{
		Id:      id,
		Type:    assetType,
		Name:    name,
		Country: country,
	}
}
