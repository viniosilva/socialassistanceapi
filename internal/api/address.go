package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

type AddressApi struct {
	service *service.AddressService
}

func NewAddressApi(router *gin.RouterGroup, service *service.AddressService) *AddressApi {
	impl := &AddressApi{service}

	router.GET("", impl.FindAll)
	router.GET("/:addressID", impl.FindOneByID)
	router.POST("", impl.Create)
	router.PATCH("/:addressID", impl.Update)
	router.DELETE("/:addressID", impl.Delete)

	return impl
}

// @Summary find all addresses
// @Tags address
// @Accept json
// @Produce json
// @Success 200 {object} service.AddressesResponse
// @Router /api/v1/addresses [get]
func (impl *AddressApi) FindAll(c *gin.Context) {
	res, err := impl.service.FindAll(c)
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
func (impl *AddressApi) FindOneByID(c *gin.Context) {
	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid addressID")
		return
	}

	res, err := impl.service.FindOneById(c, addressID)
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

// @Summary	create an address
// @Tags	address
// @Accept	json
// @Produce	json
// @Param	address		body	service.AddressDto	true	"Create address"
// @Success	201	{object}	service.AddressResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/addresses [post]
func (impl *AddressApi) Create(c *gin.Context) {
	var address service.AddressDto
	err := c.ShouldBindJSON(&address)
	if e, ok := err.(validator.ValidationErrors); ok {
		NewHttpError(c, http.StatusBadRequest, e.Error())
		return
	}

	res, err := impl.service.Create(c, address)
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
// @Param	id				path	int					true	"address ID"
// @Param	address		body	service.AddressDto	true	"Update address"
// @Success	200	{object}	service.AddressResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/addresses/{id} [patch]
func (impl *AddressApi) Update(c *gin.Context) {
	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid addressID")
		return
	}

	var address service.AddressDto
	err = c.ShouldBindJSON(&address)
	if e, ok := err.(validator.ValidationErrors); e != nil && !ok {
		NewHttpError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	res, err := impl.service.Update(c, addressID, address)
	if err != nil {
		if e, ok := err.(*exception.EmptyModelException); ok {
			NewHttpError(c, http.StatusBadRequest, e.Error())
		} else if _, ok := err.(*exception.NotFoundException); ok {
			c.JSON(http.StatusNotFound, res)
		} else {
			NewHttpInternalServerError(c)
		}

		return
	}

	c.JSON(http.StatusOK, res)
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
func (impl *AddressApi) Delete(c *gin.Context) {
	addressID, err := strconv.Atoi(c.Param("addressID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid addressID")
		return
	}

	err = impl.service.Delete(c, addressID)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.Status(http.StatusNoContent)
}
