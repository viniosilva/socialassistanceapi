package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

//go:generate mockgen -destination ../../mock/address_api_mock.go -package mock . AddressApi
type AddressApi interface {
	Configure()
}

type AddressApiImpl struct {
	Router         *gin.RouterGroup
	AddressService service.AddressService
}

func (impl *AddressApiImpl) Configure() {
	impl.Router.GET("", impl.FindAll)
	impl.Router.GET("/:addressID", impl.FindOneByID)
	impl.Router.POST("", impl.Create)
	impl.Router.PATCH("/:addressID", impl.Update)
	impl.Router.DELETE("/:addressID", impl.Delete)
}

// @Summary find all addresses
// @Tags address
// @Accept json
// @Produce json
// @Success 200 {object} service.AddressesResponse
// @Router /api/v1/addresses [get]
func (impl *AddressApiImpl) FindAll(c *gin.Context) {
	res, err := impl.AddressService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	find address by id
// @Tags	address
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"address ID"
// @Success	200	{object}	service.AddressesResponse
// @Failure	404	{object}	HttpError
// @Router	/api/v1/addresses/{id} [get]
func (impl *AddressApiImpl) FindOneByID(c *gin.Context) {
	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid addressID")
		return
	}

	res, err := impl.AddressService.FindOneById(c, addressID)
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

// @Summary	create an address
// @Tags	address
// @Accept	json
// @Produce	json
// @Param	address		body	service.CreateAddressDto	true	"Create address"
// @Success	201	{object}	service.AddressResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/addresses [post]
func (impl *AddressApiImpl) Create(c *gin.Context) {
	var dto service.CreateAddressDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := impl.AddressService.Create(c, dto)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// @Summary	update an address
// @Tags	address
// @Accept	json
// @Produce	json
// @Param	id			path	int							true	"address ID"
// @Param	address		body	service.UpdateAddressDto	true	"Update address"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/addresses/{id} [patch]
func (impl *AddressApiImpl) Update(c *gin.Context) {
	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid addressID")
		return
	}

	var dto service.UpdateAddressDto
	if err = c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = impl.AddressService.Update(c, addressID, dto); err != nil {
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

// @Summary	delete an address
// @Tags	address
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"address ID"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/addresses/{id} [delete]
func (impl *AddressApiImpl) Delete(c *gin.Context) {
	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid addressID")
		return
	}

	if err = impl.AddressService.Delete(c, addressID); err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.Status(http.StatusNoContent)
}
