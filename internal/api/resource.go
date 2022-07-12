package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

//go:generate mockgen -destination ../../mock/resource_api_mock.go -package mock . ResourceApi
type ResourceApi interface {
	Configure()
}

type ResourceApiImpl struct {
	Router          *gin.RouterGroup
	ResourceService service.ResourceService
}

func (impl *ResourceApiImpl) Configure() {
	impl.Router.GET("", impl.FindAll)
	impl.Router.GET("/:resourceID", impl.FindOneByID)
	impl.Router.POST("", impl.Create)
	impl.Router.PATCH("/:resourceID", impl.Update)
	impl.Router.PATCH("/:resourceID/quantity", impl.UpdateQuantity)
}

// @Summary find all resources
// @tags 	resource
// @Accept 	json
// @produce json
// @Success 200 {object} service.ResourcesResponse
// @Router 	/api/v1/resources [get]
func (impl *ResourceApiImpl) FindAll(c *gin.Context) {
	res, err := impl.ResourceService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary find resource by id
// @tags 	resource
// @Accept 	json
// @produce json
// Param 	id path			int true	"resource ID"
// @Success 200 {object} 	service.ResourceResponse
// Failure	404 {objetc}	HttpError
// @Router /api/v1/resources [get]
func (impl *ResourceApiImpl) FindOneByID(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid resourceID")
		return
	}

	res, err := impl.ResourceService.FindOneById(c, resourceID)
	if err != nil {
		if e, ok := err.(*exception.NotFoundException); ok {
			NewHttpError(c, http.StatusNotFound, e.Error())
		} else {
			NewHttpInternalServerError(c)
		}

		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	create a resource
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	resource		body	service.CreateResourceDto	true	"Create resource"
// @Success	201	{object}	service.ResourceResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources [post]
func (impl *ResourceApiImpl) Create(c *gin.Context) {
	var dto service.CreateResourceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := impl.ResourceService.Create(c, dto)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// @Summary	update a resource
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	id			path	int							true	"resource ID"
// @Param	resource	body	service.UpdateResourceDto	true	"Update resource"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources/{id} [patch]
func (impl *ResourceApiImpl) Update(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid resourceID")
		return
	}

	var dto service.UpdateResourceDto
	if err = c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := impl.ResourceService.Update(c, resourceID, dto); err != nil {
		if e, ok := err.(*exception.EmptyModelException); ok {
			NewHttpError(c, http.StatusBadRequest, e.Error())
		} else if e, ok := err.(*exception.NotFoundException); ok {
			NewHttpError(c, http.StatusNotFound, e.Error())
		} else {
			NewHttpInternalServerError(c)
		}

		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	update a resource quantity
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	id			path	int									true	"resource ID"
// @Param	resource	body	service.UpdateResourceQuantityDto	true	"Update resource quantity"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources/{id}/quantity [patch]
func (impl *ResourceApiImpl) UpdateQuantity(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid resourceID")
		return
	}

	var dto service.UpdateResourceQuantityDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := impl.ResourceService.UpdateQuantity(c, resourceID, dto); err != nil {
		if e, ok := err.(*exception.NotFoundException); ok {
			NewHttpError(c, http.StatusNotFound, e.Error())
		} else {
			NewHttpInternalServerError(c)
		}

		return
	}

	c.Status(http.StatusNoContent)
}
