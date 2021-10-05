package presenter

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
