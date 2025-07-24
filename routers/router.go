// @APIVersion 1.0.0
// @Title Verificación de virus en archivos
// @Description Microservicio MID para verificar virus en archivos
package routers

import (
	"github.com/udistrital/escanear_archivo/controllers"
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
