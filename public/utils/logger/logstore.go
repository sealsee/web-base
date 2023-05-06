package logger

type Store interface {
	Save(a any)
	Query(a any)
}

type LogStore interface {
	Store
}

type ErrorLogStore interface {
	Store
}
