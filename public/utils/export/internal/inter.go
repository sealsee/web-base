package internal

type ImpHandler interface {
	Headers() []string
	Row() map[string]interface{}
}

type ExpHandler interface {
	Before()
	Rows() []map[string]interface{}
	After()
}

type ImpExp interface {
	Import(bytes []byte, handler ImpHandler) error
	ImportWithUrl(url string, handler ImpHandler) error
	ExportSync(headers map[string]string, handler ExpHandler) ([]byte, error)
	ExportAsync(headers map[string]string, handler ExpHandler) (string, error)
}
