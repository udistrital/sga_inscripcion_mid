package services

import (
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/udistrital/utils_oas/requestresponse"
)

func HandlePanic(c *beego.Controller) {
	if r := recover(); r != nil {
		logs.Error("Panic: ", r)
		debug.PrintStack()
		message := fmt.Sprintf("Error service %s: An internal server error occurred.", beego.AppConfig.String("appname"))
		message += fmt.Sprintf(" Request Info: URL: %s, Method: %s", c.Ctx.Request.URL, c.Ctx.Request.Method)
		message += " Time: " + time.Now().Format(time.RFC3339)
		statusCode := http.StatusInternalServerError
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = requestresponse.APIResponseDTO(false, statusCode, nil, message)
		c.ServeJSON()
	}
}

func GenerarCodigoBarras( data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var InformacionCodigo map[string]interface{}
	//alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(data, &InformacionCodigo); err == nil {
		fmt.Println(InformacionCodigo["Prueba"])
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, InformacionCodigo)
		//alertas = append(alertas, InformacionCodigo)

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

	} else {

		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil , err.Error())
		// alerta.Type = "error"
		// alerta.Code = "400"
		// alertas = append(alertas, err.Error())
	}

	return APIResponseDTO
}