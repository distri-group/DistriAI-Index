package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
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
	Owner   string
	Name    string
	Type1   uint8
	Type2   uint8
	OrderBy string
}

type DatasetDetail struct {
	model.Dataset
	Likes     uint32
	Downloads uint32
	Clicks    uint32
}

type DatasetListResponse struct {
	List []DatasetDetail
	PageResp
}

func DatasetList(context *gin.Context) {
	var req DatasetListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response DatasetListResponse
	tx := common.Db.Table("datasets").
		Select("datasets.owner, datasets.name, datasets.scale, datasets.license, datasets.type1, datasets.type2, datasets.tags, datasets.create_time, datasets.update_time, dataset_heats.downloads, dataset_heats.likes, dataset_heats.clicks").
		Joins("LEFT JOIN dataset_heats ON datasets.owner = dataset_heats.owner AND datasets.name = dataset_heats.name")
	if req.Type1 != 0 && req.Type2 != 0 {
		tx.Where("datasets.type1 = ? AND datasets.type2 = ?", req.Type1, req.Type2)
	}
	if "" == req.Owner {
		tx.Where("dataset_heats.review = 1")
	} else {
		tx.Where("datasets.owner = ?", req.Owner)
	}
	if "" != req.Name {
		name := strings.ReplaceAll(req.Name, "%", "\\%")
		tx.Where("datasets.name LIKE ?", "%"+name+"%")
	}
	if err := tx.Count(&response.Total).Error; err != nil {
		resp.Fail(context, "Database error")
		return
	}

	switch req.OrderBy {
	case "update_time DESC", "downloads DESC", "likes DESC":
		tx.Order(req.OrderBy)
	default:
		tx.Order("`downloads` + `likes` + `clicks` DESC")
	}
	if tx.Scopes(Paginate(context)).Find(&response.List).Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

func DatasetLikes(context *gin.Context) {
	account := getAuthAccount(context)
	var req DatasetListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response DatasetListResponse
	tx := common.Db.Table("datasets").
		Select("datasets.owner, datasets.name, datasets.scale, datasets.license, datasets.type1, datasets.type2, datasets.tags, datasets.create_time, datasets.update_time, dataset_heats.downloads, dataset_heats.likes, dataset_heats.clicks").
		Joins("INNER JOIN dataset_likes ON datasets.owner = dataset_likes.owner AND datasets.name = dataset_likes.name AND dataset_likes.account = ?", account).
		Joins("LEFT JOIN dataset_heats ON datasets.owner = dataset_heats.owner AND datasets.name = dataset_heats.name")

	if req.Type1 != 0 && req.Type2 != 0 {
		tx.Where("datasets.type1 = ? AND datasets.type2 = ?", req.Type1, req.Type2)
	}
	if "" != req.Name {
		name := strings.ReplaceAll(req.Name, "%", "\\%")
		tx.Where("datasets.name LIKE ?", "%"+name+"%")
	}
	if err := tx.Count(&response.Total).Error; err != nil {
		resp.Fail(context, "Database error")
		return
	}

	switch req.OrderBy {
	case "update_time DESC", "downloads DESC", "likes DESC":
		tx.Order(req.OrderBy)
	default:
		tx.Order("`downloads` + `likes` + `clicks` DESC")
	}
	if tx.Scopes(Paginate(context)).Find(&response.List); tx.Error != nil {
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

	var datasetDetail DatasetDetail
	tx := common.Db.Table("datasets").
		Select("datasets.owner, datasets.name, datasets.scale, datasets.license, datasets.type1, datasets.type2, datasets.tags, datasets.create_time, datasets.update_time, dataset_heats.downloads, dataset_heats.likes, dataset_heats.clicks").
		Joins("LEFT JOIN dataset_heats ON datasets.owner = dataset_heats.owner AND datasets.name = dataset_heats.name").
		Where("datasets.owner = ? AND datasets.name = ?", req.Owner, req.Name).
		Take(&datasetDetail)
	if tx.Error != nil {
		resp.Fail(context, "Dataset not found")
		return
	}

	tx = common.Db.Model(&model.DatasetHeat{}).
		Where("owner = ? AND name = ?", req.Owner, req.Name).
		Update("clicks", gorm.Expr("clicks + ?", 1))
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
	}

	resp.Success(context, datasetDetail)
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

	var count int64
	tx := common.Db.Model(&model.DatasetLike{}).
		Where("account = ? AND owner = ? AND name = ?", account, req.Owner, req.Name).
		Count(&count)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	num := -1
	if req.Like {
		if count > 0 {
			resp.Fail(context, "Already liked")
			return
		}
		num = 1
	} else {
		if count == 0 {
			resp.Fail(context, "Yet not liked")
			return
		}
	}

	err := common.Db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.DatasetHeat{}).
			Where("owner = ? AND name = ?", req.Owner, req.Name).
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
	if err := context.ShouldBindQuery(&req); err != nil {
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
