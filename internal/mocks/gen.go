package mocks

//go:generate go run go.uber.org/mock/mockgen@latest -destination repo_mock.go -package mocks github.com/bruno-nascimento/translation-api/internal/repository Querier
//go:generate go run go.uber.org/mock/mockgen@latest -destination service_mock.go -package mocks github.com/bruno-nascimento/translation-api/internal/service  Translation
