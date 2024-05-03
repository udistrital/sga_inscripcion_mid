package controllers

import (
	"fmt"

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
	c.Mapping("Post", c.PostBaseLegalizacionMatricula)
	c.Mapping("GetInfoLegalizacionMatricula", c.GetInfoLegalizacionMatricula)
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

// GetInfoLegalizacionMatricula ...
// @Title GetInfoLegalizacionMatricula
// @Description consultar la información complementaria del tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router  /informacion-legalizacion/:persona_id [get]
func (c *LegalizacionController) GetInfoLegalizacionMatricula() {
	defer errorhandler.HandlePanic(&c.Controller)

	//Id de la persona
	persona_id := c.Ctx.Input.Param(":persona_id")
	fmt.Println("PERSONA ID:")
	fmt.Println(persona_id)

	respuesta := services.GetInfoLegalizacionTercero(persona_id)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}
