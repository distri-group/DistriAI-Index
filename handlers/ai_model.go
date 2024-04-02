package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

type ModelCreateReq struct {
	Name      string `binding:"required,max=50"`
	Framework uint8  `binding:"required"`
	License   uint8  `binding:"required"`
	Type1     uint32 `binding:"required"`
	Type2     uint32 `binding:"required"`
	Tags      string `binding:"max=128"`
}

func ModelCreate(context *gin.Context) {
	account := getAuthAccount(context)
	var req ModelCreateReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var count int64
	tx := common.Db.Model(&model.AiModel{}).
		Where("owner = ?", account).
		Where("name = ?", req.Name).
		Count(&count)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	if count > 0 {
		resp.Fail(context, "Duplicate model name")
		return
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	likes := uint32(rnd.Int31n(100))

	aiModel := model.AiModel{
		Owner:     account,
		Name:      req.Name,
		Framework: req.Framework,
		License:   req.License,
		Type1:     req.Type1,
		Type2:     req.Type2,
		Tags:      req.Tags,
		Downloads: likes + uint32(rnd.Int31n(1000)),
		Likes:     likes,
	}
	if err := common.Db.Create(&aiModel).Error; err != nil {
		logs.Error(fmt.Sprintf("Database error: %v \n", err))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, "")
}

type ModelListReq struct {
	Name    string
	Type1   uint32
	Type2   uint32
	OrderBy string
}

type ModelListResponse struct {
	List []model.AiModel
	PageResp
}

func ModelList(context *gin.Context) {
	var req ModelListReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response ModelListResponse
	aiModel := model.AiModel{Type1: req.Type1, Type2: req.Type2}
	tx := common.Db.Model(&aiModel).Where(&aiModel)
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
	if tx.Scopes(Paginate(context)).Find(&response.List); tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type ModelGetReq struct {
	Owner string `binding:"required"`
	Name  string `binding:"required"`
}

func ModelGet(context *gin.Context) {
	var req ModelGetReq
	if err := context.ShouldBindUri(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var aiModel model.AiModel
	tx := common.Db.
		Where("owner = ?", req.Owner).
		Where("name = ?", req.Name).
		Take(&aiModel)
	if tx.Error != nil {
		resp.Fail(context, "Model not found")
		return
	}

	resp.Success(context, aiModel)
}

type ModelPresignReq struct {
	Id       uint   `binding:"required"`
	FilePath string `binding:"required"`
	Method   string `binding:"required,oneof= PUT DELETE"`
}

func ModelPresign(context *gin.Context) {
	account := getAuthAccount(context)
	var req ModelPresignReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var aiModel model.AiModel
	if err := common.Db.Where("id = ? AND owner = ?", req.Id, account).Take(&aiModel).Error; err != nil {
		resp.Fail(context, "Model not found")
		return
	}

	objectKey := fmt.Sprintf("model/%s/%s/%s", aiModel.Owner, aiModel.Name, req.FilePath)
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
