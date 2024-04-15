package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/utils"
)

func RegistrarDoc(documento []map[string]interface{}) (status interface{}, outputError interface{}) {

	var resultadoRegistro map[string]interface{}

	fmt.Println("http://" + beego.AppConfig.String("GestorDocumental") + "document/")
	errRegDoc := utils.SendJson("http://"+beego.AppConfig.String("GestorDocumental")+"document/uploadAnyFormat", "POST", &resultadoRegistro, documento)

	fmt.Println(errRegDoc)
	if resultadoRegistro["Status"].(string) == "200" && errRegDoc == nil {
		return resultadoRegistro["res"], nil
	} else {
		return nil, resultadoRegistro["Error"].(string)
	}
}
