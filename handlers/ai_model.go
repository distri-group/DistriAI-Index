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

type ModelListReq struct {
	Name    string
	Type1   uint8
	Type2   uint8
	OrderBy string
}

type ModelListResponse struct {
	List []model.AiModel
	PageResp
}

func ModelList(context *gin.Context) {
	var req ModelListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response ModelListResponse
	//aiModel := model.AiModel{Type1: req.Type1, Type2: req.Type2}
	//tx := common.Db.Model(&aiModel).Where(&aiModel)
	tx := common.Db.Table("ai_models").
		Select("ai_models.id, ai_models.owner, ai_models.name, ai_models.framework, ai_models.license, ai_models.type1, ai_models.type2, ai_models.tags, ai_models.created_at, ai_model_heats.downloads, ai_model_heats.likes").
		Joins("LEFT JOIN ai_model_heats ON ai_models.owner = ai_model_heats.owner AND ai_models.name = ai_model_heats.name")

	if req.Type1 != 0 && req.Type2 != 0 {
		tx.Where("ai_models.type1 = ? AND ai_models.type2 = ?", req.Type1, req.Type2)
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
	/*tx := common.Db.
		Where("owner = ?", req.Owner).
		Where("name = ?", req.Name).
		Take(&aiModel)
	if tx.Error != nil {
		resp.Fail(context, "Model not found")
		return
	}*/

	err := common.Db.Transaction(func(tx *gorm.DB) error {
		tx.Where("owner = ?", req.Owner).
			Where("name = ?", req.Name).
			Take(&aiModel)
		if tx.Error != nil {
			return tx.Error
		}

		tx.Model(&model.AiModelHeat{}).
			Where("owner = ?", req.Owner).
			Where("name = ?", req.Name).
			Update("clicks", gorm.Expr("clicks + ?", 1))
		return tx.Error
	})
	if err != nil {
		resp.Fail(context, err.Error())
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

	num := -1
	if req.Like {
		num = 1
	}

	err := common.Db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.AiModelHeat{}).
			Where("owner = ?", req.Owner).
			Where("name = ?", req.Name).
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
	if err := context.ShouldBindJSON(&req); err != nil {
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
