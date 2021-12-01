package utils

import "errors"

var ErrCovertContent = errors.New("fail to covert content")

type Page struct {
	Size      int         `json:"size"`
	Page      int         `json:"page"`
	Total     int64       `json:"total"`
	TotalPage int64       `json:"totalPage"`
	Content   interface{} `json:"content"`
}

type Pageable struct {
	Size int            `json:"size"`
	Page int            `json:"page"`
	Sort map[string]int `json:"sorts"`
}
