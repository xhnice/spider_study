package entity

// Request -
type (
	ParseHandler func(*RequestResult) *ParseRequest

	RequestConfig struct {
		Retry      uint // 重试
		RetryIndex uint
	}

	StudyRequest struct {
		Depth uint
	}

	Request struct {
		URI  string
		Data interface{}
		RequestConfig
		ParseHandler ParseHandler
	}

	RequestResult struct {
		SourceURI string // 数据来源URI
		Body      interface{}
		Data      interface{}
		RequestConfig
	}

	ParseRequest struct {
		Items []interface{}
	}
)
