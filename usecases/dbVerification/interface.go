package dbverification

type Repository interface {
	VerifyRowExistence(table string, condition string) bool
}

type Application struct {
	repo Repository
}

type UseCases interface {
	RowValidation(tableName string, condition string) bool
}
