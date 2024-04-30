package controllers

import (
	"github.com/astaxie/beego"
)

// Generar_recibo_v2Controller operations for Generar_recibo_v2
type Generar_recibo_v2Controller struct {
	beego.Controller
}

// URLMapping ...
func (c *Generar_recibo_v2Controller) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Generar_recibo_v2
// @Param	body		body 	models.Generar_recibo_v2	true		"body for Generar_recibo_v2 content"
// @Success 201 {object} models.Generar_recibo_v2
// @Failure 403 body is empty
// @router / [post]
func (c *Generar_recibo_v2Controller) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Generar_recibo_v2 by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Generar_recibo_v2
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Generar_recibo_v2Controller) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Generar_recibo_v2
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Generar_recibo_v2
// @Failure 403
// @router / [get]
func (c *Generar_recibo_v2Controller) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Generar_recibo_v2
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Generar_recibo_v2	true		"body for Generar_recibo_v2 content"
// @Success 200 {object} models.Generar_recibo_v2
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Generar_recibo_v2Controller) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Generar_recibo_v2
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Generar_recibo_v2Controller) Delete() {

}
