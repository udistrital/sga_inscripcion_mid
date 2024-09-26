package controllers

import (
	"github.com/astaxie/beego"

	//"github.com/astaxie/beego/httplib"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

type InscripcionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *InscripcionesController) URLMapping() {
	c.Mapping("PostInformacionFamiliar", c.PostInformacionFamiliar)
	c.Mapping("PostReintegro", c.PostReintegro)
	c.Mapping("PostTransferencia", c.PostTransferencia)
	c.Mapping("PostInfoIcfesColegio", c.PostInfoIcfesColegio)
	c.Mapping("PostPreinscripcion", c.PostPreinscripcion)
	c.Mapping("PostInfoComplementariaUniversidad", c.PostInfoComplementariaUniversidad)
	c.Mapping("PostInfoComplementariaTercero", c.PostInfoComplementariaTercero)
	c.Mapping("GetInfoComplementariaTercero", c.GetInfoComplementariaTercero)
	c.Mapping("PostInfoIcfesColegioNuevo", c.PostInfoIcfesColegioNuevo)
	c.Mapping("ConsultarProyectosEventos", c.ConsultarProyectosEventos)
	c.Mapping("ActualizarInfoContacto", c.ActualizarInfoContacto)
	c.Mapping("GetEstadoInscripcion", c.GetEstadoInscripcion)
	c.Mapping("PostGenerarInscripcion", c.PostGenerarInscripcion)
	c.Mapping("ActualizarEstadoEstudianteMatriculado", c.ActualizarEstadoEstudianteMatriculado)
	c.Mapping("ActualizarCuposAdmitidosOpcionado", c.ActualizarCuposAdmitidosOpcionado)
	c.Mapping("ActualizarEstadoInscripcion", c.ActualizarEstadoInscripcion)
}

// GetEstadoInscripcion ...
// @Title GetEstadoInscripcion
// @Description consultar los estados de todos los recibos generados por el tercero
// @Param	persona-id	query	string	false	"Id del tercero"
// @Param	id-periodo	query	string	false	"Id del ultimo periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /estado-recibos [get]
func (c *InscripcionesController) GetEstadoInscripcion() {

	defer errorhandler.HandlePanic(&c.Controller)

	terceroId := c.GetString("persona-id")
	idPeriodo := c.GetString("id-periodo")

	respuesta := services.EstadoInscripcion(terceroId, idPeriodo)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInformacionFamiliar ...
// @Title PostInformacionFamiliar
// @Description Agregar Información Familiar
// @Param   body        body    {}  true        "body Agregar PostInformacionFamiliar content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /informacion-familiar [post]
func (c *InscripcionesController) PostInformacionFamiliar() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InformacionFamiliar(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostReintegro ...
// @Title PostReintegro
// @Description Agregar Reintegro
// @Param   body        body    {}  true        "body Agregar Reintegro content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /reintegro [post]
func (c *InscripcionesController) PostReintegro() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.Reintegro(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostTransferencia ...
// @Title PostTransferencia
// @Description Agregar Transferencia
// @Param   body        body    {}  true        "body Agregar Transferencia content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /transferencia [post]
func (c *InscripcionesController) PostTransferencia() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.TransferenciaPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInfoIcfesColegio ...
// @Title PostInfoIcfesColegio
// @Description Agregar InfoIcfesColegio
// @Param   body        body    {}  true        "body Agregar InfoIcfesColegio content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /pruebas-de-estado/informacion/saber-once [post]
func (c *InscripcionesController) PostInfoIcfesColegio() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InfoIcfesColegio(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostPreinscripcion ...
// @Title PostPreinscripcion
// @Description Agregar Preinscripcion
// @Param   body        body    {}  true        "body Agregar Preinscripcion content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /preinscripcion [post]
func (c *InscripcionesController) PostPreinscripcion() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.PreinscripcionPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInfoIcfesColegioNuevo ...
// @Title PostInfoIcfesColegioNuevo
// @Description Agregar InfoIcfesColegio
// @Param   body        body    {}  true        "body Agregar InfoIcfesColegio content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /pruebas-de-estado/informacion/saber-once-nuevo [post]
func (c *InscripcionesController) PostInfoIcfesColegioNuevo() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InfoNuevoColegioIcfes(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInfoComplementariaUniversidad ...
// @Title PostInfoComplementariaUniversidad
// @Description Agregar InfoComplementariaUniversidad
// @Param   body        body    {}  true        "body Agregar InfoComplementariaUniversidad content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /informacion-complementaria/universidad [post]
func (c *InscripcionesController) PostInfoComplementariaUniversidad() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.PutInfoComplementaria(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// ConsultarProyectosEventos ...
// @Title ConsultarProyectosEventos
// @Description get ConsultarProyectosEventos by id
// @Param	evento_padre_id	path	int	true	"Id del Evento Padre"
// @Success 200 {}
// @Failure 404 not found resource
// @router /proyectos/eventos/:evento_padre_id [get]
func (c *InscripcionesController) ConsultarProyectosEventos() {

	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la persona
	idStr := c.Ctx.Input.Param(":evento_padre_id")

	respuesta := services.ConsultarEventos(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PostInfoComplementariaTercero ...
// @Title PostInfoComplementariaTercero
// @Description Agregar PostInfoComplementariaTercero
// @Param   body        body    {}  true        "body Agregar PostInfoComplementariaTercero content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /informacion-complementaria/tercero [post]
func (c *InscripcionesController) PostInfoComplementariaTercero() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InfoComplementariaTercero(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetInfoComplementariaTercero ...
// @Title GetInfoComplementariaTercero
// @Description consultar la información complementaria del tercero
// @Param	persona_id	path	true	"ID de la persona o tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router  /informacion-complementaria/tercero/:persona_id [get]
func (c *InscripcionesController) GetInfoComplementariaTercero() {

	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la persona
	persona_id := c.Ctx.Input.Param(":persona_id")

	respuesta := services.GetInfoCompTercero(persona_id)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// ActualizarInfoContacto ...
// @Title ActualizarInfoContacto
// @Description Actualiza los datos de contacto del tercero
// @Param	body	body 	{}	true		"body for Actualizar la info de contacto del tercero content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /informacion-complementaria/tercero [put]
func (c *InscripcionesController) ActualizarInfoContacto() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.ActualizarInfoContact(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostGenerarInscripcion ...
// @Title PostGenerarInscripcion
// @Description Registra una nueva inscripción con su respectivo recibo de pago
// @Param	body	body 	{}	true		"body for información de suministrada por el usuario par la inscripción"
// @Success 200 {}
// @Failure 403 body is empty
// @router /nueva [post]
func (c *InscripcionesController) PostGenerarInscripcion() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.GenerarInscripcion(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()

}

// ActualizarEstadoEstudianteMatriculado ...
// @Title ActualizarEstadoEstudianteMatriculado
// @Description Actualiza el estado de inscripción de un aspirante matriuclado
// @Param   body        body    {}  true		"body for actualizar estado estudiante matriculado content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /actualizar-estado-estudiante-matriculado [put]
func (c *InscripcionesController) ActualizarEstadoEstudianteMatriculado() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.ActualizarEstadoMatriculado(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// ActualizarCuposAdmitidosOpcionado ...
// @Title ActualizarCuposAdmitidosOpcionado
// @Description Actualiza los cupos, los admitidos y opcionados de un proyecto
// @Param   body        body    {}  true		"body for actualizar cupos admitidos opcionados content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /actualizar-cupos-admitidos-opcionados [put]
func (c *InscripcionesController) ActualizarCuposAdmitidosOpcionado() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.ActualizarCupos(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// ActualizarEstadoInscripcion ...
// @Title ActualizarEstadoInscripcion
// @Description Actualiza la inscripcion y coordina registro de historico en cambio de estado
// @Param   body        body    {}  true		"body for actualizar inscripcion"
// @Success 200 {int}
// @Failure 400 the request contains incorrect syntax
// @router /actualizar-inscripcion [post]
func (c *InscripcionesController) ActualizarEstadoInscripcion() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.ActualizarEstadoInscripcion(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
