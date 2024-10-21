package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// CuposController operations for Cupos
type CuposController struct {
	beego.Controller
}

// URLMapping ...
func (c *CuposController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Cupos
// @Param	body		body 	models.Cupos	true		"body for Cupos content"
// @Success 201 {object} models.Cupos
// @Failure 403 body is empty
// @router / [post]
func (c *CuposController) Post() {
	data := c.Ctx.Input.RequestBody
	respuesta := services.PostCuposInscripcion(data)
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PostDocs ...
// @Title Create
// @Description create Cupos
// @Param	body		body 	models.Cupos	true		"body for Cupos content"
// @Success 201 {object} models.Cupos
// @Failure 403 body is empty
// @router /comentarios [post]
func (c *CuposController) PostDocs() {
	fmt.Println("postComentarios")
	defer errorhandler.HandlePanic(&c.Controller)
	data := c.Ctx.Input.RequestBody
	respuesta := services.PostDocCupos(data)
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Cupos by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Cupos
// @Failure 403 :id is empty
// @router /:id [get]
func (c *CuposController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Cupos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Cupos
// @Failure 403
// @router /:periodo/:proyecto/:tipo [get]
func (c *CuposController) GetAll() {
	defer errorhandler.HandlePanic(&c.Controller)

	periodo := c.Ctx.Input.Param(":periodo")
	proyecto := c.Ctx.Input.Param(":proyecto")
	tipo := c.Ctx.Input.Param(":tipo")

	respuesta := services.GetAllCuposInscripcion(periodo, proyecto, tipo)
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetAllDocs ...
// @Title GetAllDocs
// @Description get Cupos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Cupos
// @Failure 403
// @router /comentarios [get]
func (c *CuposController) GetAllDocs() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.GetAllDocCupos()
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Cupos
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Cupos	true		"body for Cupos content"
// @Success 200 {object} models.Cupos
// @Failure 403 :id is not int
// @router /:id [put]
func (c *CuposController) Put() {
	fmt.Println("Put")
	defer errorhandler.HandlePanic(&c.Controller)
	data := c.Ctx.Input.RequestBody
	respuesta := services.UpdateCuposInscripcion(data)
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Cupos
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CuposController) Delete() {

}
