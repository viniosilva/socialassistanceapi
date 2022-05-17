package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

type CustomerApi struct {
	service *service.CustomerService
}

func NewCustomerApi(router *gin.RouterGroup, service *service.CustomerService) *CustomerApi {
	impl := &CustomerApi{service}

	router.GET("", impl.FindAll)
	router.GET("/:customerID", impl.FindOneByID)
	router.POST("", impl.Create)
	router.PATCH("/:customerID", impl.Update)

	return impl
}

// @Summary find all customers
// @Tags customer
// @Accept json
// @Produce json
// @Success 200 {object} service.CustomersResponse
// @Router /api/v1/customers [get]
func (impl *CustomerApi) FindAll(c *gin.Context) {
	res, err := impl.service.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	find customer by id
// @Tags	customer
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"Customer ID"
// @Success	200	{object}	service.CustomersResponse
// @Failure	404	{object}	HttpError
// @Router	/api/v1/customers/{id} [get]
func (impl *CustomerApi) FindOneByID(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid customerID")
		return
	}

	res, err := impl.service.FindOneById(c, customerID)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	if res.Data == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	create a customer
// @Tags	customer
// @Accept	json
// @Produce	json
// @Param	customer		body	service.CustomerDto	true	"Create customer"
// @Success	201	{object}	service.CustomerResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/customers [post]
func (impl *CustomerApi) Create(c *gin.Context) {
	var customer service.CustomerDto
	err := c.ShouldBindJSON(&customer)
	if e, ok := err.(validator.ValidationErrors); ok {
		msg := e.Error()
		NewHttpError(c, http.StatusBadRequest, msg)
		return
	}

	res, err := impl.service.Create(c, customer)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// @Summary	update a customer
// @Tags	customer
// @Accept	json
// @Produce	json
// @Param	customer		body	service.CustomerDto	true	"Update customer"
// @Success	200	{object}	service.CustomerResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/customers [patch]
func (impl *CustomerApi) Update(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid customerID")
		return
	}

	var customer service.CustomerDto
	err = c.ShouldBindJSON(&customer)
	if e, ok := err.(validator.ValidationErrors); ok {
		msg := e.Error()
		NewHttpError(c, http.StatusBadRequest, msg)
		return
	}

	res, err := impl.service.Update(c, customerID, customer)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}
	if res.Data == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, res)
}
