package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/helpers"
)

// Time_bogController operations for Time_bog
type Time_bogController struct {
	beego.Controller
}

// URLMapping ...
func (c *Time_bogController) URLMapping() {
	c.Mapping("GetTimeBog", c.GetTimeBog)
}

// GetTimeBog ...
// @Title GetTimeBog
// @Description get Time_bog
// @Success 200 {object} models.Time_bog
// @Failure 500 something bad happened
// @router / [get]
func (c *Time_bogController) GetTimeBog() {
	respuesta := helpers.GetTimeBog()
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = respuesta
	c.ServeJSON()
}
