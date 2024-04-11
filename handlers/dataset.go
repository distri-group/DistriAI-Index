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
	"time"
)

type DatasetCreateReq struct {
	Name    string `binding:"required,max=50"`
	Scale   uint8  `binding:"required"`
	License uint8  `binding:"required"`
	Type1   uint32 `binging:"required"`
	Type2   uint32 `binding:"required"`
	Tags    string `binding:"required"`
}

func DataSetCreate(context *gin.Context) {
	account := getAuthAccount(context)
	var req DatasetCreateReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, "")
		return
	}

	var count int64
	tx := common.Db.Model(&model.Dataset{}).
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

	//generate random number
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	likes := uint32(rnd.Int31n(100))

	dataSet := model.Dataset{
		Owner:     account,
		Name:      req.Name,
		Scale:     req.Scale,
		License:   req.License,
		Type1:     req.Type1,
		Type2:     req.Type2,
		Tags:      req.Tags,
		Downloads: likes + uint32(rnd.Int31n(1000)),
		Likes:     likes,
	}
	if err := common.Db.Create(&dataSet).Error; err != nil {
		logs.Error(fmt.Sprintf("Database error: %v \n", err))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, "")
}

type DataSetPresignReq struct {
	Id       uint   `binding:"required"`
	FilePath string `binding:"required"`
	Method   string `binding:"required,oneof= PUT DELETE"`
}

func DatasetPresign(context *gin.Context) {
	account := getAuthAccount(context)
	var req DataSetPresignReq
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
