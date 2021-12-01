package services

import (
	"fmt"
)

type ServiceError struct {
	ServiceName string
	Code        int
	Msg         string
	Err         error
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("code(%d): Msg(%s)", e.Code, e.ServiceName+": "+e.Msg)
}

func (e ServiceError) GetCode() int {
	return e.Code
}

func (e ServiceError) GetMsg() string {
	return e.Msg
}

type ServiceEvent int

type Event struct {
	Code int
	Msg  string
}

const (
	UserNotFound           ServiceEvent = 201
	UserRegistered         ServiceEvent = 202
	UserNameBeenRegistered ServiceEvent = 203
	PasswordWrong          ServiceEvent = 204
	CreationNotFound       ServiceEvent = 301
	BrandNotFound          ServiceEvent = 401
	BrandHaveCreation      ServiceEvent = 402
	StockExisted           ServiceEvent = 501
	StockNotFound          ServiceEvent = 502
	ContractDuplicate      ServiceEvent = 601
)

func (e ServiceEvent) GetEvent() *Event {
	switch e {
	case UserNotFound:
		return &Event{int(e), "user not found"}
	case UserRegistered:
		return &Event{int(e), "user has registered"}
	case UserNameBeenRegistered:
		return &Event{int(e), "user name have been registered"}
	case PasswordWrong:
		return &Event{int(e), "password is wrong"}
	case CreationNotFound:
		return &Event{int(e), "creation not found"}
	case BrandNotFound:
		return &Event{int(e), "brand not found"}
	case BrandHaveCreation:
		return &Event{int(e), "brand have creation"}
	case StockExisted:
		return &Event{int(e), "stock has existed"}
	case StockNotFound:
		return &Event{int(e), "stock not found"}
	case ContractDuplicate:
		return &Event{int(e), "contract is duplicate"}
	default:
		return &Event{int(e), "unknown"}
	}
}
