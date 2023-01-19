package api

import (
	"fmt"
	"net/url"
)

type PaginationQuery struct {
	Limit  int `form:"limit,default=10" example:"10" binding:"lte=50"`
	Offset int `form:"offset,default=0" example:"10" binding:"gte=0"`
}

type PaginationResponse struct {
	Previous string `json:"previous" example:"localhost:8080/api/v1/families?limit=10&offset=0"`
	Next     string `json:"next" example:"localhost:8080/api/v1/families?limit=10&offset=20"`
	Total    int    `json:"total" example:"100"`
}

func BuildPreviousURL(addr string, limit, offset int) string {
	if offset == 0 {
		return ""
	}

	url, err := url.Parse(addr)
	if err != nil {
		return ""
	}

	newOffset := offset - limit

	q := url.Query()
	q.Add("limit", fmt.Sprint(limit))

	if newOffset > 0 {
		q.Add("offset", fmt.Sprint(newOffset))
	}

	url.RawQuery = q.Encode()

	return url.String()
}

func BuildNextURL(addr string, limit, offset, total int) string {
	if offset+limit >= total {
		return ""
	}

	url, err := url.Parse(addr)
	if err != nil {
		return ""
	}

	q := url.Query()
	q.Add("limit", fmt.Sprint(limit))
	q.Add("offset", fmt.Sprint(offset+limit))

	url.RawQuery = q.Encode()

	return url.String()
}
