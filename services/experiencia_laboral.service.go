package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func CreateExperienciaLaboral(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var ExperienciaLaboral map[string]interface{}
	var idInfoExperencia string
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})

	if err := json.Unmarshal(data, &ExperienciaLaboral); err == nil {
		var ExperienciaLaboralPost map[string]interface{}
		InfoComplementariaTercero := ExperienciaLaboral["InfoComplementariaTercero"].([]interface{})[0]
		Experiencia := ExperienciaLaboral["Experiencia"].(map[string]interface{})

		Dato := fmt.Sprintf("%v", InfoComplementariaTercero.(map[string]interface{})["Dato"].(string))
		var dato map[string]interface{}
		json.Unmarshal([]byte(Dato), &dato)
		Dedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Id"].(float64)
		NombreDedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Nombre"].(string)
		Vinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Id"].(float64)
		NombreVinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Nombre"].(string)
		CargoID := Experiencia["Cargo"].(map[string]interface{})["Id"].(float64)
		NombreCargo := Experiencia["Cargo"].(map[string]interface{})["Nombre"].(string)

		// GET para traer el id de experencia_labora info complementaria
		var resultadoInfoComplementaria []map[string]interface{}
		errIdInfo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=GrupoInfoComplementariaId__Id:19,CodigoAbreviacion:EXP_LABORAL,Activo:true&limit=0", &resultadoInfoComplementaria)
		if errIdInfo == nil && fmt.Sprintf("%v", resultadoInfoComplementaria[0]["System"]) != "map[]" {
			if resultadoInfoComplementaria[0]["Status"] != 404 && resultadoInfoComplementaria[0]["Id"] != nil {

				idInfoExperencia = fmt.Sprintf("%v", resultadoInfoComplementaria[0]["Id"])
			} else {
				if resultadoInfoComplementaria[0]["Message"] == "Not found resource" {
					return APIResponseDTO
				} else {
					logs.Error(resultadoInfoComplementaria)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, resultadoInfoComplementaria)
					return APIResponseDTO
				}
			}
		} else {
			logs.Error(errIdInfo)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil ,errIdInfo)
			return APIResponseDTO
		}
		intVar, _ := strconv.Atoi(idInfoExperencia)

		ExperienciaLaboralData := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": intVar},
			"Dato": "{\n    " +
				"\"Nit\": \"" + fmt.Sprintf("%v", dato["NumeroIdentificacion"]) + "\",    " +
				"\"FechaInicio\": \"" + Experiencia["FechaInicio"].(string) + "\",    " +
				"\"FechaFinalizacion\": \"" + Experiencia["FechaFinalizacion"].(string) + "\",    " +
				"\"TipoDedicacion\": { \"Id\": \"" + fmt.Sprintf("%v", Dedicacion) + "\", \"Nombre\": \"" + NombreDedicacion + "\"},    " +
				"\"TipoVinculacion\": { \"Id\": \"" + fmt.Sprintf("%v", Vinculacion) + "\", \"Nombre\": \"" + NombreVinculacion + "\"},    " +
				"\"Cargo\": { \"Id\": \"" + fmt.Sprintf("%v", CargoID) + "\", \"Nombre\": \"" + NombreCargo + "\"},    " +
				"\"Actividades\": \"" + Experiencia["Actividades"].(string) + "\",    " +
				"\"Soporte\": \"" + fmt.Sprintf("%v", Experiencia["DocumentoId"]) + "\"" +
				"\n }",
			"Activo": true,
		}
		//formatdata.JsonPrint(ExperienciaLaboralData)

		errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &ExperienciaLaboralPost, ExperienciaLaboralData)
		if errExperiencia == nil && fmt.Sprintf("%v", ExperienciaLaboralPost["System"]) != "map[]" && ExperienciaLaboralPost["Id"] != nil {
			if ExperienciaLaboralPost["Status"] != 400 {
				respuesta["FormacionAcademica"] = ExperienciaLaboralPost
				//formatdata.JsonPrint(respuesta)
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuesta, nil)
			} else {
				logs.Error(ExperienciaLaboralPost)
				APIResponseDTO = requestresponse.APIResponseDTO(true, 400, nil ,ExperienciaLaboralPost)
				return APIResponseDTO
			}
		} else {
			logs.Error(errExperiencia)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil ,errExperiencia)
			return APIResponseDTO
		}

	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil ,err)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func GetInfoEmpresa(idEmpresa string) (APIResponseDTO requestresponse.APIResponse) {
	var empresa []map[string]interface{}
	var empresaTercero map[string]interface{}
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})

	re := regexp.MustCompile("[^0-9-]")
	idEmpresa = re.ReplaceAllString(idEmpresa, "")

	partes := strings.Split(idEmpresa, "-")
	numeroNit := partes[0]

	endpoit := "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + numeroNit

	//GET que asocia el nit con la empresa
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+endpoit, &empresa)
	if errNit == nil {
		if empresa != nil && len(empresa[0]) > 0 {
			respuesta["NumeroIdentificacion"] = idEmpresa
			idEmpresa := empresa[0]["TerceroId"].(map[string]interface{})["Id"]
			//GET que trae la información de la empresa
			errUniversidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idEmpresa), &empresaTercero)
			if errUniversidad == nil && fmt.Sprintf("%v", empresaTercero["System"]) != "map[]" && empresaTercero["Id"] != nil {
				if empresaTercero["Status"] != "400" {
					respuesta["NombreCompleto"] = map[string]interface{}{
						"Id":             idEmpresa,
						"NombreCompleto": empresaTercero["NombreCompleto"],
					}

					//GET para traer los datos de la ubicación
					var lugar map[string]interface{}
					errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", empresaTercero["LugarOrigen"]), &lugar)
					if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
						if lugar["Status"] != "404" {
							respuesta["Ubicacion"] = map[string]interface{}{
								"Id":     lugar["PAIS"].(map[string]interface{})["Id"],
								"Nombre": lugar["PAIS"].(map[string]interface{})["Nombre"],
							}
						} else {
							logs.Error(lugar["Status"])
							respuesta["Ubicacion"] = nil
							//c.Data["json"] = map[string]interface{}{"Code": "400", "Body": lugar["Status"], "Type": "error"}
							//c.Data["system"] = lugar
							//c.Abort("404")
						}
					} else {
						logs.Error(errLugar)
						respuesta["Ubicacion"] = nil
						//c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errLugar.Error(), "Type": "error"}
						//c.Data["system"] = lugar
						//c.Abort("404")
					}

					//GET para traer la dirección de la empresa (info_complementaria 54)
					var resultadoDireccion []map[string]interface{}
					errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:54,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoDireccion)
					if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
						if resultadoDireccion[0]["Status"] != "404" && resultadoDireccion[0]["Id"] != nil {
							var direccionJSON map[string]interface{}
							if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &direccionJSON); err != nil {
								respuesta["Direccion"] = nil
							} else {
								respuesta["Direccion"] = direccionJSON["address"]
							}
						} else {
							if resultadoDireccion[0]["Message"] == "Not found resource" {
								respuesta["Direccion"] = nil
								//c.Data["json"] = nil
							} else {
								logs.Error(resultadoDireccion)
								respuesta["Direccion"] = nil
								//c.Data["system"] = resultadoDireccion
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(errDireccion)
						respuesta["Direccion"] = nil
						//c.Data["system"] = errDireccion
						//c.Abort("404")
					}

					// GET para traer el telefono de la empresa (info_complementaria 51)
					var resultadoTelefono []map[string]interface{}
					errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:51,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoTelefono)
					if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
						if resultadoTelefono[0]["Status"] != "404" && resultadoTelefono[0]["Id"] != nil {
							var telefonoJSON map[string]interface{}
							if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &telefonoJSON); err != nil {
								respuesta["Telefono"] = nil
							} else {
								respuesta["Telefono"] = telefonoJSON["telefono"]
							}
						} else {
							if resultadoTelefono[0]["Message"] == "Not found resource" {
								respuesta["Telefono"] = nil
								//c.Data["json"] = nil
							} else {
								logs.Error(resultadoTelefono)
								respuesta["Telefono"] = nil
								//c.Data["system"] = resultadoTelefono
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(errTelefono)
						respuesta["Telefono"] = nil
						//c.Data["system"] = errTelefono
						//c.Abort("404")
					}

					// GET para traer el correo de la empresa (info_complementaria 53)
					var resultadoCorreo []map[string]interface{}
					errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:53,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoCorreo)
					if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
						if resultadoCorreo[0]["Status"] != "404" && resultadoCorreo[0]["Id"] != nil {
							var correoJSON map[string]interface{}
							if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJSON); err != nil {
								respuesta["Correo"] = nil
							} else {
								respuesta["Correo"] = correoJSON["email"]
							}
						} else {
							if resultadoCorreo[0]["Message"] == "Not found resource" {
								respuesta["Correo"] = nil
								//c.Data["json"] = nil
							} else {
								logs.Error(resultadoCorreo)
								respuesta["Correo"] = nil
								//c.Data["system"] = resultadoCorreo
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(errCorreo)
						respuesta["Correo"] = nil
						//c.Data["system"] = errCorreo
						//c.Abort("404")
					}

					// GET para traer la organizacion de la empresa (info_complementaria 110)
					var resultadoOrganizacion []map[string]interface{}
					errorganizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero/?limit=1&query=TerceroId__Id:"+fmt.Sprintf("%.f", idEmpresa), &resultadoOrganizacion)
					if errorganizacion == nil && fmt.Sprintf("%v", resultadoOrganizacion[0]["System"]) != "map[]" {
						if resultadoOrganizacion[0]["Status"] != "404" && resultadoOrganizacion[0]["Id"] != nil {
							respuesta["TipoTerceroId"] = map[string]interface{}{
								"Id":     resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Id"],
								"Nombre": resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Nombre"],
							}
						} else {
							if resultadoOrganizacion[0]["Message"] == "Not found resource" {
								respuesta["TipoTerceroId"] = nil
								//c.Data["json"] = nil
							} else {
								logs.Error(resultadoOrganizacion)
								respuesta["TipoTerceroId"] = nil
								//c.Data["system"] = resultadoOrganizacion
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(resultadoOrganizacion)
						respuesta["TipoTerceroId"] = nil
						//c.Data["system"] = resultadoOrganizacion
						//c.Abort("404")
					}

					APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuesta, nil)

				} else {
					logs.Error(empresaTercero["Status"])
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, empresaTercero)
					return APIResponseDTO
				}
			} else {
				logs.Error(errUniversidad)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, empresaTercero)
				return APIResponseDTO
			}
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, empresa)
			return APIResponseDTO
		}
	} else {
		logs.Error(errNit)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil ,empresa)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func GetExperienciaLaboralByPersona(idTercero string) (APIResponseDTO requestresponse.APIResponse) {

	var empresa []map[string]interface{}
	var resultado []map[string]interface{}
	resultado = make([]map[string]interface{}, 0)
	var empresaTercero map[string]interface{}
	var errorGetAll bool

	var Data []map[string]interface{}

	fmt.Println("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero?query=TerceroId__Id:" + fmt.Sprintf("%v", idTercero) + ",InfoComplementariaId__CodigoAbreviacion:EXP_LABORAL,Activo:true&limit=0&sortby=Id&order=asc")
	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__CodigoAbreviacion:EXP_LABORAL,Activo:true&limit=0&sortby=Id&order=asc", &Data)
	if errData == nil {
		if Data != nil && fmt.Sprintf("%v", Data) != "[map[]]" {
			var experiencia map[string]interface{}
			for i := 0; i < len(Data); i++ {
				resultadoAux := make(map[string]interface{})
				if err := json.Unmarshal([]byte(Data[i]["Dato"].(string)), &experiencia); err == nil {
					resultadoAux["Id"] = Data[i]["Id"]
					resultadoAux["Actividades"] = experiencia["Actividades"]
					resultadoAux["Cargo"] = experiencia["Cargo"]
					resultadoAux["Soporte"] = experiencia["Soporte"]
					resultadoAux["TipoVinculacion"] = experiencia["TipoVinculacion"]
					resultadoAux["TipoDedicacion"] = experiencia["TipoDedicacion"]
					resultadoAux["FechaFinalizacion"] = experiencia["FechaFinalizacion"]
					resultadoAux["FechaInicio"] = experiencia["FechaInicio"]
					resultadoAux["Nit"] = experiencia["Nit"]

					if reflect.TypeOf(experiencia["Nit"]).Kind() == reflect.Float64 {
						experiencia["Nit"] = fmt.Sprintf("%.f", experiencia["Nit"])
					}

					var endpoit string
					if strings.Contains(fmt.Sprintf("%v", experiencia["Nit"]), "-") {
						var auxNit = strings.Split(fmt.Sprintf("%v", experiencia["Nit"]), "-")
						endpoit = "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + auxNit[0] + ",DigitoVerificacion:" + auxNit[1]
					} else {
						endpoit = "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + fmt.Sprintf("%v", experiencia["Nit"])
					}

					errDatosIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+endpoit, &empresa)
					if errDatosIdentificacion == nil {
						if empresa != nil && len(empresa[0]) > 0 {
							idEmpresa := empresa[0]["TerceroId"].(map[string]interface{})["Id"]

							//GET que trae la información de la empresa
							errEmpresa := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", idEmpresa), &empresaTercero)
							if errEmpresa == nil && fmt.Sprintf("%v", empresaTercero["System"]) != "map[]" && empresaTercero["Id"] != nil {
								if empresaTercero["Status"] != "400" {
									resultadoAux["NombreEmpresa"] = map[string]interface{}{
										"Id":             idEmpresa,
										"NombreCompleto": empresaTercero["NombreCompleto"],
									}
									var lugar map[string]interface{}
									//GET para traer los datos de la ubicación
									errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", empresaTercero["LugarOrigen"]), &lugar)
									if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
										if lugar["Status"] != "404" {
											resultadoAux["Ubicacion"] = map[string]interface{}{
												"Id":     lugar["PAIS"].(map[string]interface{})["Id"],
												"Nombre": lugar["PAIS"].(map[string]interface{})["Nombre"],
											}

											//GET para traer la dirección de la empresa (info_complementaria 54)
											var resultadoDireccion []map[string]interface{}
											errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:54,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoDireccion)
											if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
												if resultadoDireccion[0]["Status"] != "404" && resultadoDireccion[0]["Id"] != nil {
													var direccionJSON map[string]interface{}
													if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &direccionJSON); err != nil {
														resultadoAux["Direccion"] = nil
													} else {
														resultadoAux["Direccion"] = direccionJSON["address"]
													}
												} else {
													resultadoAux["Direccion"] = nil
												}
											} else {
												errorGetAll = true
												APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errDireccion.Error())
											}

											// GET para traer el telefono de la empresa (info_complementaria 51)
											var resultadoTelefono []map[string]interface{}
											errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:51,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoTelefono)
											if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
												if resultadoTelefono[0]["Status"] != "404" && resultadoTelefono[0]["Id"] != nil {
													var telefonoJSON map[string]interface{}
													if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &telefonoJSON); err != nil {
														resultadoAux["Telefono"] = nil
													} else {
														resultadoAux["Telefono"] = telefonoJSON["telefono"]
													}
												} else {
													resultadoAux["Telefono"] = nil
												}
											} else {
												errorGetAll = true
												APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errTelefono.Error())
											}

											// GET para traer el correo de la empresa (info_complementaria 53)
											var resultadoCorreo []map[string]interface{}
											errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:53,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoCorreo)
											if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
												if resultadoCorreo[0]["Status"] != "404" && resultadoCorreo[0]["Id"] != nil {
													var correoJSON map[string]interface{}
													if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJSON); err != nil {
														resultadoAux["Correo"] = nil
													} else {
														resultadoAux["Correo"] = correoJSON["email"]
													}
												} else {
													resultadoAux["Correo"] = nil
												}
											} else {
												errorGetAll = true
												APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errCorreo.Error())
											}

											// GET para traer la organizacion de la empresa (info_complementaria 110)
											var resultadoOrganizacion []map[string]interface{}
											errorganizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero/?limit=1&query=TerceroId__Id:"+fmt.Sprintf("%.f", idEmpresa), &resultadoOrganizacion)
											if errorganizacion == nil && fmt.Sprintf("%v", resultadoOrganizacion[0]["System"]) != "map[]" {
												if resultadoOrganizacion[0]["Status"] != "404" && resultadoOrganizacion[0]["Id"] != nil {

													resultadoAux["TipoTerceroId"] = map[string]interface{}{
														"Id":     resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Id"],
														"Nombre": resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Nombre"],
													}
												} else {
													resultadoAux["TipoTerceroId"] = nil
												}
											} else {
												errorGetAll = true
												APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errorganizacion.Error())
											}

										} else {
											resultadoAux["Ubicacion"] = nil
											resultadoAux["Direccion"] = nil
											resultadoAux["Telefono"] = nil
											resultadoAux["Correo"] = nil
											resultadoAux["TipoTerceroId"] = nil
										}
									} else {
										errorGetAll = true
										APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errLugar.Error())
									}
								} else {
									resultadoAux["NombreCompleto"] = nil
									resultadoAux["Ubicacion"] = nil
									resultadoAux["Direccion"] = nil
									resultadoAux["Telefono"] = nil
									resultadoAux["Correo"] = nil
									resultadoAux["TipoTerceroId"] = nil
								}
							} else {
								errorGetAll = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errEmpresa.Error())
							}
						} else {
							resultadoAux["NombreEmpresa"] = nil
							resultadoAux["Ubicacion"] = nil
							resultadoAux["Direccion"] = nil
							resultadoAux["Telefono"] = nil
							resultadoAux["Correo"] = nil
							resultadoAux["TipoTerceroId"] = nil
						}
					} else {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errDatosIdentificacion.Error())
					}

					resultado = append(resultado, resultadoAux)
				} else {
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				}
			}
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		return APIResponseDTO
	} else {
		return APIResponseDTO
	}
}

func ActualizarExperienciaLaboral(idTercero string, data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var Data []map[string]interface{}
	var Put map[string]interface{}
	var ExperienciaLaboral interface{}
	var Experiencia map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var errorGetAll bool

	if err := json.Unmarshal(data, &Experiencia); err == nil {
		InfoComplementariaTercero := Experiencia["InfoComplementariaTercero"].([]interface{})[0]
		ExperienciaLaboral = Experiencia["Experiencia"]
		Dato := fmt.Sprintf("%v", InfoComplementariaTercero.(map[string]interface{})["Dato"].(string))
		var dato map[string]interface{}
		json.Unmarshal([]byte(Dato), &dato)
		Dedicacion := ExperienciaLaboral.(map[string]interface{})["TipoDedicacion"].(map[string]interface{})["Id"]
		NombreDedicacion := ExperienciaLaboral.(map[string]interface{})["TipoDedicacion"].(map[string]interface{})["Nombre"].(string)
		Vinculacion := ExperienciaLaboral.(map[string]interface{})["TipoVinculacion"].(map[string]interface{})["Id"]
		NombreVinculacion := ExperienciaLaboral.(map[string]interface{})["TipoVinculacion"].(map[string]interface{})["Nombre"].(string)
		CargoID := ExperienciaLaboral.(map[string]interface{})["Cargo"].(map[string]interface{})["Id"]
		NombreCargo := ExperienciaLaboral.(map[string]interface{})["Cargo"].(map[string]interface{})["Nombre"].(string)

		errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+idTercero, &Data)
		if errData == nil {
			if Data != nil {
				Data[0]["Dato"] = "{\n    " +
					"\"Nit\": " + fmt.Sprintf("%v", dato["NumeroIdentificacion"]) + ",    " +
					"\"FechaInicio\": \"" + ExperienciaLaboral.(map[string]interface{})["FechaInicio"].(string) + "\",    " +
					"\"FechaFinalizacion\": \"" + ExperienciaLaboral.(map[string]interface{})["FechaFinalizacion"].(string) + "\",    " +
					"\"TipoDedicacion\": { \"Id\": \"" + fmt.Sprintf("%v", Dedicacion) + "\", \"Nombre\": \"" + NombreDedicacion + "\"},    " +
					"\"TipoVinculacion\": { \"Id\": \"" + fmt.Sprintf("%v", Vinculacion) + "\", \"Nombre\": \"" + NombreVinculacion + "\"},    " +
					"\"Cargo\": { \"Id\": \"" + fmt.Sprintf("%v", CargoID) + "\", \"Nombre\": \"" + NombreCargo + "\"},    " +
					"\"Actividades\": \"" + ExperienciaLaboral.(map[string]interface{})["Actividades"].(string) + "\",    " +
					"\"Soporte\": \"" + fmt.Sprintf("%v", ExperienciaLaboral.(map[string]interface{})["DocumentoId"]) + "\"" +
					"\n }"
			}
		}

		errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idTercero, "PUT", &Put, Data[0])
		if errPut == nil {
			if Put != nil {
				resultado = Put
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
			}
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPut.Error())

		}

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil ,ExperienciaLaboral)
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		return APIResponseDTO
	} else {
		return APIResponseDTO
	}
}

func GetExperienciaLaboralById(idExperiencia string) (APIResponseDTO requestresponse.APIResponse) {
	var resultado map[string]interface{}
	//resultado experiencia
	var experiencia []map[string]interface{}

	errExperiencia := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral/?query=Id:"+idExperiencia, &experiencia)
	if errExperiencia == nil && fmt.Sprintf("%v", experiencia[0]["System"]) != "map[]" {
		if experiencia[0]["Status"] != 404 {
			//buscar soporte_experiencia_laboral
			var soporte []map[string]interface{}

			errSoporte := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/?query=ExperienciaLaboral:"+idExperiencia+"&fields=Documento", &soporte)
			if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
				if soporte[0]["Status"] != 404 {
					experiencia[0]["Documento"] = soporte[0]["Documento"]

				} else {
					if soporte[0]["Message"] == "Not found resource" {
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
						return APIResponseDTO
					} else {
						logs.Error(soporte)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404,nil ,errSoporte)
						return APIResponseDTO
					}
				}
			} else {
				logs.Error(soporte)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil ,errSoporte)
				return APIResponseDTO
			}

			//buscar organizacion_experiencia_laboral
			var organizacion []map[string]interface{}
			errOrganizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=1&query=Id:"+
				// fmt.Sprintf("%v", experiencia[u]["Id"])+"&fields=Documento", &soporte)
				fmt.Sprintf("%v", experiencia[0]["Organizacion"]), &organizacion)
			if errOrganizacion == nil && fmt.Sprintf("%v", organizacion[0]["System"]) != "map[]" {
				if organizacion[0]["Status"] != 404 && organizacion[0]["Id"] != nil {

					// unmarshall dato
					var organizacionJson map[string]interface{}
					if err := json.Unmarshal([]byte(organizacion[0]["Dato"].(string)), &organizacionJson); err != nil {
						experiencia[0]["Organizacion"] = nil
					} else {
						experiencia[0]["Organizacion"] = organizacionJson
					}

				} else {
					if organizacion[0]["Message"] == "Not found resource" {
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
						return APIResponseDTO
					} else {
						logs.Error(organizacion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errOrganizacion)
						return APIResponseDTO
					}
				}
			} else {
				logs.Error(organizacion)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errOrganizacion)
				return APIResponseDTO
			}

			resultado = experiencia[0]
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)

		} else {
			if experiencia[0]["Message"] == "Not found resource" {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				return APIResponseDTO
			} else {
				logs.Error(experiencia)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil ,errExperiencia)
				return APIResponseDTO
			}
		}
	} else {
		logs.Error(experiencia)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil ,errExperiencia)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func DeleteExperienciaById(idExperiencia string) (APIResponseDTO requestresponse.APIResponse) {
	var Data []map[string]interface{}
	var Put map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var errorGetAll bool

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+idExperiencia, &Data)
	if errData == nil {
		if Data != nil {
			Data[0]["Activo"] = false

			errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idExperiencia, "PUT", &Put, Data[0])
			if errPut == nil {
				if Put != nil {
					resultado = Put
				} else {
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				}
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPut.Error())
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errData.Error())
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		return APIResponseDTO
	}else {
		return APIResponseDTO
	}

}