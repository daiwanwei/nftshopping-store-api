package items

type Response interface {
	GetCode() int
	GetMsg() string
}

type NonDataResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r NonDataResponse) GetCode() int {
	return r.Code
}

func (r NonDataResponse) GetMsg() string {
	return r.Msg
}

type DataResponse struct {
	NonDataResponse
	Data interface{} `json:"data"`
}

type PageRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
