package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

//go:generate mockgen -destination ../../mock/person_api_mock.go -package mock . PersonApi
type PersonApi interface {
	Configure()
}

type PersonApiImpl struct {
	Router          *gin.RouterGroup
	PersonService   service.PersonService
	TraceMiddleware func(c *gin.Context)
}

func (impl *PersonApiImpl) Configure() {
	impl.Router.GET("", impl.TraceMiddleware, impl.FindAll)
	impl.Router.GET("/:personID", impl.TraceMiddleware, impl.FindOneByID)
	impl.Router.POST("", impl.TraceMiddleware, impl.Create)
	impl.Router.PATCH("/:personID", impl.TraceMiddleware, impl.Update)
	impl.Router.DELETE("/:personID", impl.TraceMiddleware, impl.Delete)
}

// @Summary find all persons
// @Tags person
// @Accept json
// @Produce json
// @Success 200 {object} service.PersonResponse
// @Router /api/v1/persons [get]
func (impl *PersonApiImpl) FindAll(c *gin.Context) {
	res, err := impl.PersonService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	find person by id
// @Tags	person
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"person ID"
// @Success	200	{object}	service.PersonsResponse
// @Failure	404	{object}	HttpError
// @Router	/api/v1/persons/{id} [get]
func (impl *PersonApiImpl) FindOneByID(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	res, err := impl.PersonService.FindOneById(c, personID)
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

// @Summary	create a person
// @Tags	person
// @Accept	json
// @Produce	json
// @Param	person		body	service.CreatePersonDto	true	"Create person"
// @Success	201	{object}	service.PersonResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/persons [post]
func (impl *PersonApiImpl) Create(c *gin.Context) {
	var dto service.PersonCreateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := impl.PersonService.Create(c, dto)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// @Summary	update a person
// @Tags	person
// @Accept	json
// @Produce	json
// @Param	id			path	int					true	"person ID"
// @Param	person		body	service.UpdatePersonDto	true	"Update person"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/persons/{id} [patch]
func (impl *PersonApiImpl) Update(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	var dto service.PersonUpdateDto
	if err = c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}
	dto.ID = personID

	if err = impl.PersonService.Update(c, dto); err != nil {
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

// @Summary	delete a person
// @Tags	person
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"person ID"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/persons/{id} [delete]
func (impl *PersonApiImpl) Delete(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	if err = impl.PersonService.Delete(c, personID); err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.Status(http.StatusNoContent)
}
