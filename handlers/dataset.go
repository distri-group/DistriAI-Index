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
	"math/rand"
	"strings"
	"time"
)

type DatasetCreateReq struct {
	Name    string `binding:"required,max=50"`
	Scale   uint8  `binding:"required"`
	License uint8  `binding:"required"`
	Type1   uint32 `binging:"required"`
	Type2   uint32 `binding:"required"`
	Tags    string `binding:"max=128"`
}

// DatasetCreate user create a dataset.
func DatasetCreate(context *gin.Context) {
	account := getAuthAccount(context)
	var req DatasetCreateReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
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
		resp.Fail(context, "Duplicate dataset name")
		return
	}

	//generate random number
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	likes := uint32(rnd.Int31n(100))
	// create new dataset struct
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
	/*if err := common.Db.Create(&dataSet).Error; err != nil {
		logs.Error(fmt.Sprintf("Database error: %v \n", err))
		resp.Fail(context, "Database error")
		return
	}*/

	datasetHeat := model.DatasetHeat{
		Owner:     account,
		Name:      req.Name,
		Likes:     0,
		Downloads: 0,
		Clicks:    0,
	}
	/*if err := common.Db.Create(&datasetHeat).Error; err != nil {
		logs.Error(fmt.Sprintf("Database error: %v \n", err))
		resp.Fail(context, "Database error")
		return
	}*/

	err := common.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dataSet).Error; err != nil {
			return err
		}
		return tx.Create(&datasetHeat).Error
	})

	if err != nil {
		logs.Error(fmt.Sprintf("Database error: %v \n", err))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, "")
}

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
	Type1   uint32
	Type2   uint32
	OrderBy string
}

type DatasetRsp struct {
	Id        uint   `gorm:"primarykey"`
	Owner     string `gorm:"size:44;not null;index:idx_dataset_owner_name"`
	Name      string `gorm:"size:50;not null;index:idx_dataset_owner_name"`
	Scale     uint8  `gorm:"not null"`
	License   uint8  `gorm:"not null"`
	Type1     uint32 `gorm:"not null"`
	Type2     uint32 `gorm:"not null"`
	Tags      string `gorm:"size:128;not null"`
	Downloads uint32 `gorm:"not null"`
	Likes     uint32 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Islike    bool
}

type DatasetList struct {
	List []model.Dataset
	PageResp
}

type DatasetListResponse struct {
	List []DatasetRsp
	PageResp
}

func DatasetListGet(context *gin.Context) {
	account := getAuthAccount(context)
	var req DatasetListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var list DatasetList
	//dataset := model.Dataset{Type1: req.Type1, Type2: req.Type2}
	//tx := common.Db.Model(&dataset).Where(&dataset)
	tx := common.Db.Table("datasets").
		Select("datasets.id, datasets.owner, datasets.name, datasets.Scale, datasets.license, datasets.type1, datasets.type2, datasets.tags, datasets.created_at, dataset_heats.downloads, dataset_heats.likes").
		Joins("INNER JOIN dataset_heats ON datasets.owner = dataset_heats.owner AND datasets.name = dataset_heats.name")
	if req.Type1 != 0 && req.Type2 != 0 {
		tx.Where("datasets.type1 = ? AND datasets.type2 = ?", req.Type1, req.Type2)
	}
	if "" != req.Name {
		name := strings.ReplaceAll(req.Name, "%", "\\%")
		tx.Where("name LIKE ?", "%"+name+"%")
	}
	if err := tx.Count(&list.Total).Error; err != nil {
		resp.Fail(context, "Database error")
		return
	}

	switch req.OrderBy {
	case "downloads DESC", "likes DESC", "updated_at DESC":
		tx.Order(req.OrderBy)
	default:
		//tx.Order("updated_at DESC")
		tx.Order(gorm.Expr("dataset_heats.likes + dataset_heats.downloads + dataset_heats.clicks"))
	}
	if tx.Scopes(Paginate(context)).Find(&list.List).Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	likeList, err := addIsLike(list.List, account)
	if err != nil {
		resp.Fail(context, err.Error())
		return
	}
	res := DatasetListResponse{
		List:     likeList,
		PageResp: list.PageResp,
	}
	resp.Success(context, res)
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
	Owner string `binding:"required" json:"Owner"`
	Name  string `binding:"required" json:"Name"`
	Like  bool   `json:"Like"`
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
			Where("owner = ?", req.Owner).
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

func addIsLike(datasets []model.Dataset, account string) ([]DatasetRsp, error) {
	mp := make(map[string]string)
	if account != "" {
		var datasetLike []model.DatasetLike
		tx := common.Db.Model(&model.DatasetLike{}).
			Where("account = ?", account).
			Find(&datasetLike)
		if tx.Error != nil {
			return nil, tx.Error
		}
		for _, like := range datasetLike {
			mp[like.Owner] = like.Name
		}
	} else {
		mp[""] = ""
	}
	//var mp map[string]string

	var resp []DatasetRsp
	for _, dataset := range datasets {
		resp = append(resp, DatasetRsp{
			Id:        dataset.Id,
			Owner:     dataset.Owner,
			Name:      dataset.Name,
			Scale:     dataset.Scale,
			License:   dataset.License,
			Type1:     dataset.Type1,
			Type2:     dataset.Type2,
			Tags:      dataset.Tags,
			Downloads: dataset.Downloads,
			Likes:     dataset.Likes,
			CreatedAt: dataset.CreatedAt,
			UpdatedAt: dataset.UpdatedAt,
			Islike:    mp[dataset.Owner] == dataset.Name,
		})
	}

	return resp, nil
}
