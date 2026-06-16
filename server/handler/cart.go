package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/myerrors"
)

type CartHandler struct {
	DB *db.CartDB
}

func NewCartHandler(cartDB *db.CartDB) *CartHandler {
	return &CartHandler{DB: cartDB}
}

type AddToCartRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

func (h *CartHandler) View(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
		return
	}

	items, err := h.DB.GetCart(userID)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	if items == nil {
		items = []db.CartItem{}
	}

	JSONSuccess(c, http.StatusOK, gin.H{"cart": items})
}

func (h *CartHandler) Add(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.DB.AddItem(userID, req.ProductID, req.Quantity)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, gin.H{"item": item})
}

func (h *CartHandler) Update(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
		return
	}

	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid product id")
		return
	}

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Quantity == 0 {
		if err := h.DB.RemoveItem(userID, productID); err != nil {
			JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
			return
		}
		JSONSuccess(c, http.StatusOK, gin.H{"message": "item removed"})
		return
	}

	item, err := h.DB.UpdateItem(userID, productID, req.Quantity)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, gin.H{"item": item})
}

func (h *CartHandler) Remove(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
		return
	}

	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.DB.RemoveItem(userID, productID); err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, gin.H{"message": "item removed"})
}
