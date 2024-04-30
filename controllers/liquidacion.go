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
