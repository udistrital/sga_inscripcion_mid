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
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /:tercero/:tipo_produccion [post]
func (c *SolicitudProduccionController) PostAlertSolicitudProduccion() {

	defer errorhandler.HandlePanic(&c.Controller)

	idTercero := c.Ctx.Input.Param(":tercero")
	idTipoProduccionSrt := c.Ctx.Input.Param(":tipo_produccion")
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
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /coincidencia/:id_solicitud/:id_coincidencia/:id_tercero [post]
func (c *SolicitudProduccionController) PostSolicitudEvaluacionCoincidencia() {
	
	defer errorhandler.HandlePanic(&c.Controller)
	
	idSolicitud := c.Ctx.Input.Param(":id_solicitud")
	idSolicitudCoincidencia := c.Ctx.Input.Param(":id_coincidencia")
	idTercero := c.Ctx.Input.Param(":id_tercero")

	data := c.Ctx.Input.RequestBody

	respuesta := services.PostSolicitudEvaluacion(idSolicitud, idSolicitudCoincidencia, idTercero, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
