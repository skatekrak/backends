package database

import (
	"log"
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	PerPage      int         `json:"perPage" query:"perPage"`
	Page         int         `json:"page" query:"page"`
	TotalResults int64       `json:"totalResults"`
	TotalPages   int         `json:"totalPages"`
	Items        interface{} `json:"items"`
} // @name Pagination

func (p *Pagination) GetPageSize() int {
	if p.PerPage == 0 {
		return 20
	}
	return p.PerPage
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

func (p *Pagination) Scope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := p.GetPageSize()

		tx := db.Session(&gorm.Session{})

		var totalItems int64
		tx.Count(&totalItems)

		p.TotalResults = totalItems
		p.TotalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
		log.Println("pagination", p)

		return db.Offset(p.GetOffset()).Limit(p.GetPageSize())
	}
}
