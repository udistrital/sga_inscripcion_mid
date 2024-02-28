package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// GeneradorCodigoBarrasController ...
type GeneradorCodigoBarrasController struct {
	beego.Controller
}

// URLMapping ...
func (c *GeneradorCodigoBarrasController) URLMapping() {
	c.Mapping("GenerarCodigoBarras", c.GenerarCodigoBarras)
}

// GenerarCodigoBarras ...
// @Title GenerarCodigoBarras
// @Description Creacion de codigo de barras
// @Param   body        body    {}  true        "body Agregar ProduccionAcademica content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *GeneradorCodigoBarrasController) GenerarCodigoBarras() {

	defer errorhandler.HandlePanic(&c.Controller)
	
	data := c.Ctx.Input.RequestBody

	respuesta := services.GenerarCodigoBarras(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}