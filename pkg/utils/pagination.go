package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultSize = 10
)

type PaginationQuery struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// Set page size
func (q *PaginationQuery) SetSize(limitQuery string) error {
	if limitQuery == "" {
		q.Limit = defaultSize
		return nil
	}
	n, err := strconv.Atoi(limitQuery)
	if err != nil {
		return err
	}
	q.Limit = n

	return nil
}

// Set page number
func (q *PaginationQuery) SetPage(offsetQuery string) error {
	if offsetQuery == "" {
		q.Offset = 0
		return nil
	}
	n, err := strconv.Atoi(offsetQuery)
	if err != nil {
		return err
	}
	q.Offset = n

	return nil
}

// Get pagination query struct from
func GetPaginationFromCtx(c *gin.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.Query("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.Query("size")); err != nil {
		return nil, err
	}

	return q, nil
}
