package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/verificar_archivo_mid/services"
	"github.com/udistrital/verificar_archivo_mid/models"
	"github.com/udistrital/utils_oas/errorhandler"
	//"github.com/udistrital/utils_oas/requestresponse"
	//"fmt"
	"encoding/json"
)

// VerificarController operations for Verificar
type VerificarController struct {
	beego.Controller
}

// URLMapping ...
func (c *VerificarController) URLMapping() {
	c.Mapping("PostVerificar", c.PostVerificar)
}

// PostVerificar ...
// @Title PostVerificar
// @Description Verifica si un archivo tiene virus.
// @Param	body	body	models.VerificarRequest	true	"Datos de verificación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 body is empty or invalid
// @router / [post]
func (c *VerificarController) PostVerificar() {
	defer errorhandler.HandlePanic(&c.Controller)

	var payload models.VerificarRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &payload); err != nil {
		res := buildLambdaFormatResponse(400, services.LambdaResponse{
			Status:    "error",
			RawOutput: "Error al parsear el body: " + err.Error(),
		})
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if payload.PdfBase64 == "" {
		res := buildLambdaFormatResponse(400, services.LambdaResponse{
			Status:    "error",
			RawOutput: "No se encontraron datos PDF en la solicitud",
		})
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	resp, err := services.VerificarArchivo(payload.PdfBase64)
	if err != nil {
		res := buildLambdaFormatResponse(500, services.LambdaResponse{
			Status:    "error",
			RawOutput: "Error interno al verificar archivo: " + err.Error(),
		})
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res := buildLambdaFormatResponse(200, *resp)
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = res
	c.ServeJSON()
}

func buildLambdaFormatResponse(statusCode int, response services.LambdaResponse) map[string]interface{} {
	bodyBytes, _ := json.Marshal(response)
	return map[string]interface{}{
		"statusCode": statusCode,
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
		"body": string(bodyBytes),
	}
}
