package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

//go:generate mockgen -destination ../../mock/donate_resource_api_mock.go -package mock . DonateResourceApi
type DonateResourceApi interface {
	Configure()
}

type DonateResourceApiImpl struct {
	Router                *gin.RouterGroup
	DonateResourceService service.DonateResourceService
	TraceMiddleware       func(c *gin.Context)
}

func (impl *DonateResourceApiImpl) Configure() {
	impl.Router.POST("/:resourceID/donate", impl.TraceMiddleware, impl.Donate)
	impl.Router.DELETE("/:resourceID/return", impl.TraceMiddleware, impl.Return)
}

// @Summary	donate a resource
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	id				path	int							true	"resource ID"
// @Param	resource		body	service.DonateResourceDto	true	"Donate a resource"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources/{id}/donate [post]
func (impl *DonateResourceApiImpl) Donate(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid resourceID")
		return
	}

	var dto service.DonateResourceDonateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}
	dto.ResourceID = resourceID

	err = impl.DonateResourceService.Donate(c, dto)
	if err != nil {
		if _, ok := err.(*exception.NegativeException); ok {
			NewHttpError(c, http.StatusBadRequest, err.Error())
		} else if _, ok := err.(*exception.NotFoundException); ok {
			NewHttpError(c, http.StatusNotFound, err.Error())
		} else {
			NewHttpInternalServerError(c)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Return a doneted resource
// @Tags	resource
// @Accept	json
// @Produce	json
// @Param	id				path	int							true	"resource ID"
// @Param	resource		body	service.DonateResourceDto	true	"Return a doneted resource"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/resources/{id}/return [delete]
func (impl *DonateResourceApiImpl) Return(c *gin.Context) {
	resourceID, err := strconv.Atoi(c.Param("resourceID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid resourceID")
		return
	}

	err = impl.DonateResourceService.Return(c, resourceID)
	if err != nil {
		if _, ok := err.(*exception.NegativeException); ok {
			NewHttpError(c, http.StatusBadRequest, err.Error())
		} else if _, ok := err.(*exception.NotFoundException); ok {
			NewHttpError(c, http.StatusNotFound, err.Error())
		} else {
			NewHttpInternalServerError(c)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
