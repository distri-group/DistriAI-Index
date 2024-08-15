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

type ModelListReq struct {
	Owner   string
	Name    string
	Type1   uint8
	Type2   uint8
	OrderBy string
}

type ModelDetail struct {
	model.AiModel
	Likes     uint32
	Downloads uint32
	Clicks    uint32
	Status    uint8
	Reason    string
	Size      uint32
}

type ModelListResponse struct {
	List []ModelDetail
	PageResp
}

// ModelList handles the request to list models with various filtering and sorting options.
func ModelList(context *gin.Context) {
	var req ModelListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response ModelListResponse
	tx := common.Db.Table("ai_models").
		Select("ai_models.*, ai_model_heats.likes, ai_model_heats.downloads, ai_model_heats.clicks,  ai_model_heats.status,  ai_model_heats.reason, ai_model_heats.size").
		Joins("LEFT JOIN ai_model_heats ON ai_models.owner = ai_model_heats.owner AND ai_models.name = ai_model_heats.name")

	// Filter based on conditions from the request
	if req.Type1 != 0 && req.Type2 != 0 {
		tx.Where("ai_models.type1 = ? AND ai_models.type2 = ?", req.Type1, req.Type2)
	}
	if "" == req.Owner {
		tx.Where("ai_model_heats.status = 2")
	} else {
		tx.Where("ai_models.owner = ?", req.Owner)
	}
	if "" != req.Name {
		// Handle special characters in fuzzy search
		name := strings.ReplaceAll(req.Name, "%", "\\%")
		tx.Where("ai_models.name LIKE ?", "%"+name+"%")
	}
	if err := tx.Count(&response.Total).Error; err != nil {
		resp.Fail(context, "Database error")
		return
	}
	// Order query results based on orderBy field
	switch req.OrderBy {
	case "update_time DESC", "downloads DESC", "likes DESC":
		tx.Order(req.OrderBy)
	default:
		tx.Order("`downloads` + `likes` + `clicks` DESC")
	}
	// Paginate query and store results in response.List
	if tx.Scopes(Paginate(context)).Find(&response.List); tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

// ModelLikes handles the request to list models that the authenticated account has liked.
func ModelLikes(context *gin.Context) {
	account := getAuthAccount(context)
	var req ModelListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response ModelListResponse
	tx := common.Db.Table("ai_models").
		Select("ai_models.*, ai_model_heats.likes, ai_model_heats.downloads, ai_model_heats.clicks,  ai_model_heats.status,  ai_model_heats.reason, ai_model_heats.size").
		Joins("INNER JOIN ai_model_likes ON ai_models.owner = ai_model_likes.owner AND ai_models.name = ai_model_likes.name AND ai_model_likes.account = ?", account).
		Joins("LEFT JOIN ai_model_heats ON ai_models.owner = ai_model_heats.owner AND ai_models.name = ai_model_heats.name")

	if req.Type1 != 0 && req.Type2 != 0 {
		tx.Where("ai_models.type1 = ? AND ai_models.type2 = ?", req.Type1, req.Type2)
	}
	if "" != req.Name {
		name := strings.ReplaceAll(req.Name, "%", "\\%")
		tx.Where("ai_models.name LIKE ?", "%"+name+"%")
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

type ModelGetReq struct {
	Owner string `binding:"required"`
	Name  string `binding:"required"`
}

// ModelGet handles the GET request to retrieve model details.
func ModelGet(context *gin.Context) {
	var req ModelGetReq
	if err := context.ShouldBindUri(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var modelDetail ModelDetail
	tx := common.Db.Table("ai_models").
		Select("ai_models.*, ai_model_heats.likes, ai_model_heats.downloads, ai_model_heats.clicks,  ai_model_heats.status,  ai_model_heats.reason, ai_model_heats.size").
		Joins("LEFT JOIN ai_model_heats ON ai_models.owner = ai_model_heats.owner AND ai_models.name = ai_model_heats.name").
		Where("ai_models.owner = ? AND ai_models.name = ?", req.Owner, req.Name).
		Take(&modelDetail)
	if tx.Error != nil {
		resp.Fail(context, "Model not found")
		return
	}

	tx = common.Db.Model(&model.AiModelHeat{}).
		Where("owner = ? AND name = ?", req.Owner, req.Name).
		Update("clicks", gorm.Expr("clicks + ?", 1))
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
	}

	resp.Success(context, modelDetail)
}

type ModelPresignReq struct {
	Id       uint   `binding:"required"`
	FilePath string `binding:"required"`
	Method   string `binding:"required,oneof= PUT DELETE"`
}

// ModelPresign handles the generation of a presigned URL for S3 operations on AI models.
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
		// Generate presigned URL for uploading object with a validity of 3600 seconds (1 hour)
		presignedPutRequest, err = common.S3Presigner.PutObject("distriai", objectKey, 3600)
	case "DELETE":
		// Generate presigned URL for deleting object
		presignedPutRequest, err = common.S3Presigner.DeleteObject("distriai", objectKey)
	}

	if err != nil {
		resp.Fail(context, "Create presigned URL error")
		return
	}

	resp.Success(context, presignedPutRequest.URL)
}

func ModelDownload(context *gin.Context) {
	var req ModelGetReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	tx := common.Db.Model(&model.AiModelHeat{}).
		Where("owner = ?", req.Owner).
		Where("name = ?", req.Name).
		Update("downloads", gorm.Expr("downloads + ?", 1))
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	resp.Success(context, "")
}

type ModelLikeReq struct {
	Owner string `binding:"required"`
	Name  string `binding:"required"`
	Like  bool
}

func ModelLike(context *gin.Context) {
	account := getAuthAccount(context)
	var req ModelLikeReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var count int64
	tx := common.Db.Model(&model.AiModelLike{}).
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
		tx.Model(&model.AiModelHeat{}).
			Where("owner = ? AND name = ?", req.Owner, req.Name).
			Update("likes", gorm.Expr("likes + ?", num))
		if tx.Error != nil {
			return tx.Error
		}

		like := model.AiModelLike{
			Account: account,
			Owner:   req.Owner,
			Name:    req.Name,
		}
		if req.Like {
			return tx.Create(&like).Error
		}

		return tx.Where(&like).
			Delete(&model.AiModelLike{}).Error
	})

	if err != nil {
		resp.Fail(context, err.Error())
		return
	}
	resp.Success(context, "")
}

func ModelIsLike(context *gin.Context) {
	account := getAuthAccount(context)
	var req ModelGetReq
	if err := context.ShouldBindQuery(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var count int64
	tx := common.Db.Model(&model.AiModelLike{}).
		Where("account = ? and owner = ? and name = ?", account, req.Owner, req.Name).
		Count(&count)

	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	resp.Success(context, count > 0)
}

func ModelLikeCount(context *gin.Context) {
	var req ModelGetReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}
	var modelHeat model.AiModelHeat
	tx := common.Db.Model(&model.AiModelHeat{}).
		Where("owner = ?", req.Owner).
		Where("name = ?", req.Name).
		Take(&modelHeat)

	if tx.Error != nil {
		resp.Fail(context, tx.Error.Error())
		return
	}

	resp.Success(context, modelHeat.Likes)
}

type ModelUpdateStatusReq struct {
	Name   string `binding:"required"`
	Public bool
}

func ModelStatusUpdate(context *gin.Context) {
	account := getAuthAccount(context)
	var req ModelUpdateStatusReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var aiModelHeat model.AiModelHeat
	tx := common.Db.Model(&model.AiModelHeat{}).
		Where("owner = ? AND name = ?", account, req.Name).
		Take(&aiModelHeat)
	if tx.Error != nil {
		resp.Fail(context, "AiModel not found")
		return
	}
	if req.Public {
		aiModelHeat.Status = 1
	} else {
		aiModelHeat.Status = 0
	}
	aiModelHeat.Reason = ""
	if tx := common.Db.Save(aiModelHeat); tx.Error != nil {
		resp.Fail(context, "Database error")
	}
	resp.Success(context, "")
}

type ModelUpdateSizeReq struct {
	Name string `binding:"required"`
	Size uint32
}

func ModelSizeUpdate(context *gin.Context) {
	account := getAuthAccount(context)
	var req ModelUpdateSizeReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	tx := common.Db.Model(&model.AiModelHeat{}).
		Where("owner = ? AND name = ?", account, req.Name).
		Update("size", req.Size)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	resp.Success(context, "")
}

func ModelSizeTotal(context *gin.Context) {
	account := getAuthAccount(context)

	var size uint32
	tx := common.Db.Model(&model.AiModelHeat{}).
		Select("IFNULL(SUM(size), 0)").
		Where("owner = ?", account).
		Find(&size)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	resp.Success(context, size)
}
