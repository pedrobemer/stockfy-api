package presenter

type SectorBody struct {
	Sector string `json:"sector"`
}

type Sector struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func ConvertSectorToApiReturn(id string, name string) *Sector {
	if id == "" && name == "" {
		return nil
	}

	return &Sector{
		Id:   id,
		Name: name,
	}
}
