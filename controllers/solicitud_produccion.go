package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// SolicitudProduccionController ...
type SolicitudProduccionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudProduccionController) URLMapping() {
	c.Mapping("PostAlertSolicitudProduccion", c.PostAlertSolicitudProduccion)
	c.Mapping("PostSolicitudEvaluacionCoincidencia", c.PostSolicitudEvaluacionCoincidencia)
	c.Mapping("PutResultadoSolicitud", c.PutResultadoSolicitud)
}

// PostAlertSolicitudProduccion ...
// @Title PostAlertSolicitudProduccion
// @Description Agregar Alerta en Solicitud docente en casos necesarios
// @Param   body    body    {}  true        "body Agregar SolicitudProduccion content"
// @Param	tercero	query	        string	false	"Id del tercero"
// @Param	tipo-produccion	query	string	false	"Id del tipo de produccion"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /alerta-docente [post]
func (c *SolicitudProduccionController) PostAlertSolicitudProduccion() {

	defer errorhandler.HandlePanic(&c.Controller)

	idTercero := c.GetString("persona")
	idTipoProduccionSrt := c.GetString("tipo-produccion")
	data := c.Ctx.Input.RequestBody

	respuesta := services.PostAlertSolicitud(idTercero, idTipoProduccionSrt, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PutResultadoSolicitud ...
// @Title PutResultadoSolicitud
// @Description Modificar resultaado solicitud docente
// @Param	id		path 	int	true		"el id de la produccion"
// @Param   body        body    {}  true        "body Modificar resultado en produccionAcaemica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *SolicitudProduccionController) PutResultadoSolicitud() {

	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.PutResultado(idStr, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostSolicitudEvaluacionCoincidencia ...
// @Title PostSolicitudEvaluacionCoincidencia
// @Description Agregar Alerta en Solicitud docente en casos necesarios
// @Param   body    body    {}  true        "body Agregar SolicitudProduccion content"
// @Param	id-solicitud	query	string	false	"Se recibe parametro Id de la solicitud"
// @Param	id-coincidencia	query	string	false	"Se recibe parametro Id de la coincidencia"
// @Param	id-tercero	query	string	false	"Se recibe parametro Id del tercero"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /coincidencias [post]
func (c *SolicitudProduccionController) PostSolicitudEvaluacionCoincidencia() {

	defer errorhandler.HandlePanic(&c.Controller)

	idSolicitud := c.GetString("id-solicitud")
	idSolicitudCoincidencia := c.GetString("id-coincidencia")
	idTercero := c.GetString("id-tercero")

	data := c.Ctx.Input.RequestBody

	respuesta := services.PostSolicitudEvaluacion(idSolicitud, idSolicitudCoincidencia, idTercero, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
