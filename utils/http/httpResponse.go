package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	errorStatus "go-store/utils/errors"
)

func SendResponse(c *gin.Context, data interface{}, handleErr error) {
	// w.Header().Add("Content-Type", "application/json")
	// if data == nil {
	// 	log.Info("Data in 'sendResponse' is null")
	// 	return
	// }
	var status int
	var message string
	if handleErr != nil {

		var apiErr errorStatus.APIError
		if errors.As(handleErr, &apiErr) {
			status, message = apiErr.APIError()
		} else {
			status = http.StatusInternalServerError
			message = "internal server Error"
		}
		c.JSON(status, message)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
