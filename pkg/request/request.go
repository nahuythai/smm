package request

type Pagination struct {
	Skip  int64 `json:"-"`
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
	Total int64 `json:"total"`
}

func NewPagination(page, limit int64) *Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}
	return &Pagination{
		Skip:  (page - 1) * limit,
		Limit: limit,
		Page:  page,
	}
}

func (p *Pagination) SetTotal(total int64) {
	p.Total = total
}
