package dbverification

type Mock struct {
}

func NewMockRepo() *Mock {
	return &Mock{}
}

func (m *Mock) VerifyRowExistence(table string, condition string) bool {

	if condition == "True" {
		return true
	} else {
		return false
	}

}
