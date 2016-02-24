package base

import (
	cmn "circle/common"
)

type MyPage struct {
	Page cmn.Page `json:"page"`
	Id   uint32   `json:"id"`
}
