package common

type RetStatus struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

var SuccessRet = RetStatus{
	Code: 0,
	Msg:  "ok",
}
