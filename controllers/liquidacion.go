package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// LiquidacionController operations for Liquidacion
type LiquidacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *LiquidacionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
}

// Post ...
// @Title Create
// @Description create Liquidacion
// @Param	body		body 	models.Liquidacion	true		"body for Liquidacion content"
// @Success 201 {object} models.Liquidacion
// @Failure 403 body is empty
// @router / [post]
func (c *LiquidacionController) Post() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.GenerarReciboPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// Get ...
// @Title Get
// @Description get Liquidacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Liquidacion
// @Failure 403 :id is empty
// @router /:id [get]

func (c *LiquidacionController) Get() {
	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")

	respuesta := services.ValidarReciboGet(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}
