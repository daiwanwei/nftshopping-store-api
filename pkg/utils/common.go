package utils

type CustomError interface {
	error
	GetCode() int
	GetMsg() string
}

func GetTotalPage(size int64, total int64) (totalPage int64) {
	var result int64
	if size > total {
		return 1
	} else {
		result = total / size
		if total%size > 0 {
			result += 1
		}
	}
	return result
}
