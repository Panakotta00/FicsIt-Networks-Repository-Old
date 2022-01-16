package Util

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
