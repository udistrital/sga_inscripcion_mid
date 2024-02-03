package helpers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_inscripcion/utils"
)

func RegistrarDoc(documento []map[string]interface{}) (status interface{}, outputError interface{}) {

	var resultadoRegistro map[string]interface{}

	errRegDoc := utils.SendJson("http://"+beego.AppConfig.String("GestorDocumental")+"document/upload", "POST", &resultadoRegistro, documento)

	if resultadoRegistro["Status"].(string) == "200" && errRegDoc == nil {
		return resultadoRegistro["res"], nil
	} else {
		return nil, resultadoRegistro["Error"].(string)
	}
}
