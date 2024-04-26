package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// LegalizacionController operations for Legalizacion
type LegalizacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *LegalizacionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// PostBaseLegalizacionMatricula ...
// @Title PostBaseLegalizacionMatricula
// @Description create Legalizacion
// @Param   body        body    {}  true		"body for Legalizacion content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /base [post]
func (c *LegalizacionController) PostBaseLegalizacionMatricula() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.CrearInfolegalizacionMatricula(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Legalizacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Legalizacion
// @Failure 403 :id is empty
// @router /:id [get]
func (c *LegalizacionController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Legalizacion
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Legalizacion
// @Failure 403
// @router / [get]
func (c *LegalizacionController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Legalizacion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Legalizacion	true		"body for Legalizacion content"
// @Success 200 {object} models.Legalizacion
// @Failure 403 :id is not int
// @router /:id [put]
func (c *LegalizacionController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Legalizacion
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *LegalizacionController) Delete() {

}
