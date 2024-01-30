package controllers

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
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
	var InformacionCodigo map[string]interface{}
	//var alerta models.Alert
	//alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InformacionCodigo); err == nil {
		fmt.Println(InformacionCodigo["Prueba"])

		CodigoRecibido := InformacionCodigo["Prueba"].(string)
		fmt.Println("Generando code128 barcode para : ", CodigoRecibido)
		bcode, _ := code128.Encode(CodigoRecibido)

		if err != nil {
			fmt.Printf("String %s cannot be encoded", CodigoRecibido)
			os.Exit(1)
		}

		// Scale the barcode to 500x200 pixels
		ScCode, _ := barcode.Scale(bcode, 400, 40)

		// create the output file
		file, _ := os.Create("Codigo_generado.png")
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, ScCode)

		fmt.Println("Code128 code generated and saved to Codigo_generado.png")

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": InformacionCodigo}

	} else {
		logs.Error(err)
		c.Data["message"] = "Error service GenerarCodigoBarras: " + err.Error()
		c.Abort("400")
	}
	c.ServeJSON()
}
