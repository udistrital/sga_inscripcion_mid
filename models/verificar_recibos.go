package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func VerificarRecibos(personaId string, periodoId string) ( resultadoAuxResponse map[string]interface{}, Error string) {
	var Inscripciones []map[string]interface{}
	var ReciboXML map[string]interface{}
	var resultadoAux []map[string]interface{}
	var resultado = make(map[string]interface{})
	var Estado string

	//Se consultan todas las inscripciones relacionadas a ese tercero
	errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=Activo:true,PersonaId:"+personaId+",PeriodoId:"+periodoId, &Inscripciones)
	if errInscripcion == nil {
		if Inscripciones != nil && fmt.Sprintf("%v", Inscripciones[0]) != "map[]" {
			fmt.Print(Inscripciones)
			// Ciclo for que recorre todas las inscripciones del tercero
			resultadoAux = make([]map[string]interface{}, len(Inscripciones))
			for i := 0; i < len(Inscripciones); i++ {
				if Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"] == "Transferencia interna" || Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"] == "Transferencia externa" || Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"] == "Reingreso" {
					Inscripciones = append(Inscripciones[:i], Inscripciones[i+1:]...)
					i = i - 1
				} else {
					ReciboInscripcion := fmt.Sprintf("%v", Inscripciones[i]["ReciboInscripcion"])
					if ReciboInscripcion != "0/<nil>" {
						errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ConsultarReciboJbpmService")+"consulta_recibo/"+ReciboInscripcion, &ReciboXML)
						if errRecibo == nil {
							if ReciboXML != nil && fmt.Sprintf("%v", ReciboXML) != "map[reciboCollection:map[]]" && fmt.Sprintf("%v", ReciboXML) != "map[]" {
								//Fecha límite de pago extraordinario
								FechaLimite := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_extraordinario"].(string)
								EstadoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["estado"].(string)
								PagoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["pago"].(string)
								//Verificación si el recibo de pago se encuentra activo y pago
								if EstadoRecibo == "A" && PagoRecibo == "S" {
									Estado = "Pago"
								} else {
									//Verifica si el recibo está vencido o no
									ATiempo, err := VerificarFechaLimite(FechaLimite)
									if err == nil {
										if ATiempo {
											Estado = "Pendiente pago"
										} else {
											Estado = "Vencido"
										}
									} else {
										Estado = "Vencido"
									}
								}

								resultadoAux[i] = map[string]interface{}{
									"Id":                  Inscripciones[i]["Id"],
									"ProgramaAcademicoId": Inscripciones[i]["ProgramaAcademicoId"],
									"ReciboInscripcion":   Inscripciones[i]["ReciboInscripcion"],
									"FechaCreacion":       Inscripciones[i]["FechaCreacion"],
									"Estado":              Estado,
									"EstadoInscripcion":   Inscripciones[i]["EstadoInscripcionId"].(map[string]interface{})["Nombre"],
								}
							} else {
								if fmt.Sprintf("%v", resultadoAux) != "map[]" {
									resultado["Inscripciones"] = resultadoAux
								} else {
									return resultado, "404"
									// errorGetAll = true
									// alertas = append(alertas, "No data found")
									// alerta.Code = "404"
									// alerta.Type = "error"
									// alerta.Body = alertas
									// c.Data["json"] = map[string]interface{}{"Response": alerta}
								}
							}
						} else {
							return resultado, "400"
							// errorGetAll = true
							// alertas = append(alertas, errRecibo.Error())
							// alerta.Code = "400"
							// alerta.Type = "error"
							// alerta.Body = alertas
							// c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					}
				}
			}

			for i := 0; i < len(resultadoAux); i++ {
				if resultadoAux[i] == nil {
					resultadoAux = append(resultadoAux[:i], resultadoAux[i+1:]...)
				}
			}

			resultado["Inscripciones"] = resultadoAux
		} else if (fmt.Sprintf("%v", Inscripciones[0]) == "map[]"){
			fmt.Println("Nueva inscripción")
		} else {
			return resultado, "404"
			// errorGetAll = true
			// alertas = append(alertas, "No data found")
			// alerta.Code = "404"
			// alerta.Type = "error"
			// alerta.Body = alertas
			// c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
	}else {
		return resultado, "400"
		// errorGetAll = true
		// alertas = append(alertas, errInscripcion.Error())
		// alerta.Code = "400"
		// alerta.Type = "error"
		// alerta.Body = alertas
		// c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	return resultado, Error
}