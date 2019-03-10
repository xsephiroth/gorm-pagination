package pagination

import (
	"math"

	"github.com/jinzhu/gorm"
)

// Response is a json response struct
type Response struct {
	Total     int         `json:"total"`
	TotalPage int         `json:"total_page"`
	Page      int         `json:"page"`
	Next      bool        `json:"next"`
	Prev      bool        `json:"prev"`
	Results   interface{} `json:"results"`
}

// PagePagination accept page and page_size params,
// use gorm limit and offset implement pagination
// page: request page
// size: page items size
// out: gorm.DB.Find(&out)
func PagePagination(db *gorm.DB, page int, size int, out interface{}) *Response {
	// setup default size
	if size == 0 {
		size = 10
	}

	// count total, with user setup where, before limit and offset
	var total int
	db.Count(&total)
	tpage := totalPage(total, size)

	// valid request page
	if page == 0 {
		page = 1
	}
	// page request <= total_page
	if page > tpage {
		page = tpage
	}

	offset := (page - 1) * size
	query(db, offset, size, out)

	return &Response{
		Total:     total,
		TotalPage: tpage,
		Page:      page,
		Next:      hasNext(total, offset, size),
		Prev:      hasPrev(tpage, page),
		Results:   out,
	}
}

func query(db *gorm.DB, offset int, limit int, out interface{}) *gorm.DB {
	return db.Offset(offset).Limit(limit).Find(out)
}

func hasNext(total int, offset int, limit int) bool {
	return offset+limit < total
}

func hasPrev(totalPage, page int) bool {
	if page <= 1 {
		return false
	}
	return totalPage >= page
}

func totalPage(total int, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
