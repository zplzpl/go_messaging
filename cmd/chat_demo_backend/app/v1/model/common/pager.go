package common

type Pager struct {
	PageSize   int64 `json:"page_size"`
	PageNo     int64 `json:"page_no"`
	TotalCount int64 `json:"total_count"`
	TotalPage  int64 `json:"total_page"`
}
