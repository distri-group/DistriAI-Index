package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type HttpHeader struct {
	Account string
}

type PageReq struct {
	Page     int
	PageSize int
}

type PageResp struct {
	Total int64
}

func Paginate(context *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var req PageReq
		_ = context.ShouldBindBodyWith(&req, binding.JSON)
		page := req.Page
		pageSize := req.PageSize
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
