package http

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"go-store/internal/entity"
	dto "go-store/internal/product/dto"
	errorstatus "go-store/utils/errors"
	httphelper "go-store/utils/http"
)

// ProductHandler  represent the httphandler for article
type ProductHandler struct {
	PrUsecase entity.ProductUsecase
	user      entity.UserUsecase
	srvLog    *logrus.Entry
}

const (
	sucsess = "sucsess"
)

// NewProductHandler will initialize the articles/ resources endpoint
func NewProductHandler(handler *gin.RouterGroup, mdw gin.HandlerFunc, uc *entity.Usecases, srvLog *logrus.Entry) {
	ph := &ProductHandler{
		PrUsecase: uc.ProductUsecase,
		user:      uc.UserUsecase,
		srvLog:    srvLog,
	}
	h := handler.Group("/product")
	{
		h.GET("", ph.products)
		h.GET("/:skuCode", ph.singleProduct)
		h.GET("/skuValue/:skuValueId", ph.optionBySkuValue)
	}
	a := handler.Group("/admin")
	{
		a.GET("/product", ph.productBySkus)
		a.POST("/product", ph.createProduct)
		a.PUT("/product/:productId", ph.updateProduct)
		a.DELETE("/product/:productId", ph.deleteProduct)
		a.POST("/product/sku/:skuId/option", ph.createProductOption)
		a.DELETE("/product/skuValue/:skuValueId", ph.deleteProductOption)
		a.POST("/product/:productId/sku", ph.createSku)
		a.PUT("/product/sku/:skuId", ph.updateSku)
		a.DELETE("/product/sku/:skuId", ph.deleteSku)
	}

}

func (ph *ProductHandler) products(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "products"})
	var listReq dto.ProductListRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&listReq); err != nil {
		srvLog.WithError(err).Error("format is wrong")
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	srvLog.WithFields(log.Fields{"limit": listReq.Limit, "offset": listReq.Offset, "filter": listReq.Filter})

	result, err := ph.PrUsecase.GetSku(c, listReq.Limit, listReq.Offset, listReq.Filter)
	if err != nil {
		srvLog.Warning("Cannot GetSkuProduct, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ph *ProductHandler) productBySkus(c *gin.Context) {
	var err error

	srvLog := log.WithFields(log.Fields{"func": "productBySkus"})

	pageParam, err := httphelper.SortingParams(c)
	if err != nil {
		srvLog.Warning(err)
		httphelper.SendResponse(c, nil, errorstatus.ErrBadReq)
		return
	}

	srvLog.WithFields(log.Fields{"limit": pageParam.Limit, "offset": pageParam.Offset, "categoryId": pageParam.Category})

	result, err := ph.PrUsecase.GetProductSkus(c, pageParam.Limit, pageParam.Offset, pageParam.Category)

	if err != nil {
		srvLog.Warning("Cannot get GetProduct, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ph *ProductHandler) singleProduct(c *gin.Context) {
	skuCode := c.Param("skuCode")

	srvLog := log.WithFields(log.Fields{"skuCode": skuCode, "func": "GetSingleProductHandler"})

	result, err := ph.PrUsecase.GetSingleProduct(c, skuCode)
	if err != nil {
		srvLog.Warning(err)
	}

	if err != nil {
		srvLog.Warning("Cannot get product, Err: ", err)
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, result, nil)
}

func (ph *ProductHandler) createProduct(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProduct"})

	prodForm, err := httphelper.ProdCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}
	productId, err := ph.PrUsecase.CreateProduct(c, prodForm)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, productId, nil)
}

func (ph *ProductHandler) updateProduct(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "ProductRemove"})

	productIdStr := c.Param("productId")

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.DeleteProduct.productId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	prodForm, err := httphelper.ProdCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	prodForm.Id = productId
	err = ph.PrUsecase.UpdateProduct(c, prodForm)
	if err != nil {
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) deleteProduct(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "ProductRemove"})

	productIdStr := c.Param("productId")

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		srvLog.WithError(err).Warning("utils.DeleteProduct.productId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	err = ph.PrUsecase.DeleteProduct(c, productId)
	if err != nil {
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) createProductOption(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProductOption"})

	skuIdStr := c.Param("skuId")

	skuId, err := strconv.Atoi(skuIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.OptionForm.optionId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	optionId, optionValueId, err := httphelper.OptionForm(c)

	err = ph.PrUsecase.CreateProductOption(c, skuId, *optionId, *optionValueId)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) deleteProductOption(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "ProductRemove"})

	skuValueIdStr := c.Param("skuValueId")

	skuValueId, err := strconv.Atoi(skuValueIdStr)
	if err != nil {
		srvLog.WithError(err).Warning("utils.DeleteProductOption.skuValueIdStr")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	err = ph.PrUsecase.DeleteProductOption(c, skuValueId)
	if err != nil {
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) createSku(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProduct"})

	productIdStr := c.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.CreateSku.productId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	skuForm, err := httphelper.SkuCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	skuForm.ProductId = productId
	err = ph.PrUsecase.CreateSku(c, skuForm)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) updateSku(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "CreateProduct"})

	skuIdStr := c.Param("skuId")
	skuId, err := strconv.Atoi(skuIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.UpdateSku.skuId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	skuForm, err := httphelper.SkuCreateForm(c)
	if err != nil {
		srvLog.WithError(err).Warning("httphelper.ProdCreateForm")
		httphelper.SendResponse(c, nil, err)
		return
	}

	skuForm.Id = skuId

	err = ph.PrUsecase.UpdateSku(c, skuForm)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.CreateProduct")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) deleteSku(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "DeleteSku"})

	skuIdStr := c.Param("skuId")
	skuId, err := strconv.Atoi(skuIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.DeleteSku.skuId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	err = ph.PrUsecase.DeleteSku(c, skuId)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.DeleteSku")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, sucsess, nil)
}

func (ph *ProductHandler) optionBySkuValue(c *gin.Context) {
	srvLog := log.WithFields(log.Fields{"func": "OptionBySkuValue"})

	skuValueIdStr := c.Param("skuValueId")
	skuValueId, err := strconv.Atoi(skuValueIdStr)
	if err != nil {
		logrus.WithError(err).Warning("utils.OptionBySkuValue.skuValueId")
		httphelper.SendResponse(c, errorstatus.ErrBadReq, err)
		return
	}

	option, err := ph.PrUsecase.GetSkuOption(c, skuValueId)
	if err != nil {
		srvLog.WithError(err).Warning("ph.PrUsecase.DeleteSku")
		httphelper.SendResponse(c, nil, err)
		return
	}

	httphelper.SendResponse(c, option, nil)
}
