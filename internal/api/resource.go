package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

type ResourceApi struct {
	service *service.ResourceService
}

func NewResourceApi(router *gin.RouterGroup, service *service.ResourceService) *ResourceApi {
	impl := &ResourceApi{service}

	router.GET("", impl.FindAll)
	router.GET("/:resourceID", impl.FindOneByID)
	router.POST("", impl.Create)
	router.PATCH("/:resourceID", impl.Update)
	router.DELETE("/:resourceID", impl.Delete)

	return impl
}

// @Summary find all resources
// @tags 	resource
// @Accept 	json
// @produce json
// @Success 200 {object} service.ResourcesResponse
// @Router 	/api/v1/resources [get]
func (impl *ResourceApi) FindAll(c *gin.Context) {
	res, err := impl.service.FindAll(c)
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
func (impl *ResourceApi) FindOneByID(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "Invalid resourceID")
		return
	}

	res, err := impl.service.FindOneById(c, resourceID)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	if res.Data == nil {
		c.JSON(http.StatusNotFound, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	create a resource
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	resource		body	service.ResourceDto	true	"Create resource"
// @Success	201	{object}	service.ResourceResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources [post]
func (impl *ResourceApi) Create(c *gin.Context) {
	var resource service.ResourceDto
	err := c.ShouldBindJSON(&resource)
	if e, ok := err.(validator.ValidationErrors); ok {
		msg := e.Error()
		NewHttpError(c, http.StatusBadRequest, msg)
		return
	}

	res, err := impl.service.Create(c, resource)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	update a resource
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	id				path	int			true	"resource ID"
// @Param	resource	body	service.ResourceDto	true	"Update resource"
// @Success	200	{object}	service.ResourceResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources/{id} [patch]
func (impl *ResourceApi) Update(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid resourceID")
		return
	}

	var resource service.ResourceDto
	err = c.ShouldBindJSON(&resource)
	if e, ok := err.(validator.ValidationErrors); ok {
		msg := e.Error()
		NewHttpError(c, http.StatusBadRequest, msg)
		return
	}

	res, err := impl.service.Update(c, resourceID, resource)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}
	if res.Data == nil {
		c.JSON(http.StatusNotFound, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	delete a resource
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"resource ID"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources/{id} [delete]
func (impl *ResourceApi) Delete(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid resourceID")
		return
	}

	err = impl.service.Delete(c, resourceID)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.Status(http.StatusNoContent)
}
