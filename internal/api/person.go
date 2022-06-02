package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// @Summary find all people
// @Tags person
// @Accept json
// @Produce json
// @Success 200 {object} service.PeopleResponse
// @Router /api/v1/people [get]
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
// @Success	200	{object}	service.PeopleResponse
// @Failure	404	{object}	HttpError
// @Router	/api/v1/people/{id} [get]
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
// @Router	/api/v1/people [post]
func (impl *PersonApi) Create(c *gin.Context) {
	var person service.PersonDto
	err := c.ShouldBindJSON(&person)
	if e, ok := err.(validator.ValidationErrors); ok {
		NewHttpError(c, http.StatusBadRequest, e.Error())
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
// @Router	/api/v1/people/{id} [patch]
func (impl *PersonApi) Update(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	var person service.PersonDto
	err = c.ShouldBindJSON(&person)
	if e, ok := err.(validator.ValidationErrors); ok {
		msg := e.Error()
		NewHttpError(c, http.StatusBadRequest, msg)
		return
	}

	res, err := impl.service.Update(c, personID, person)
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

// @Summary	delete a person
// @Tags	person
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"person ID"
// @Success	200	{object}	service.PersonResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/people/{id} [delete]
func (impl *PersonApi) Delete(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("personID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid personID")
		return
	}

	res, err := impl.service.Delete(c, personID)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, res)
}
