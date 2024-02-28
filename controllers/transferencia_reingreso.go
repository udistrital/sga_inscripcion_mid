package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"

	"github.com/udistrital/utils_oas/errorhandler"

)

// Transferencia_reingresoController operations for Transferencia_reingreso
type Transferencia_reingresoController struct {
	beego.Controller
}

// URLMapping ...
func (c *Transferencia_reingresoController) URLMapping() {
	c.Mapping("Post", c.PostSolicitud)
	c.Mapping("Put", c.PutInfoSolicitud)
	c.Mapping("PutInscripcion", c.PutInscripcion)
	c.Mapping("PutSolicitud", c.PutSolicitud)
	c.Mapping("GetInscripcion", c.GetInscripcion)
	c.Mapping("GetSolicitudesInscripcion", c.GetSolicitudesInscripcion)
	c.Mapping("GetConsultarPeriodo", c.GetConsultarPeriodo)
	c.Mapping("GetConsultarParametros", c.GetConsultarParametros)
	c.Mapping("GetEstados", c.GetEstados)
}

// PostSolicitud ...
// @Title Create
// @Description create Transferencia_reingreso
// @Param	body		body 	helpers.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 201 {object} helpers.Transferencia_reingreso
// @Failure 403 body is empty
// @router / [post]
func (c *Transferencia_reingresoController) PostSolicitud() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.SolicitudPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PutInfoSolicitud ...
// @Title Create
// @Description create Transferencia_reingreso
// @Param	body		body 	helpers.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 201 {object} helpers.Transferencia_reingreso
// @Failure 403 body is empty
// @Failure 404 not found resource
// @router /:id [put]
func (c *Transferencia_reingresoController) PutInfoSolicitud() {
	
	defer errorhandler.HandlePanic(&c.Controller)
	
	id_solicitud := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.PutSolicitudInfo(id_solicitud, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PutInscripcion ...
// @Title PutInscripcion
// @Description crear la inscripción y actualizar solicitud
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	helpers.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 200 {object} helpers.Transferencia_reingreso
// @Failure 400 the request contains incorrect syntax
// @router /actualizar_estado/:id [put]
func (c *Transferencia_reingresoController) PutInscripcion() {
	
	defer errorhandler.HandlePanic(&c.Controller)
	
	id_solicitud := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.PutInscripcion(id_solicitud, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PutSolicitud ...
// @Title PutSolicitud
// @Description crear la inscripción y actualizar solicitud
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	helpers.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 200 {object} helpers.Transferencia_reingreso
// @Failure 400 the request contains incorrect syntax
// @router /respuesta_solicitud/:id [put]
func (c *Transferencia_reingresoController) PutSolicitud() {

	defer errorhandler.HandlePanic(&c.Controller)
	
	id_solicitud := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.SolicitudPut(id_solicitud, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetInscripcion ...
// @Title GetInscripcion
// @Description get Transferencia_reingreso by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} helpers.Transferencia_reingreso
// @Failure 403 :id is empty
// @router /inscripcion/:id [get]
func (c *Transferencia_reingresoController) GetInscripcion() {
	
	defer errorhandler.HandlePanic(&c.Controller)
	
	idInscripcion := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.GetInscripcionById(idInscripcion, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetSolicitudesInscripcion ...
// @Title GetSolicitudesInscripcion
// @Description get Transferencia_reingreso by id
// @Success 200 {object} helpers.Transferencia_reingreso
// @router /solicitudes/ [get]
func (c *Transferencia_reingresoController) GetSolicitudesInscripcion() {
	defer errorhandler.HandlePanic(&c.Controller)

	respuesta := services.GetSolicitudes()

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetConsultarPeriodo ...
// @Title GetConsultarPeriodo
// @Description get información necesaria para crear un solicitud de transferencias
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_periodo/ [get]
func (c *Transferencia_reingresoController) GetConsultarPeriodo() {
	
	defer errorhandler.HandlePanic(&c.Controller)

	respuesta := services.ConsultarPeriodo()

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetConsultarParametros ...
// @Title GetConsultarParametros
// @Description get información necesaria para crear un solicitud de transferencias
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_parametros/:id_calendario/:persona_id [get]
func (c *Transferencia_reingresoController) GetConsultarParametros() {
	defer errorhandler.HandlePanic(&c.Controller)
	
	idCalendario := c.Ctx.Input.Param(":id_calendario")
	idPersona := c.Ctx.Input.Param(":persona_id")

	respuesta := services.ConsultarParametros(idCalendario, idPersona)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetEstadoInscripcion ...
// @Title GetEstadoInscripcion
// @Description consultar los estados de todos los recibos generados por el tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Success 200 {}
// @Failure 403 body is empty
// @Failure 404 not found resource
// @Failure 400 not found resource
// @router /estado_recibos/:persona_id [get]
func (c *Transferencia_reingresoController) GetEstadoInscripcion() {

	defer errorhandler.HandlePanic(&c.Controller)
	
	persona_id := c.Ctx.Input.Param(":persona_id")

	respuesta := services.EstadoInscripcionGet(persona_id)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetEstados ...
// @Title GetEstados
// @Description get Transferencia_reingreso by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} helpers.Transferencia_reingreso
// @Failure 403 :id is empty
// @router /estados [get]
func (c *Transferencia_reingresoController) GetEstados() {
	defer errorhandler.HandlePanic(&c.Controller)

	respuesta := services.EstadosGet()

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
