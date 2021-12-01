package adapter

type NonDataResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type DataResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
