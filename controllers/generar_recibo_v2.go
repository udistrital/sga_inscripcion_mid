package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// Generar_recibo_v2Controller operations for Generar_recibo_v2
type GenerarReciboV2Controller struct {
	beego.Controller
}

// URLMapping ...
func (c *GenerarReciboV2Controller) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Create
// @Description create Generar_recibo_v2
// @Param	body		body 	models.Generar_recibo_v2	true		"body for Generar_recibo_v2 content"
// @Success 201 {object} models.Generar_recibo_v2
// @Failure 403 body is empty
// @router / [post]
func (c *GenerarReciboV2Controller) Post() {
	fmt.Println("Post recibo v2")
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.GenerarReciboV2(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}
