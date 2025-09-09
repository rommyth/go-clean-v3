package response

import "github.com/labstack/echo/v4"

type Response struct {
	StatusCode int		`json:"status_code"`
	Data interface{} `json:"data,omitempty"`
	Error string `json:"error,omitempty"`

}

func JSON(c echo.Context, status int, data interface{}) error {
response := Response{
		StatusCode: status,
		Data: data,
	}
	return c.JSON(status, response)
}

func Error(c echo.Context, status int, message string, err error) error {
	response := Response{
		StatusCode: status,
		Error: message,
	}

	return c.JSON(status, response)
}