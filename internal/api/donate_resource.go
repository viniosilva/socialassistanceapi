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
}

func (impl *DonateResourceApiImpl) Configure() {
	impl.Router.POST("/:resourceID/donate", impl.Donate)
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
