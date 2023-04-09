package entity

import "time"

type Usecases struct {
	UserUsecase     UserUsecase
	OrderUsecase    OrderUsecase
	ProductUsecase  ProductUsecase
	CategoryUsecase CategoryUsecase
	CartUsecase     CartUsecase
}

type PaginationParam struct {
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	Category int `json:"category"`
}

type State string

const (
	Enabled  State = "enabled"
	Disabled State = "disabled"
	Deleted  State = "deleted"
)

func NowUTC() time.Time {
	return time.Now().Round(time.Microsecond).UTC()
}
