package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/resp"
	"fmt"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"strings"
)

type DatasetPresignReq struct {
	Id       uint   `binding:"required"`
	FilePath string `binding:"required"`
	Method   string `binding:"required,oneof= PUT DELETE"`
}

func DatasetPresign(context *gin.Context) {
	account := getAuthAccount(context)
	var req DatasetPresignReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var dataset model.Dataset
	if err := common.Db.Where("id = ? AND owner = ?", req.Id, account).Model(&dataset).Error; err != nil {
		resp.Fail(context, "Dataset not found")
	}
	objectKey := fmt.Sprintf("dataset/%s/%s/%s", dataset.Owner, dataset.Name, req.FilePath)
	presignedPutRequest := new(v4.PresignedHTTPRequest)
	var err error
	switch req.Method {
	case "PUT":
		presignedPutRequest, err = common.S3Presigner.PutObject("distriai", objectKey, 3600)
	case "DELETE":
		presignedPutRequest, err = common.S3Presigner.DeleteObject("distriai", objectKey)
	}
	if err != nil {
		resp.Fail(context, "Create presigned URL error")
		return
	}
	resp.Success(context, presignedPutRequest.URL)
}

type DatasetListReq struct {
	Name    string
	Type1   uint8
	Type2   uint8
	OrderBy string
}

type DatasetListResponse struct {
	List []model.Dataset
	PageResp
}

func DatasetList(context *gin.Context) {
	var req DatasetListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response DatasetListResponse
	//dataset := model.Dataset{Type1: req.Type1, Type2: req.Type2}
	//tx := common.Db.Model(&dataset).Where(&dataset)
	tx := common.Db.Table("datasets").
		Select("datasets.owner, datasets.name, datasets.scale, datasets.license, datasets.type1, datasets.type2, datasets.tags, datasets.update_time, dataset_heats.downloads, dataset_heats.likes").
		Joins("LEFT JOIN dataset_heats ON datasets.owner = dataset_heats.owner AND datasets.name = dataset_heats.name")
	if req.Type1 != 0 && req.Type2 != 0 {
		tx.Where("datasets.type1 = ? AND datasets.type2 = ?", req.Type1, req.Type2)
	}
	if "" != req.Name {
		name := strings.ReplaceAll(req.Name, "%", "\\%")
		tx.Where("name LIKE ?", "%"+name+"%")
	}
	if err := tx.Count(&response.Total).Error; err != nil {
		resp.Fail(context, "Database error")
		return
	}

	switch req.OrderBy {
	case "downloads DESC", "likes DESC":
		tx.Order(req.OrderBy)
	default:
		tx.Order("updated_at DESC")
	}
	if tx.Scopes(Paginate(context)).Find(&response.List).Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type DatasetGetReq struct {
	Owner string `binding:"required"`
	Name  string `binding:"required"`
}

func DatasetGet(context *gin.Context) {
	var req DatasetGetReq
	if err := context.ShouldBindUri(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	// query through owner and name
	var dataset model.Dataset
	err := common.Db.Transaction(func(tx *gorm.DB) error {
		tx.Where("owner = ?", req.Owner).
			Where("name = ?", req.Name).
			Take(&dataset)
		if tx.Error != nil {
			return tx.Error
		}

		tx.Model(&model.DatasetHeat{}).
			Where("owner = ?", req.Owner).
			Where("name = ?", req.Name).
			Update("clicks", gorm.Expr("clicks + ?", 1))
		return tx.Error
	})
	if err != nil {
		resp.Fail(context, "Database error")
		return
	}
	/*
		tx := common.Db.
			Where("owner = ?", req.Owner).
			Where("name = ?", req.Name).
			Take(&dataset)
		if tx.Error != nil {
			resp.Fail(context, "Dataset not found")
			return
		}*/

	resp.Success(context, dataset)
}

func DatasetDownload(context *gin.Context) {
	//var dataset model.DatasetHeat
	var req DatasetGetReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}
	tx := common.Db.Model(&model.DatasetHeat{}).
		Where("name = ?", req.Name).
		Where("owner = ?", req.Owner).
		Update("downloads", gorm.Expr("downloads + ?", 1))
	if tx.Error != nil {
		resp.Fail(context, "dataset not found")
		return
	}

	resp.Success(context, "")
}

type DatasetLikeReq struct {
	Owner string `binding:"required"`
	Name  string `binding:"required"`
	Like  bool
}

func DatasetLike(context *gin.Context) {
	account := getAuthAccount(context)
	var req DatasetLikeReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	num := -1
	if req.Like {
		num = 1
	}

	err := common.Db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.DatasetHeat{}).
			Where("owner = ?", account).
			Where("name = ?", req.Name).
			Update("likes", gorm.Expr("likes + ?", num))
		if tx.Error != nil {
			return tx.Error
		}
		like := model.DatasetLike{
			Account: account,
			Owner:   req.Owner,
			Name:    req.Name,
		}

		if req.Like {
			// tx.Where(like).FirstOrCreate(&model.DatasetLike{}, like)
			return tx.Create(&like).Error
		}
		return tx.Where(&like).
			Delete(&model.DatasetLike{}).
			Error
	})

	if err != nil {
		resp.Fail(context, err.Error())
		return
	}
	resp.Success(context, "")
}

func DatasetIsLike(context *gin.Context) {
	account := getAuthAccount(context)
	var req DatasetGetReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var count int64
	tx := common.Db.Model(&model.DatasetLike{}).
		Where("account = ? and owner = ? and name = ?", account, req.Owner, req.Name).
		Count(&count)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, count > 0)
}

func DatasetLikeCount(context *gin.Context) {
	var req DatasetGetReq

	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var datasetHeat model.DatasetHeat
	tx := common.Db.Model(&model.DatasetHeat{}).
		Where("owner = ?", req.Owner).
		Where("name = ?", req.Name).
		Take(&datasetHeat)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, datasetHeat.Likes)
}
