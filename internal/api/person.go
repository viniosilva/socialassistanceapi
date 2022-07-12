package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

type PersonApi struct {
	service *service.PersonService
}

func NewPersonApi(router *gin.RouterGroup, service *service.PersonService) *PersonApi {
	impl := &PersonApi{service}

	router.GET("", impl.FindAll)
	router.GET("/:personID", impl.FindOneByID)
	router.POST("", impl.Create)
	router.PATCH("/:personID", impl.Update)
	router.DELETE("/:personID", impl.Delete)

	return impl
}

// @Summary find all persons
// @Tags person
// @Accept json
// @Produce json
// @Success 200 {object} service.PersonResponse
// @Router /api/v1/persons [get]
func (impl *PersonApi) FindAll(c *gin.Context) {
	res, err := impl.service.FindAll(c)
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
func (impl *PersonApi) FindOneByID(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	res, err := impl.service.FindOneById(c, personID)
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

// @Summary	create a person
// @Tags	person
// @Accept	json
// @Produce	json
// @Param	person		body	service.PersonDto	true	"Create person"
// @Success	201	{object}	service.PersonResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/persons [post]
func (impl *PersonApi) Create(c *gin.Context) {
	var person service.PersonDto
	if err := c.ShouldBindJSON(&person); err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			NewHttpError(c, http.StatusBadRequest, e.Error())
			return
		}
		NewHttpError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	res, err := impl.service.Create(c, person)
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
// @Param	id				path	int					true	"person ID"
// @Param	person		body	service.PersonDto	true	"Update person"
// @Success	200	{object}	service.PersonResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/persons/{id} [patch]
func (impl *PersonApi) Update(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	var person service.PersonDto
	err = c.ShouldBindJSON(&person)
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	res, err := impl.service.Update(c, personID, person)
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

// @Summary	delete a person
// @Tags	person
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"person ID"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/persons/{id} [delete]
func (impl *PersonApi) Delete(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	err = impl.service.Delete(c, personID)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.Status(http.StatusNoContent)
}
