package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// ExperienciaLaboralController ...
type ExperienciaLaboralController struct {
	beego.Controller
}

// URLMapping ...
func (c *ExperienciaLaboralController) URLMapping() {
	c.Mapping("PostExperienciaLaboral", c.PostExperienciaLaboral)
	c.Mapping("PutExperienciaLaboral", c.PutExperienciaLaboral)
	c.Mapping("GetExperienciaLaboral", c.GetExperienciaLaboral)
	c.Mapping("GetInformacionEmpresa", c.GetInformacionEmpresa)
	c.Mapping("GetExperienciaLaboralByTercero", c.GetExperienciaLaboralByTercero)
	c.Mapping("DeleteExperienciaLaboral", c.DeleteExperienciaLaboral)
}

// PostExperienciaLaboral ...
// @Title PostExperienciaLaboral
// @Description Agregar Formacion Academica ud
// @Param   body        body    {}  true		"body Agregar Experiencia Laboral content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *ExperienciaLaboralController) PostExperienciaLaboral() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.CreateExperienciaLaboral(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetInformacionEmpresa ...
// @Title GetInformacionEmpresa
// @Description Obtener la información de la empresa por el nit
// @Param	Id		query 	int	true		"nit de la empresa"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /informacion_empresa/ [get]
func (c *ExperienciaLaboralController) GetInformacionEmpresa() {

	
	defer errorhandler.HandlePanic(&c.Controller)

	//Numero del nit de la empresa
	idStr := c.GetString("Id")


	respuesta := services.GetInfoEmpresa(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetExperienciaLaboralByTercero ...
// @Title GetExperienciaLaboralByTercero
// @Description Obtener la información de la empresa por el nit
// @Param	Id		query 	int	true		"nit de la empresa"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /by_tercero/ [get]
func (c *ExperienciaLaboralController) GetExperienciaLaboralByTercero() {
	defer errorhandler.HandlePanic(&c.Controller)

	terceroID := c.GetString("Id")

	respuesta := services.GetExperienciaLaboralByPersona(terceroID)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PutExperienciaLaboral ...
// @Title PutExperienciaLaboral
// @Description Modificar Formacion Academica ud
// @Param   body        body    {}  true		"body Agregar Experiencia Laboral content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *ExperienciaLaboralController) PutExperienciaLaboral() {
	
	defer errorhandler.HandlePanic(&c.Controller)

	idTercero := c.GetString(":id")
	
	data := c.Ctx.Input.RequestBody

	respuesta := services.ActualizarExperienciaLaboral(idTercero, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetExperienciaLaboral ...
// @Title GetExperienciaLaboral
// @Description consultar Experiencia Laboral por id
// @Param	id		path 	int	true		"Id de la experiencia"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *ExperienciaLaboralController) GetExperienciaLaboral() {	
	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la experiencia
	idStr := c.Ctx.Input.Param(":id")

	respuesta := services.GetExperienciaLaboralById(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// DeleteExperienciaLaboral ...
// @Title DeleteExperienciaLaboral
// @Description eliminar Experiencia Laboral por id
// @Param   id      path    int  true        "Id de la Experiencia Laboral"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *ExperienciaLaboralController) DeleteExperienciaLaboral() {
	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la experiencia
	idStr := c.Ctx.Input.Param(":id")

	respuesta := services.DeleteExperienciaById(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}
