package dbverification

type Repository interface {
	VerifyRowExistence(table string, condition string) bool
}

type Application struct {
	repo Repository
}
