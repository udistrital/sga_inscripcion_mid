package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"golang.org/x/sync/errgroup"

	"github.com/k0kubun/pp"
)

// Funcion para solicitar descuentos academicos
func SolicitarDescuentoAcademico(data []byte) (APIResponseDTO requestresponse.APIResponse) {

	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var solicitud map[string]interface{}
	var solicitudPost map[string]interface{}
	var tipoDescuento []map[string]interface{}
	alertas := []interface{}{}

	if err := json.Unmarshal(data, &solicitud); err == nil {
		IDTipoDescuento := fmt.Sprintf("%v", solicitud["DescuentosDependenciaId"].(map[string]interface{})["Id"])
		IDDependencia := fmt.Sprintf("%v", solicitud["DescuentosDependenciaId"].(map[string]interface{})["Dependencia"])
		errDescuentosDependencia := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia?query=TipoDescuentoId__Id:"+IDTipoDescuento+",DependenciaId:"+IDDependencia, &tipoDescuento)
		if errDescuentosDependencia == nil && fmt.Sprintf("%v", tipoDescuento[0]["System"]) != "map[]" {

			// DescuentosDependenciaID := map[string]interface{}{
			// 	"Activo":          solicitud["DescuentosDependenciaId"].(map[string]interface{})["Activo"],
			// 	"DependenciaId":   solicitud["DescuentosDependenciaId"].(map[string]interface{})["Dependencia"],
			// 	"PeriodoId":       solicitud["DescuentosDependenciaId"].(map[string]interface{})["Periodo"],
			// 	"TipoDescuentoId": tipoDescuento,
			// }

			solicituddescuento := map[string]interface{}{
				"Id":                      0,
				"TerceroId":               solicitud["PersonaId"],
				"Estado":                  "Por aprobar",
				"PeriodoId":               solicitud["PeriodoId"],
				"Activo":                  true,
				"DescuentosDependenciaId": tipoDescuento[0],
			}
			formatdata.JsonPrint(solicituddescuento)
			// fmt.Println(solicituddescuento)

			errSolicitud := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento", "POST", &solicitudPost, solicituddescuento)
			if errSolicitud == nil && fmt.Sprintf("%v", solicitudPost["System"]) != "map[]" && solicitudPost["Id"] != nil {
				if solicitudPost["Status"] != 400 {
					//soporte de descuento
					var soporte map[string]interface{}

					soportedescuento := map[string]interface{}{
						"SolicitudDescuentoId": solicitudPost,
						"Activo":               true,
						"DocumentoId":          solicitud["DocumentoId"],
					}

					errSoporte := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento", "POST", &soporte, soportedescuento)
					if errSoporte == nil && fmt.Sprintf("%v", soporte["System"]) != "map[]" && soporte["Id"] != nil {
						if soporte["Status"] != 400 {
							resultado = map[string]interface{}{"Id": solicitudPost["Id"], "PersonaId": solicitudPost["PersonaId"], "Estado": solicitudPost["Estado"], "PeriodoId": solicitudPost["PeriodoId"], "DescuentosDependenciaId": solicitudPost["DescuentosDependenciaId"]}
							resultado["DocumentoId"] = soporte["DocumentoId"]
							APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)

						} else {
							//resultado solicitud de descuento
							var resultado2 map[string]interface{}
							request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/%.f", solicitudPost["Id"]), "DELETE", &resultado2, nil)
							logs.Error(errSoporte)
							//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
							alertas = append(alertas, soporte)
						}
					} else {
						logs.Error(errSoporte)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						alertas = append(alertas, soporte)
					}
				} else {
					logs.Error(errSolicitud)
					//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
					alertas = append(alertas, solicitudPost)
				}
			} else {
				logs.Error(errSolicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				alertas = append(alertas, solicitudPost)
			}
		} else {
			logs.Error(errDescuentosDependencia)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			alertas = append(alertas, errDescuentosDependencia)
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		alertas = append(alertas, err)
	}

	if len(alertas) > 0 {
		return requestresponse.APIResponseDTO(false, 400, nil, alertas)
	} else {
		return APIResponseDTO
	}

}

func ActualizarDescuentoAcademico(data []byte, id string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var solicitud map[string]interface{}

	if err := json.Unmarshal(data, &solicitud); err == nil {
		//soporte de descuento
		var soporte []map[string]interface{}
		var soportePut map[string]interface{}

		errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=Activo:true,SolicitudDescuentoId:"+id, &soporte)
		if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
			if soporte[0]["Status"] != 404 {
				soporte[0]["DocumentoId"] = solicitud["DocumentoId"]

				errSoportePut := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/"+
					fmt.Sprintf("%v", soporte[0]["Id"]), "PUT", &soportePut, soporte[0])
				if errSoportePut == nil && fmt.Sprintf("%v", soportePut["System"]) != "map[]" && soportePut["Id"] != nil {
					if soportePut["Status"] != 400 {
						resultado = solicitud
						APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
					} else {
						logs.Error(errSoportePut)
						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errSoportePut)
						return APIResponseDTO
					}
				} else {
					logs.Error(errSoportePut)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoportePut)
					return APIResponseDTO
				}

			} else {
				if soporte[0]["Message"] == "Not found resource" {
					APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil, "No data found")
				} else {
					logs.Error(soporte)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporte)
					return APIResponseDTO
				}
			}
		} else {
			logs.Error(soporte)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporte)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func GetDescuentoAcademicoById(idTercero string, idSolicitud string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado consulta
	var resultado map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}
	validData := []interface{}{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento?query=TerceroId:"+idTercero+",Id:"+idSolicitud+"&fields=Id,TerceroId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)

	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if len(solicitud[0]) >= 1 {
			resultado = solicitud[0]

			//resultado descuento dependencia
			var descuento map[string]interface{}
			errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+fmt.Sprintf("%v", solicitud[0]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
			if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
				if descuento["Status"] != 404 {
					//resultado tipo descuento
					var tipo map[string]interface{}

					errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)

					if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
						if tipo["Status"] != 404 {
							descuento["TipoDescuentoId"] = tipo
							resultado["DescuentosDependenciaId"] = descuento

							//resultado soporte descuento
							var soporte []map[string]interface{}

							errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+idSolicitud+"&fields=DocumentoId", &soporte)

							if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
								if soporte[0]["Status"] != 404 {
									pp.Println("$$$$$$$$$$$$$$$$$$$$$$$$$")
									//fmt.Println("el resultado de los documentos es: ", resultado4)
									resultado["DocumentoId"] = soporte[0]["DocumentoId"]
									validData = append(validData, resultado)
									APIResponseDTO = requestresponse.APIResponseDTO(true, 200, validData)

								} else {
									if soporte[0]["Message"] == "Not found resource" {
										APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
										return APIResponseDTO
									} else {
										logs.Error(soporte)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporte)
										return APIResponseDTO
									}
								}
							} else {
								logs.Error(soporte)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporte)
								return APIResponseDTO
							}
						} else {
							if tipo["Message"] == "Not found resource" {
								validData = append(validData, nil)
							} else {
								logs.Error(tipo)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errTipo)
								return APIResponseDTO
							}
						}
					} else {
						logs.Error(tipo)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errTipo)
						return APIResponseDTO
					}
				} else {
					if descuento["Message"] == "Not found resource" {
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
						return APIResponseDTO
					} else {
						logs.Error(descuento)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDescuento)
						return APIResponseDTO
					}
				}
			} else {
				logs.Error(descuento)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDescuento)
				return APIResponseDTO
			}
		} else {
			if solicitud[0]["Message"] == "Not found resource" {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				return APIResponseDTO
			} else {
				logs.Error(solicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSolicitud)
				return APIResponseDTO
			}
		}
	} else {
		logs.Error(solicitud)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSolicitud)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func GetDescuentoByDpendencia(idDependencia string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado consulta
	var resultados []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}
	//var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})
	//Definición de el group para las gorutines
	wge := new(errgroup.Group)

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia?limit=0&query=Activo:true,DependenciaId:"+idDependencia, &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {

			for _, solici := range  solicitud{
				wge.Go(func () error{
					fmt.Println("Entra hilo")
					var tipoDescuento map[string]interface{}
					errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", solici["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipoDescuento)
					if errDescuento == nil && fmt.Sprintf("%v", tipoDescuento["System"]) != "map[]" {
						resultados = append(resultados, tipoDescuento)
					} else {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDescuento.Error())
					}
					return errDescuento
				})
			}
			//Si existe error, se realiza
			if err := wge.Wait(); err != nil {
				errorGetAll = true
			}
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSolicitud.Error())
		alertas = append(alertas, errSolicitud.Error())
	}
	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultados, nil)
		return APIResponseDTO
	} else {
		return APIResponseDTO
	}
}

func GetDescuentoAcademicoByTercero(idTercero string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado solicitud descuento
	var resultado []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/?query=PersonaId:"+idTercero+"&fields=Id,PersonaId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {

			for u := 0; u < len(solicitud); u++ {
				//resultado solicitud descuento
				var descuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+
					fmt.Sprintf("%v", solicitud[u]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
				if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
					if descuento["Status"] != 404 {
						//resultado tipo descuento
						var tipo map[string]interface{}
						errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
						if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
							if tipo["Status"] != 404 {
								descuento["TipoDescuentoId"] = tipo
								solicitud[u]["DescuentosDependenciaId"] = descuento

								//resultado soporte descuento
								var soporte []map[string]interface{}
								errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+fmt.Sprintf("%v", solicitud[u]["Id"])+"&fields=DocumentoId", &soporte)
								if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
									if soporte[0]["Status"] != 404 {
										//fmt.Println("el resultado de los documentos es: ", resultado4)
										solicitud[u]["DocumentoId"] = soporte[0]["DocumentoId"]
									} else {
										if soporte[0]["Message"] == "Not found resource" {
											APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
											return APIResponseDTO
										} else {
											logs.Error(soporte)
											//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
											APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporte)
											return APIResponseDTO
										}
									}
								} else {
									logs.Error(soporte)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporte)
									return APIResponseDTO
								}
							} else {
								if tipo["Message"] == "Not found resource" {
									APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
									return APIResponseDTO
								} else {
									logs.Error(tipo)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errTipo)
									return APIResponseDTO
								}
							}
						} else {
							logs.Error(tipo)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errTipo)
							return APIResponseDTO
						}
					} else {
						if descuento["Message"] == "Not found resource" {
							APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
							return APIResponseDTO
						} else {
							logs.Error(descuento)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDescuento)
							return APIResponseDTO
						}
					}
				} else {
					logs.Error(descuento)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDescuento)
					return APIResponseDTO
				}
			}
			resultado = solicitud
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		} else {
			if solicitud[0]["Message"] == "Not found resource" {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				return APIResponseDTO
			} else {
				logs.Error(solicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSolicitud)
				return APIResponseDTO
			}
		}
	} else {
		logs.Error(solicitud)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSolicitud)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func GetDescuentoByTerceroPeriodoDependencia(idTercero string, idPeriodo string, idDependencia string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado solicitud descuento
	var resultado []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}
	var errorGetAll bool

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento?query=Activo:true,TerceroId:"+idTercero+",PeriodoId:"+idPeriodo+",DescuentosDependenciaId.DependenciaId:"+idDependencia+"&fields=Id,TerceroId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {
			for u := 0; u < len(solicitud); u++ {
				//resultado solicitud descuento
				var descuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+
					fmt.Sprintf("%v", solicitud[u]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
				if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
					if descuento["Status"] != 404 {
						//resultado tipo descuento
						var tipo map[string]interface{}
						errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
						if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
							if tipo["Status"] != 404 {
								descuento["TipoDescuentoId"] = tipo
								solicitud[u]["DescuentosDependenciaId"] = descuento

								//resultado soporte descuento
								var soporte []map[string]interface{}
								errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=Activo:true,SolicitudDescuentoId:"+fmt.Sprintf("%v", solicitud[u]["Id"])+"&fields=DocumentoId", &soporte)
								if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
									if soporte[0]["Status"] != 404 {
										//fmt.Println("el resultado de los documentos es: ", resultado4)
										solicitud[u]["DocumentoId"] = soporte[0]["DocumentoId"]
									}
								} else {
									errorGetAll = true
									APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errSoporte.Error())
								}
							}
						} else {
							errorGetAll = true
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errTipo.Error())
						}
					}
				} else {
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errDescuento.Error())
				}
			}
			resultado = solicitud
			// c.Data["json"] = resultado
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errSolicitud.Error())
	}
	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 200, resultado, nil)
		return APIResponseDTO
	} else {
		return APIResponseDTO
	}
}
