package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/myerrors"
)

type ProductHandler struct {
	DB *db.ProductDB
}

func NewProductHandler(productDB *db.ProductDB) *ProductHandler {
	return &ProductHandler{DB: productDB}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	ImagePath   string  `json:"image_path"`
	Stock       int     `json:"stock" binding:"min=0"`
}

func (h *ProductHandler) List(c *gin.Context) {
	limit := 20
	offset := 0

	if l, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	if o, err := strconv.Atoi(c.DefaultQuery("offset", "0")); err == nil && o >= 0 {
		offset = o
	}

	products, err := h.DB.ListProducts(limit, offset)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	if products == nil {
		products = []db.Product{}
	}

	JSONSuccess(c, http.StatusOK, gin.H{"products": products, "limit": limit, "offset": offset})
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.DB.GetProductByID(id)
	if err != nil {
		JSONError(c, http.StatusNotFound, myerrors.ErrProductNotFound.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, gin.H{"product": product})
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.DB.CreateProduct(req.Name, req.Description, req.Price, req.ImagePath, req.Stock)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusCreated, gin.H{"product": product})
}
