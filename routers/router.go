// @APIVersion 1.0.0
// @Title Verificación externa de firma digital
// @Description Microservicio MID para verificar archivos con firma digital
package routers

import (
	"github.com/udistrital/funcion_lambda_verificar_archivo/controllers"
	"github.com/udistrital/utils_oas/errorhandler"

	"github.com/astaxie/beego"
)

func init() {
	beego.ErrorController(&errorhandler.ErrorHandlerController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/verificar",
			beego.NSInclude(
				&controllers.VerificarController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
