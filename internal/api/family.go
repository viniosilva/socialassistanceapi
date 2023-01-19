package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/service"
)

//go:generate mockgen -destination ../../mock/family_api_mock.go -package mock . FamilyApi
type FamilyApi interface {
	Configure()
}

type FamilyApiImpl struct {
	Router          *gin.RouterGroup
	FamilyService   service.FamilyService
	TraceMiddleware func(c *gin.Context)
	Addr            string
}

func (impl *FamilyApiImpl) Configure() {
	impl.Router.GET("", impl.TraceMiddleware, impl.FindAll)
	impl.Router.GET("/:familyID", impl.TraceMiddleware, impl.FindOneByID)
	impl.Router.POST("", impl.TraceMiddleware, impl.Create)
	impl.Router.PATCH("/:familyID", impl.TraceMiddleware, impl.Update)
	impl.Router.DELETE("/:familyID", impl.TraceMiddleware, impl.Delete)
}

// @Summary find all families
// @Tags family
// @Accept json
// @Produce json
// @Param limit query integer false "limit pagination"
// @Param offset query integer false "offset pagination"
// @Success 200 {object} service.FamiliesResponse
// @Router /api/v1/families [get]
func (impl *FamilyApiImpl) FindAll(c *gin.Context) {
	var p PaginationQuery
	if err := c.BindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, total, err := impl.FamilyService.FindAll(c, p.Limit, p.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	data := []Family{}
	for _, d := range res {
		data = append(data, *impl.Scan(d))
	}

	c.JSON(http.StatusOK, FamiliesResponse{
		PaginationResponse: PaginationResponse{
			Previous: BuildPreviousURL(impl.Addr, p.Limit, p.Offset),
			Next:     BuildNextURL(impl.Addr, p.Limit, p.Offset, total),
			Total:    total,
		},
		Data: data,
	})
}

// @Summary	find family by id
// @Tags	family
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"family ID"
// @Success	200	{object}	service.FamiliesResponse
// @Failure	404	{object}	HttpError
// @Router	/api/v1/families/{id} [get]
func (impl *FamilyApiImpl) FindOneByID(c *gin.Context) {
	familyID, err := strconv.Atoi(c.Param("familyID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid familyID")
		return
	}

	res, err := impl.FamilyService.FindOneById(c, familyID)
	if err != nil {
		if e, ok := err.(*exception.NotFoundException); ok {
			NewHttpError(c, http.StatusNotFound, e.Error())
		} else {
			NewHttpInternalServerError(c)
		}

		return
	}

	c.JSON(http.StatusOK, FamilyResponse{Data: impl.Scan(*res)})
}

// @Summary	create an family
// @Tags	family
// @Accept	json
// @Produce	json
// @Param	family		body	service.FamilyCreateDto	true	"Create family"
// @Success	201	{object}	service.FamilyResponse
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/families [post]
func (impl *FamilyApiImpl) Create(c *gin.Context) {
	var dto service.FamilyCreateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := impl.FamilyService.Create(c, dto)
	if err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, FamilyResponse{Data: impl.Scan(*res)})
}

// @Summary	update an family
// @Tags	family
// @Accept	json
// @Produce	json
// @Param	id			path	int							true	"family ID"
// @Param	family		body	service.FamilyUpdateDto		true	"Update family"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/families/{id} [patch]
func (impl *FamilyApiImpl) Update(c *gin.Context) {
	familyID, err := strconv.Atoi(c.Param("familyID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid familyID")
		return
	}

	var dto service.FamilyUpdateDto
	if err = c.ShouldBindJSON(&dto); err != nil {
		NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}
	dto.ID = familyID

	if err = impl.FamilyService.Update(c, dto); err != nil {
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

// @Summary	delete an family
// @Tags	family
// @Accept	json
// @Produce	json
// @Param	id	path		int	true	"family ID"
// @Success	204
// @Failure	400	{object}	HttpError
// @Failure	500	{object}	HttpError
// @Router	/api/v1/families/{id} [delete]
func (impl *FamilyApiImpl) Delete(c *gin.Context) {
	familyID, err := strconv.Atoi(c.Param("familyID"))
	if err != nil {
		NewHttpError(c, http.StatusBadRequest, "invalid familyID")
		return
	}

	if err = impl.FamilyService.Delete(c, familyID); err != nil {
		NewHttpInternalServerError(c)
		return
	}

	c.Status(http.StatusNoContent)
}

func (impl *FamilyApiImpl) Scan(data model.Family) *Family {
	return &Family{
		ID:           data.ID,
		CreatedAt:    data.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:    data.UpdatedAt.Format("2006-01-02T15:04:05"),
		Name:         data.Name,
		Country:      data.Country,
		State:        data.State,
		City:         data.City,
		Neighborhood: data.Neighborhood,
		Street:       data.Street,
		Number:       data.Number,
		Complement:   data.Complement,
		Zipcode:      data.Zipcode,
	}
}
