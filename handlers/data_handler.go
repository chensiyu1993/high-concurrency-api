package handlers

import (
	"net/http"
	"time"

	"high-concurrency-api/dao"
	"high-concurrency-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DataHandler struct {
	dao *dao.DataDAO
}

func NewDataHandler(dao *dao.DataDAO) *DataHandler {
	return &DataHandler{
		dao: dao,
	}
}

func (h *DataHandler) Create(c *gin.Context) {
	var req models.CreateDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request",
		})
		return
	}

	data := &models.Data{
		ID:        uuid.New().String(),
		Content:   req.Content,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.dao.Create(c.Request.Context(), data); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to create data",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data:    data,
	})
}

func (h *DataHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request",
		})
		return
	}

	data := &models.Data{
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}

	if err := h.dao.Update(c.Request.Context(), id, data); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to update data",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
	})
}

func (h *DataHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.dao.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to delete data",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
	})
}

func (h *DataHandler) Get(c *gin.Context) {
	id := c.Param("id")

	data, err := h.dao.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Data not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data:    data,
	})
} 