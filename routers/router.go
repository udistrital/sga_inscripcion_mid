// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/utils_oas/errorhandler"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/controllers"
)

func init() {
	beego.ErrorController(&errorhandler.ErrorHandlerController{})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/academico/descuento",
			beego.NSInclude(
				&controllers.DescuentoController{},
			),
		),
		beego.NSNamespace("/academico/formacion",
			beego.NSInclude(
				&controllers.FormacionController{},
			),
		),
		beego.NSNamespace("/academico/produccion",
			beego.NSInclude(
				&controllers.ProduccionAcademicaController{},
			),
		),
		beego.NSNamespace("/experiencia-laboral",
			beego.NSInclude(
				&controllers.ExperienciaLaboralController{},
			),
		),
		beego.NSNamespace("/codigo",
			beego.NSInclude(
				&controllers.GeneradorCodigoBarrasController{},
			),
		),
		beego.NSNamespace("/recibos",
			beego.NSInclude(
				&controllers.GenerarReciboController{},
			),
		),
		beego.NSNamespace("/inscripciones",
			beego.NSInclude(
				&controllers.InscripcionesController{},
			),
		),
		beego.NSNamespace("/solicitud-produccion",
			beego.NSInclude(
				&controllers.SolicitudProduccionController{},
			),
		),
		beego.NSNamespace("/transferencia",
			beego.NSInclude(
				&controllers.Transferencia_reingresoController{},
			),
		),
		beego.NSNamespace("/cupos",
			beego.NSInclude(
				&controllers.CuposController{},
			),
		),
		beego.NSNamespace("/time_bog",
			beego.NSInclude(
				&controllers.Time_bogController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
