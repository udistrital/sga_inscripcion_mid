package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"

)

type GenerarReciboController struct {
	beego.Controller
}

// URLMapping ...
func (c *GenerarReciboController) URLMapping() {
	c.Mapping("PostGenerarRecibo", c.PostGenerarRecibo)
	c.Mapping("PostGenerarEstudianteRecibo", c.PostGenerarEstudianteRecibo)
	c.Mapping("PostGenerarComprobanteInscripcion", c.PostGenerarComprobanteInscripcion)
}

// PostGenerarEstudianteRecibo ...
// @Title PostGenerarEstudianteRecibo
// @Description Genera un recibo de pago
// @Param	body		body 	{}	true		"body Datos del recibo content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /estudiantes/ [post]
func (c *GenerarReciboController) PostGenerarEstudianteRecibo() {

	defer errorhandler.HandlePanic(&c.Controller)
	
	data := c.Ctx.Input.RequestBody

	respuesta := services.GenerarReciboPago(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PostGenerarRecibo ...
// @Title PostGenerarRecibo
// @Description Genera un recibo de pago
// @Param	body		body 	{}	true		"body Datos del recibo content"
// @Success 200 {}
// @Failure 400 body is empty
// @router / [post]
func (c *GenerarReciboController) PostGenerarRecibo() {

	defer errorhandler.HandlePanic(&c.Controller)
	
	data := c.Ctx.Input.RequestBody

	respuesta := services.GenerarReciboPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PostGenerarComprobanteInscripcion ...
// @Title PostGenerarComprobanteInscripcion
// @Description Genera un comprobante de inscripcion
// @Param	body		body 	{}	true		"Informacion para el comprobante"
// @Success 200 {}
// @Failure 400 body is empty
// @router /comprobante_inscripcion/ [post]
func (c *GenerarReciboController) PostGenerarComprobanteInscripcion() {
	defer errorhandler.HandlePanic(&c.Controller)
	
	data := c.Ctx.Input.RequestBody

	respuesta := services.GenerarComprobante(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}