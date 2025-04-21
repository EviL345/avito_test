package handler

import (
	openapi "github.com/EviL345/avito_test/internal/gen"
	"github.com/EviL345/avito_test/internal/middleware"
	"github.com/EviL345/avito_test/internal/model/dto/request"
	"github.com/EviL345/avito_test/internal/model/dto/response"
	"github.com/EviL345/avito_test/internal/model/entity"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

const InvalidCity = "invalid city"

type PvzService interface {
	CreatePvz(pvz *entity.Pvz) (*entity.Pvz, error)
	GetPvz(startDate, endDate *time.Time, page, limit *int) ([]response.PvzInfo, error)
}

func (h *Handler) PostPvz(c *gin.Context) {
	log.SetPrefix("handler.PostPvz")

	middleware.Auth(entity.ModeratorRole)(c)
	if c.IsAborted() {
		return
	}

	var req request.Pvz
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "Field validation") {
			c.JSON(400, gin.H{"error": InvalidCity})

			return
		}

		log.Printf("error: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	pvz := &entity.Pvz{
		City: req.City,
	}

	pvz, err := h.pvzService.CreatePvz(pvz)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, pvz.ToResponse())
}

func (h *Handler) GetPvz(c *gin.Context, params openapi.GetPvzParams) {
	log.SetPrefix("handler.GetPvz")

	middleware.Auth(entity.EmployeeRole, entity.ModeratorRole)(c)
	if c.IsAborted() {
		return
	}

	if params.Page == nil {
		params.Page = new(int)
		*params.Page = 1
	}

	if params.Limit == nil {
		params.Limit = new(int)
		*params.Limit = 10
	} else if *params.Limit > 30 {
		*params.Limit = 10
	}

	pvzList, err := h.pvzService.GetPvz(params.StartDate, params.EndDate, params.Page, params.Limit)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, pvzList)
}
