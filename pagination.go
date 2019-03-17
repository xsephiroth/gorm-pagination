package pagination

import (
	"math"

	"github.com/jinzhu/gorm"
)

// Response is a json response struct
type Response struct {
	Total     *int  `json:"total,omitempty"`
	TotalPage *int  `json:"total_page,omitempty"`
	Page      *int  `json:"page,omitempty"`
	Next      *bool `json:"next,omitempty"`
	Prev      *bool `json:"prev,omitempty"`
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
	totalPage := TotalPage(total, size)

	// valid request page
	if page == 0 {
		page = 1
	}
	// page request <= total_page
	if page > totalPage {
		page = totalPage
	}

	offset := (page - 1) * size
	query(db, offset, size, out)

	hn := HasNext(totalPage, page)
	hp := HasPrev(totalPage, page)

	return &Response{
		Total:     &total,
		TotalPage: &totalPage,
		Page:      &page,
		Next:      &hn,
		Prev:      &hp,
	}
}

func query(db *gorm.DB, offset int, limit int, out interface{}) *gorm.DB {
	return db.Offset(offset).Limit(limit).Find(out)
}

func HasNext(totalPage, page int) bool {
	return totalPage > page
}

func HasPrev(totalPage, page int) bool {
	if page <= 1 {
		return false
	}
	return totalPage >= page
}

func TotalPage(total int, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
