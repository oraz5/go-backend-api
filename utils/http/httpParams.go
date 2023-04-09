package http

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"go-store/internal/entity"
	errorStatus "go-store/utils/errors"
)

func PaginationParams(c *gin.Context) (cr *entity.PaginationParam, err error) {
	limit_s, ok := c.GetQuery("limit")
	if !ok {
		err = errors.New("params not include 'limit'")
		return nil, err
	}
	offset_s, ok := c.GetQuery("offset")
	if !ok {
		err = errors.New("params not include 'offset'")
		return nil, err
	}

	limit, err := strconv.Atoi(limit_s)
	if err != nil {
		err = fmt.Errorf("Limit is not integer, Err: %w", err)
		return
	}

	offset, err := strconv.Atoi(offset_s)
	if err != nil {
		err = fmt.Errorf("Offset is not integer, Err: %w", err)
		return
	}

	cr = &entity.PaginationParam{
		Limit:  limit,
		Offset: offset,
	}
	return cr, nil
}

func SortingParams(c *gin.Context) (cr *entity.PaginationParam, err error) {
	limit_s, ok := c.GetQuery("limit")
	if !ok {
		err = errors.New("params not include 'limit'")
		return nil, err
	}
	offset_s, ok := c.GetQuery("offset")
	if !ok {
		err = errors.New("params not include 'offset'")
		return nil, err
	}
	categoryId_s, ok := c.GetQuery("categoryId")
	if !ok {
		err = errors.New("params not include 'categoryId'")
		return nil, err
	}

	limit, err := strconv.Atoi(limit_s)
	if err != nil {
		err = fmt.Errorf("Limit is not integer, Err: %w", err)
		return
	}

	offset, err := strconv.Atoi(offset_s)
	if err != nil {
		err = fmt.Errorf("Offset is not integer, Err: %w", err)
		return
	}

	categoryId, err := strconv.Atoi(categoryId_s)
	if err != nil {
		err = fmt.Errorf("Category ID wrong! Err: %w", err)
		return
	}

	cr = &entity.PaginationParam{
		Limit:    limit,
		Offset:   offset,
		Category: categoryId,
	}
	return cr, nil
}

func ProdCreateForm(c *gin.Context) (*entity.Product, error) {
	productName := c.PostForm("productName")
	productDescr := c.PostForm("productDescr")
	categoryID := c.PostForm("categoryId")
	brandID := c.PostForm("brandId")
	regionID := c.PostForm("regionId")

	categoryIdInt, err := strconv.Atoi(categoryID)
	if err != nil {
		logrus.WithError(err).Warning("utils.ProdCreateForm.categoryID")
		return nil, err
	}
	brandIdInt, err := strconv.Atoi(brandID)
	if err != nil {
		logrus.WithError(err).Warning("utils.ProdCreateForm.brandID")
		return nil, err
	}
	regionIdInt, err := strconv.Atoi(regionID)
	if err != nil {
		logrus.WithError(err).Warning("utils.ProdCreateForm.regionID")
		return nil, err
	}

	prodCreate := &entity.Product{
		ProductName: productName,
		Description: productDescr,
		CategoryId:  categoryIdInt,
		BrandId:     brandIdInt,
		RegionId:    regionIdInt,
		CreateTs:    time.Now(),
		UpdateTs:    time.Now(),
		State:       entity.Enabled,
		Version:     0,
	}

	return prodCreate, nil
}

func SkuCreateForm(c *gin.Context) (*entity.Sku, error) {
	productId := c.Param("productId")
	skuCode := c.PostForm("skuCode")
	price := c.PostForm("price")
	quantity := c.PostForm("quantity")
	state := c.PostForm("state")
	priceFl, err := strconv.ParseFloat(price, 32)
	if err != nil {
		logrus.WithError(err).Warning("utils.ProdCreateForm.price")
		return nil, err
	}
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		logrus.WithError(err).Warning("utils.ProdCreateForm.quantity")
		return nil, err
	}
	if skuCode == "" {
		skuCode = c.Param("skuCode")
	}
	if skuCode == "" {
		skuCode = "sku-" + productId + strconv.FormatInt(time.Now().UnixMicro(), 10)
	}
	form, err := c.MultipartForm()
	if err != nil {
		logrus.WithError(err).Warning("utils.ProdCreateForm")
		return nil, err
	}
	images := form.File["images"]
	var imagesName string
	for idx, image := range images {
		image.Filename = "img-" + strconv.Itoa(idx) + ".jpg"
		path := filepath.Join("static", "images", skuCode, image.Filename)
		if idx < len(images)-1 {
			imagesName = imagesName + path + ","
		} else {
			imagesName = imagesName + path
		}
		if _, err := os.Stat(filepath.Dir(path)); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				logrus.WithError(err).Warning("utils.ProdCreateForm")
				return nil, err
			}
		}
		err := c.SaveUploadedFile(image, path)
		if err != nil {
			logrus.WithError(err).Warning("utils.ProdCreateForm")
			return nil, err
		}
	}
	sku := &entity.Sku{
		Sku:        skuCode,
		Price:      float32(priceFl),
		Quantity:   quantityInt,
		LargeImage: imagesName,
		CreateTs:   time.Now(),
		UpdateTs:   time.Now(),
		State:      entity.State(state),
		Version:    0,
	}

	if state == "" {
		sku.State = entity.Enabled
	}

	return sku, nil
}

func OptionForm(c *gin.Context) (optionId *int, optionValueId *int, err error) {
	optionIdStr, ok := c.GetQuery("optionId")
	if !ok {
		logrus.WithError(err).Warning("utils.OptionForm.optionIdStr")
		return nil, nil, err
	}

	optionValueIdStr, ok := c.GetQuery("optionValueId")
	if !ok {
		logrus.WithError(err).Warning("utils.OptionForm.optionValueIdStr")
		return nil, nil, err
	}

	optionID, err := strconv.Atoi(optionIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.OptionForm.optionId")
		return nil, nil, err
	}

	optionValueID, err := strconv.Atoi(optionValueIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.OptionForm.optionValueId")
		return nil, nil, err
	}

	return &optionID, &optionValueID, nil
}

func CategoryForm(c *gin.Context) (category *entity.Category, err error) {
	name := c.PostForm("name")
	parentId := c.PostForm("parentId")

	if name != "" && parentId != "" {
		parentIdInt, err := strconv.Atoi(parentId)
		if err != nil {
			logrus.WithError(err).Warning("utils.CategoryForm.parentIdInt")
			return nil, err
		}

		icon, err := c.FormFile("icon")
		if err != nil {
			logrus.WithError(err).Warning("utils.CategoryForm.MkdirAll")
			return nil, err
		}
		pathIcon := filepath.Join("static", "images", "category", name, "icon", icon.Filename)
		if _, err := os.Stat(filepath.Dir(pathIcon)); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(filepath.Dir(pathIcon), os.ModePerm)
			if err != nil {
				logrus.WithError(err).Warning("utils.CategoryForm.icon.MkdirAll")
				return nil, err
			}
		}
		err = c.SaveUploadedFile(icon, pathIcon)
		if err != nil {
			logrus.WithError(err).Warning("utils.CategoryForm.icon.SaveUploadedFile")
			return nil, err
		}

		image, err := c.FormFile("image")
		if err != nil {
			logrus.WithError(err).Warning("utils.CategoryForm.MkdirAll")
			return nil, err
		}
		pathImage := filepath.Join("static", "images", "category", name, "image", image.Filename)
		if _, err := os.Stat(filepath.Dir(pathImage)); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(filepath.Dir(pathImage), os.ModePerm)
			if err != nil {
				logrus.WithError(err).Warning("utils.CategoryForm.image.MkdirAll")
				return nil, err
			}
		}
		err = c.SaveUploadedFile(image, pathImage)
		if err != nil {
			logrus.WithError(err).Warning("utils.CategoryForm.image.SaveUploadedFile")
			return nil, err
		}

		category = &entity.Category{
			Name:   name,
			Parent: parentIdInt,
			Icon:   pathIcon,
			Image:  pathImage,
		}

	} else {
		logrus.WithError(err).Warning("utils.CategoryForm.FormEmpty")
		return nil, errorStatus.ErrBadReq
	}

	return category, nil
}

func OrderForm(c *gin.Context) (order *entity.Order, err error) {

	phone, ok := c.GetPostForm("phone")
	if !ok {
		logrus.WithError(err).Warning("utils.OrderForm.phone")
		return nil, errorStatus.ErrBadReq
	}

	address, ok := c.GetPostForm("address")
	if !ok {
		logrus.WithError(err).Warning("utils.OrderForm.address")
		return nil, errorStatus.ErrBadReq
	}

	comment, ok := c.GetPostForm("comment")
	if !ok {
		logrus.WithError(err).Warning("utils.OrderForm.comment")
		return nil, errorStatus.ErrBadReq
	}

	notes, ok := c.GetPostForm("notes")
	if !ok {
		logrus.WithError(err).Warning("utils.OrderForm.notes")
		return nil, errorStatus.ErrBadReq
	}

	status, ok := c.GetPostForm("status")

	order = &entity.Order{
		Phone:   phone,
		Address: address,
		Comment: comment,
		Status:  status,
		Notes:   notes,
	}

	return order, nil
}

func UserCreateForm(c *gin.Context) (users *entity.Users, sendMethod *string, err error) {
	username := c.PostForm("username")
	fullName := c.PostForm("fullName")
	password := c.PostForm("password")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	address := c.PostForm("address")
	regionID := c.PostForm("regionId")
	sendMet := c.PostForm("sendMethod")

	regionIdInt, err := strconv.Atoi(regionID)
	if err != nil {
		logrus.WithError(err).Warning("utils.ProdCreateForm.regionID")
		return nil, nil, err
	}

	photo, err := c.FormFile("photo")
	if err != nil {
		logrus.WithError(err).Warning("utils.UserCreateForm.FormFile")
		return nil, nil, err
	}
	pathPhoto := filepath.Join("static", "images", "user", username, "photo", photo.Filename)
	if _, err := os.Stat(filepath.Dir(pathPhoto)); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(pathPhoto), os.ModePerm)
		if err != nil {
			logrus.WithError(err).Warning("utils.CategoryForm.image.MkdirAll")
			return nil, nil, err
		}
	}

	userCreate := &entity.Users{
		Username:    username,
		FullName:    fullName,
		Password:    password,
		Email:       email,
		PhoneNumber: phone,
		Address:     address,
		Photo:       pathPhoto,
		RegionId:    regionIdInt,
		CreateTs:    time.Now(),
		UpdateTs:    time.Now(),
		State:       entity.Enabled,
		Version:     0,
	}

	return userCreate, &sendMet, nil
}
