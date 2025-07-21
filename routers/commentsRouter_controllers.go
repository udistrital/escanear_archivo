package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/funcion_lambda_verificar_archivo/controllers:VerificarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/funcion_lambda_verificar_archivo/controllers:VerificarController"],
		beego.ControllerComments{
			Method:           "PostVerificar",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
