package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// FormacionController ...
type FormacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *FormacionController) URLMapping() {
	c.Mapping("PostFormacionAcademica", c.PostFormacionAcademica)
	c.Mapping("PutFormacionAcademica", c.PutFormacionAcademica)
	c.Mapping("GetFormacionAcademica", c.GetFormacionAcademica)
	c.Mapping("GetFormacionAcademicaByTercero", c.GetFormacionAcademicaByTercero)
	c.Mapping("DeleteFormacionAcademica", c.DeleteFormacionAcademica)
	c.Mapping("GetInfoUniversidad", c.GetInfoUniversidad)
	c.Mapping("GetInfoUniversidadByNombre", c.GetInfoUniversidadByNombre)
	c.Mapping("PostTercero", c.PostTercero)
}

// PostFormacionAcademica ...
// @Title PostFormacionAcademica
// @Description Agregar Formacion Academica ud
// @Param   body        body    {}  true		"body Agregar Formacion Academica content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *FormacionController) PostFormacionAcademica() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.CrearFormacion(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetInfoUniversidad ...
// @Title GetInfoUniversidad
// @Description Obtener la información de la universidad por el nit
// @Param	Id		query 	int	true		"nit de la universidad"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /informacion-universidad/nit/:id [get]
func (c *FormacionController) GetInfoUniversidad() {

	defer errorhandler.HandlePanic(&c.Controller)

	//Numero del nit de la Universidad
	idStr := c.GetString("Id")

	respuesta := services.GetUniversidadInfo(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()

}

// GetInfoUniversidadByNombre ...
// @Title GetInfoUniversidadByNombre
// @Description Obtener la información de la universidad por el nombre
// @Param	nombre	query 	string	true		"nombre universidad"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /informacion-universidad/nombre/:nombre [get]
func (c *FormacionController) GetInfoUniversidadByNombre() {

	defer errorhandler.HandlePanic(&c.Controller)

	//Nombre de la Universidad
	idStr := c.GetString("nombre")

	respuesta := services.GetUniversidadInfo(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()

}

// PutFormacionAcademica ...
// @Title PutFormacionAcademica
// @Description Modificar Formacion Academica
// @Param	Id			query	int true		"Id del registro de formación"
// @Param	body		body 	{}	true		"body Modificar Formacion Academica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router / [put]
func (c *FormacionController) PutFormacionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	idFormacion := c.GetString("Id")

	data := c.Ctx.Input.RequestBody

	respuesta := services.ActualizarFormacionAcademica(idFormacion, data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetFormacionAcademica ...
// @Title GetFormacionAcademica
// @Description consultar Formacion Academica por id
// @Param	Id			query	int true		"Id del registro de formación"
// @Success 200 {}
// @Failure 404 not found resource
// @router /informacion-complementaria [get]
func (c *FormacionController) GetFormacionAcademica() {
	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.GetString("Id")

	respuesta := services.GetFormacionAcademicaById(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetFormacionAcademicaByTercero ...
// @Title GetFormacionAcademicaByTercero
// @Description consultar la Formacion Academica por id del tercero
// @Param	Id		query 	int	true		"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *FormacionController) GetFormacionAcademicaByTercero() {
	defer errorhandler.HandlePanic(&c.Controller)

	terceroId := c.GetString("Id")

	respuesta := services.GetFormacionAcademicaByIdTercero(terceroId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostTercero ...
// @Title PostTercero
// @Description Agregar nuevo tercero
// @Param   body        body    {}  true		"body Agregar nuevo tercero content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /tercero [post]
func (c *FormacionController) PostTercero() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.NuevoTercero(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()

}

// DeleteFormacionAcademica ...
// @Title DeleteFormacionAcademica
// @Description eliminar Formacion Academica por id de la formacion
// @Param	id		path 	int	true		"Id de la formacion academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *FormacionController) DeleteFormacionAcademica() {

	defer errorhandler.HandlePanic(&c.Controller)

	FormacionId := c.GetString("id")

	respuesta := services.EliminarFormacion(FormacionId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
