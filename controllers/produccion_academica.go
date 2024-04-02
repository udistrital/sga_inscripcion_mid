package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// ProduccionAcademicaController ...
type ProduccionAcademicaController struct {
	beego.Controller
}

// URLMapping ...
func (c *ProduccionAcademicaController) URLMapping() {
	c.Mapping("PostProduccionAcademica", c.PostProduccionAcademica)
	c.Mapping("PutProduccionAcademica", c.PutProduccionAcademica)
	c.Mapping("GetAllProduccionAcademica", c.GetAllProduccionAcademica)
	c.Mapping("GetOneProduccionAcademica", c.GetOneProduccionAcademica)
	c.Mapping("GetProduccionAcademica", c.GetProduccionAcademica)
	c.Mapping("GetIdProduccionAcademica", c.GetIdProduccionAcademica)
	c.Mapping("DeleteProduccionAcademica", c.DeleteProduccionAcademica)
	c.Mapping("PutEstadoAutorProduccionAcademica", c.PutEstadoAutorProduccionAcademica)
}

// PostProduccionAcademica ...
// @Title PostProduccionAcademica
// @Description Agregar Producción academica
// @Param   body    body    {}  true        "body Agregar ProduccionAcademica content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *ProduccionAcademicaController) PostProduccionAcademica() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.ProduccionAcademicaPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PutEstadoAutorProduccionAcademica ...
// @Title PutEstadoAutorProduccionAcademica
// @Description Modificar Estado de Autor de Producción Academica
// @Param	id		path 	int	true		"el id del autor a modificar"
// @Param   body        body    {}  true        "body Modificar AutorProduccionAcademica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /autor/:id [put]
func (c *ProduccionAcademicaController) PutEstadoAutorProduccionAcademica() {
	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.EstadoAutorProduccion(idStr, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PutProduccionAcademica ...
// @Title PutProduccionAcademica
// @Description Modificar Producción Academica
// @Param	id		path 	int	true		"el id de la Produccion academica a modificar"
// @Param   body        body    {}  true        "body Modificar ProduccionAcademica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *ProduccionAcademicaController) PutProduccionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.ProduccionAcademicaPut(idStr, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()

}

// GetOneProduccionAcademica ...
// @Title GetOneProduccionAcademica
// @Description consultar Produccion Academica por id
// @Param   id      path    int  true        "Id"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *ProduccionAcademicaController) GetOneProduccionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la producción
	idProduccion := c.Ctx.Input.Param(":id")

	respuesta := services.GetProduccionById(idProduccion)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetAllProduccionAcademica ...
// @Title GetAllProduccionAcademica
// @Description consultar todas las Producciones académicas
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *ProduccionAcademicaController) GetAllProduccionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	respuesta := services.GetAllProducciones()

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetIdProduccionAcademica ...
// @Title GetIdProduccionAcademica
// @Description consultar Produccion Academica por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:tercero [get]
func (c *ProduccionAcademicaController) GetIdProduccionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	idTercero := c.Ctx.Input.Param(":tercero")

	respuesta := services.GetIdProduccion(idTercero)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetProduccionAcademica ...
// @Title GetProduccionAcademica
// @Description consultar Produccion Academica por tercero
// @Param   tercero path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /tercero/:tercero [get]
func (c *ProduccionAcademicaController) GetProduccionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	idTercero := c.Ctx.Input.Param(":tercero")

	respuesta := services.GetProduccion(idTercero)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// DeleteProduccionAcademica ...
// @Title DeleteProduccionAcademica
// @Description eliminar Produccion Academica por id
// @Param   id      path    int  true        "Id de la Produccion Academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *ProduccionAcademicaController) DeleteProduccionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")

	respuesta := services.DeleteProduccion(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
