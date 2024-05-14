package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/utils"
	"github.com/udistrital/utils_oas/request"
)

func SetInactivo(url string) (exito bool) {
	exito = false
	var payload1 map[string]interface{}
	fmt.Println(url)
	errGet := request.GetJson(url, &payload1)
	if errGet == nil {
		fmt.Println(payload1)
		var idDisable string = ""
		var body map[string]interface{}
		if payload1["Id"] != nil {
			fmt.Println("is by id only")
			idDisable = fmt.Sprintf("%v", payload1["Id"])
			body = payload1
		}
		if payload1["Data"] != nil {
			fmt.Println("is is inside data")
			idDisable = fmt.Sprintf("%v", payload1["Data"].(map[string]interface{})["Id"])
			body = payload1["Data"].(map[string]interface{})
		}

		fmt.Println("id is:", idDisable)

		if idDisable != "" {
			body["Activo"] = false
			fmt.Println("body is:", body)
			var payload2 map[string]interface{}
			errSet := request.SendJson(url, "PUT", &payload2, body)
			if errSet == nil {
				if payload2["Id"] != nil {
					if fmt.Sprintf("%v", payload2["Id"]) == idDisable {
						exito = true
					} else {
						exito = false
					}
				} else if payload1["Data"] != nil {
					if fmt.Sprintf("%v", payload2["Data"].(map[string]interface{})["Id"]) == idDisable {
						exito = true
					} else {
						exito = false
					}
				} else {
					exito = false
				}
			} else {
				exito = false
			}
		} else {
			exito = false
		}
	} else {
		exito = false
	}

	return exito
}

// IdInfoCompTercero is ...
func IdInfoCompTercero(grupo string, codAbrev string) (Id string, ok bool) {
	var resp []map[string]interface{}
	errResp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=GrupoInfoComplementariaId__Id:"+grupo+",CodigoAbreviacion:"+codAbrev+"&fields=Id", &resp)
	if errResp == nil && fmt.Sprintf("%v", resp) != "[map[]]" {
		Id = fmt.Sprintf("%v", resp[0]["Id"].(float64))
		ok = true
	} else {
		Id = "0"
		ok = false
	}
	return Id, ok
}

// Verificar estado de lso recibos ...
func VerificarRecibos(personaId string, periodoId string) (resultadoAuxResponse map[string]interface{}, Error string) {
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
									ATiempo, err := utils.VerificarFechaLimite(FechaLimite)
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

								//Estado = "Pago"

								resultadoAux[i] = map[string]interface{}{
									"Id": Inscripciones[i]["Id"],
									//"ProgramaAcademicoId": Inscripciones[i]["ProgramaAcademicoId"],
									"ReciboInscripcion": Inscripciones[i]["ReciboInscripcion"],
									"FechaCreacion":     Inscripciones[i]["FechaCreacion"],
									"Estado":            Estado,
									"EstadoInscripcion": Inscripciones[i]["EstadoInscripcionId"].(map[string]interface{})["Nombre"],
								}
							} else {
								if fmt.Sprintf("%v", resultadoAux) != "map[]" {
									resultado["Inscripciones"] = resultadoAux
								} else {
									return resultado, "404"
								}
							}
						} else {
							return resultado, "400"
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
		} else if fmt.Sprintf("%v", Inscripciones[0]) == "map[]" {

			fmt.Println("Nueva inscripción")

		} else {
			return resultado, "404"
		}
	} else {
		return resultado, "400"
	}

	return resultado, Error
}

// Generacion credencial inscripciones pregrado
func GenerarCredencialInscripcionPregrado(periodoId float64) (credencial int) {
	var parametros []map[string]interface{}
	periodoIdInt := int(periodoId)

	errParam := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?limit=1&query=PeriodoId:"+fmt.Sprintf("%d", periodoIdInt)+"&fields=Credencial&sortby=Credencial&order=desc", &parametros)
	if errParam == nil {

		credencialMaxima := parametros[0]["Credencial"].(float64)

		if credencialMaxima > 100 {
			credencial = int(credencialMaxima + 1)
		} else {
			credencial = 101
		}

		return credencial
	} else {
		return 0
	}
}
