package utils

type Pagination struct {
	Page        int32 `json:"page"`
	Limit       int32 `json:"limit"`
	TotalRecord int32 `json:"total_record"`
	TotalPage   int32 `json:"total_page"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

func NewPagination(page, limit, total_record int32) *Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limitInt := GetIntEnv("LIMIT_ITEM_ON_PER_PAGE", 10)
		limit = int32(limitInt)
	}
	totalPages := (total_record + limit - 1) / limit
	return &Pagination{
		Page:        page,
		Limit:       limit,
		TotalRecord: total_record,
		TotalPage:   totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
}

func NewPaginationResponse(data any, page, limit, total_record int32) map[string]any {
	return map[string]any{
		"data":       data,
		"Pagination": NewPagination(page, limit, total_record),
	}
}
