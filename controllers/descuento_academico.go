package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// DescuentoController ...
type DescuentoController struct {
	beego.Controller
}

// URLMapping ...
func (c *DescuentoController) URLMapping() {
	c.Mapping("PostDescuentoAcademico", c.PostDescuentoAcademico)
	c.Mapping("PutDescuentoAcademico", c.PutDescuentoAcademico)
	c.Mapping("GetDescuentoAcademico", c.GetDescuentoAcademico)
	c.Mapping("GetDescuentoAcademicoByPersona", c.GetDescuentoAcademicoByPersona)
	c.Mapping("GetDescuentoByPersonaPeriodoDependencia", c.GetDescuentoByPersonaPeriodoDependencia)
	c.Mapping("GetDescuentoAcademicoByDependenciaID", c.GetDescuentoAcademicoByDependenciaID)
}

// PostDescuentoAcademico ...
// @Title PostDescuentoAcademico
// @Description Agregar Descuento Academico
// @Param	body		body 	{}	true		"body Agregar Descuento Academico content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *DescuentoController) PostDescuentoAcademico() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.SolicitarDescuentoAcademico(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()

}

// PutDescuentoAcademico ...
// @Title PutDescuentoAcademico
// @Description Modificar Descuento Academico
// @Param	id	path 	int	true		"el id de la solicitud de descuento a modificar"
// @Param	body		body 	{}	true		"body Modificar Descuento Academico content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *DescuentoController) PutDescuentoAcademico() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	idStr := c.Ctx.Input.Param(":id")

	respuesta := services.ActualizarDescuentoAcademico(data, idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetDescuentoAcademico ...
// @Title GetDescuentoAcademico
// @Description consultar Descuento Academico por userid
// @Param	PersonaId		query 	int	true		"Id de la persona"
// @Param	SolicitudId		query 	int	true		"Id de la solicitud"
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *DescuentoController) GetDescuentoAcademico() {

	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la persona
	idStr := c.GetString("PersonaId")

	//Id de la solicitud
	idSolitudDes := c.GetString("SolicitudId")

	respuesta := services.GetDescuentoAcademicoById(idStr, idSolitudDes)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetDescuentoAcademicoByDependenciaID ...
// @Title GetDescuentoAcademicoByDependenciaID
// @Description consultar Descuento Academico por DependenciaId
// @Param	dependencia_id		path 	int	true		"DependenciaId"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:dependencia_id [get]
func (c *DescuentoController) GetDescuentoAcademicoByDependenciaID() {

	defer errorhandler.HandlePanic(&c.Controller)
	//Id de la persona
	idStr := c.Ctx.Input.Param(":dependencia_id")

	respuesta := services.GetDescuentoByDpendencia(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetDescuentoAcademicoByPersona ...
// @Title GetDescuentoAcademicoByPersona
// @Description consultar Descuento Academico por userid
// @Param	persona_id		path 	int	true		"Id de la persona"
// @Success 200 {}
// @Failure 404 not found resource
// @router /persona/:persona_id [get]
func (c *DescuentoController) GetDescuentoAcademicoByPersona() {
	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la persona
	idStr := c.Ctx.Input.Param(":persona_id")

	respuesta := services.GetDescuentoAcademicoByTercero(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
	c.ServeJSON()
}

// GetDescuentoByPersonaPeriodoDependencia ...
// @Title GetDescuentoByPersonaPeriodoDependencia
// @Description consultar Descuento Academico por userid
// @Param	PersonaId		query 	int	true		"Id de la persona"
// @Param	DependenciaId		query 	int	true		"Id de la dependencia"
// @Param	PeriodoId		query 	int	true		"Id del periodo académico"
// @Success 200 {}
// @Failure 404 not found resource
// @router /detalle [get]
func (c *DescuentoController) GetDescuentoByPersonaPeriodoDependencia() {
	defer errorhandler.HandlePanic(&c.Controller)

	//Captura de parámetros
	idPersona := c.GetString("PersonaId")
	idDependencia := c.GetString("DependenciaId")
	idPeriodo := c.GetString("PeriodoId")

	respuesta := services.GetDescuentoByTerceroPeriodoDependencia(idPersona, idPeriodo, idDependencia)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}
