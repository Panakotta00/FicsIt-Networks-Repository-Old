package util

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

func GetDefaultInt(c echo.Context, param string, fallback int) int {
	val, err := strconv.Atoi(c.QueryParam(param))
	if err != nil {
		return fallback
	}
	return val
}

func GetIntRange(c echo.Context, param string, min int, max int, fallback int) int {
	val := GetDefaultInt(c, param, fallback)
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

func GetPagination(c echo.Context) (int, int) {
	page := GetDefaultInt(c, "page", 0)
	if page < 0 {
		page = 0
	}
	count := GetIntRange(c, "count", 1, 100, 50)
	return page, count
}

func GetSnowflake(c echo.Context, param string) (int64, error) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	return id, err
}
