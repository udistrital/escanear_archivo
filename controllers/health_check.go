package controllers

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/escanear_archivo/services"
)

func statusResponse(status string) map[string]interface{} {
	return map[string]interface{}{
		"Status": status,
	}
}

var defaultStatusResponse map[string]interface{} = statusResponse("Ok")

const defaultErrorString string = "UNHEALTHY_STATE"

func formatErrorResponse(errorMsg interface{}) map[string]interface{} {
	return statusResponse(fmt.Sprintf("%s -  %v", defaultErrorString, errorMsg))
}

// InitWithHandler accepts a (handler) function that, once performs the
// healthcheck, returns "nil" when everything is OK.
func InitWithHandler(statusCheckHandler func() (statusCheckError interface{})) {
	beego.Get("/", func(ctx *context.Context) {
		var responseError interface{}

		defer func() {

			// "catch"
			if err := recover(); err != nil {
				responseError = err
			}

			// "finally"
			response := defaultStatusResponse
			if responseError != nil {
				response = formatErrorResponse(responseError)
				ctx.Output.SetStatus(http.StatusServiceUnavailable) // 503
			}

			if err := services.PingClamd(); err != nil {
				logs.Critical(defaultErrorString, err)
				response = formatErrorResponse(err)
				ctx.Output.SetStatus(http.StatusServiceUnavailable)
			}

			ctx.Output.JSON(response, true, true)
		}()

		// "try"
		if statusCheckHandler != nil {
			if err := statusCheckHandler(); err != nil {
				responseError = err
			}
		}
	})
}

func InitHealthCheck() {
	InitWithHandler(nil)
}
