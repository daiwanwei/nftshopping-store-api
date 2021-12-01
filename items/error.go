package items

import "fmt"

type StatusError struct {
	Code int
	Msg  string
}

func (e StatusError) Error() string {
	return fmt.Sprintf("code(%d): Msg(%s)", e.Code, e.Msg)
}

func (e StatusError) GetCode() int {
	return e.Code
}

func (e StatusError) GetMsg() string {
	return e.Msg
}

func StatusHandler(code int) (err error) {
	statusCode := StatusCode(code)
	err = statusCode.GetError()
	return
}

type StatusCode int

const (
	Success      StatusCode = 200
	Unauthorized StatusCode = 401
	Forbidden    StatusCode = 403
	NotFound     StatusCode = 404
)

func (e StatusCode) GetError() (err error) {
	switch e {
	case Success:
		return
	case Unauthorized:
		err = &ItemError{Code: 1, Msg: "status code:Unauthorized"}
	case Forbidden:
		err = &ItemError{Code: 2, Msg: "status code:Forbidden"}
	case NotFound:
		err = &ItemError{Code: 3, Msg: "status code:NotFound"}
	default:
		err = &ItemError{Code: 0, Msg: "status code:unknown"}
	}
	return
}

func ResponseHandler(response Response) (err error) {
	if response.GetCode() == 200 {
		return
	}
	err = NewMarketError(response)
	return
}

type ResponseCode int

const (
	ContractExisted  ResponseCode = 201
	ContractNotFound ResponseCode = 202
	ItemExisted      ResponseCode = 301
	ItemNotFound     ResponseCode = 302
	OrderNotFound    ResponseCode = 401
)

type ItemError struct {
	Code int
	Msg  string
	Err  error
}

func NewMarketError(res Response) error {
	switch ResponseCode(res.GetCode()) {
	case ContractExisted:
		return &ItemError{Code: res.GetCode(), Msg: "contract have existed"}
	case ContractNotFound:
		return &ItemError{Code: res.GetCode(), Msg: "contract not found"}
	case ItemExisted:
		return &ItemError{Code: res.GetCode(), Msg: "item have existed"}
	case ItemNotFound:
		return &ItemError{Code: res.GetCode(), Msg: "item not found"}
	case OrderNotFound:
		return &ItemError{Code: res.GetCode(), Msg: "order not found"}
	default:
		return &ItemError{Code: res.GetCode(), Msg: res.GetMsg()}
	}
}

func (e ItemError) Error() string {
	return fmt.Sprintf("code(%d): Msg(%s)", e.Code, "market service"+":"+e.Msg)
}

func (e ItemError) GetCode() int {
	return e.Code
}

func (e ItemError) GetMsg() string {
	return e.Msg
}
