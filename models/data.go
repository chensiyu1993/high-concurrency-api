package models

import (
	"time"
)

type Data struct {
	ID        string    `json:"id" gorm:"primary_key"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDataRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateDataRequest struct {
	Content string `json:"content" binding:"required"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
} 