package dbverification

//NewApplication create new use case
func NewApplication(r Repository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) RowValidation(tableName string, condition string) bool {
	return a.repo.VerifyRowExistence(tableName, condition)
}
