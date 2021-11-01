package dbverification

type MockApplication struct {
	repo Repository
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (a *MockApplication) RowValidation(tableName string, condition string) bool {

	if condition == "symbol='SYMBOL_EXIST'" {
		return true
	} else {
		return false
	}
}
