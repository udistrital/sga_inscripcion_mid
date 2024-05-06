package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func CrearFormacion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var FormacionAcademica map[string]interface{}
	var idInfoFormacion string
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})

	if err := json.Unmarshal(data, &FormacionAcademica); err == nil {
		var FormacionAcademicaPost map[string]interface{}

		NombrePrograma := fmt.Sprintf("%v", FormacionAcademica["ProgramaAcademicoId"])
		FechaI := fmt.Sprintf("%q", FormacionAcademica["FechaInicio"])
		FechaF := fmt.Sprintf("%q", FormacionAcademica["FechaFinalizacion"])
		TituloTG := fmt.Sprintf("%q", FormacionAcademica["TituloTrabajoGrado"])
		DescripcionTG := fmt.Sprintf("%q", FormacionAcademica["DescripcionTrabajoGrado"])
		DocumentoId := fmt.Sprintf("%v", FormacionAcademica["DocumentoId"])
		NitU := fmt.Sprintf("%q", FormacionAcademica["NitUniversidad"])
		// NivelFormacion := fmt.Sprintf("%v", FormacionAcademica["NivelFormacion"])

		// GET para traer el id de experencia_labora info complementaria
		var resultadoInfoComplementaria []map[string]interface{}
		errIdInfo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=GrupoInfoComplementariaId__Id:18,CodigoAbreviacion:FORM_ACADEMICA,Activo:true&limit=0", &resultadoInfoComplementaria)
		if errIdInfo == nil && fmt.Sprintf("%v", resultadoInfoComplementaria[0]["System"]) != "map[]" {
			if resultadoInfoComplementaria[0]["Status"] != 404 && resultadoInfoComplementaria[0]["Id"] != nil {

				idInfoFormacion = fmt.Sprintf("%v", resultadoInfoComplementaria[0]["Id"])
			} else {
				if resultadoInfoComplementaria[0]["Message"] == "Not found resource" {
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
					return APIResponseDTO
				} else {
					logs.Error(resultadoInfoComplementaria)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, resultadoInfoComplementaria)
					return APIResponseDTO
				}
			}
		} else {
			logs.Error(errIdInfo)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errIdInfo)
			return APIResponseDTO
		}
		intVar, _ := strconv.Atoi(idInfoFormacion)

		FormacionAcademicaData := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": intVar},
			"Dato": "{\n    " +
				"\"ProgramaAcademico\": " + NombrePrograma + ",    " +
				"\"FechaInicio\": " + FechaI + ",    " +
				"\"FechaFin\": " + FechaF + ",    " +
				"\"TituloTrabajoGrado\": " + TituloTG + ",    " +
				"\"DesTrabajoGrado\": " + DescripcionTG + ",    " +
				"\"DocumentoId\": " + DocumentoId + ",    " +
				"\"NitUniversidad\": " + NitU +
				// "\"NivelFormacion\": " + NivelFormacion + ", \n " +
				"\n }",
			"Activo": true,
		}

		errFormacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &FormacionAcademicaPost, FormacionAcademicaData)
		if errFormacion == nil && fmt.Sprintf("%v", FormacionAcademicaPost["System"]) != "map[]" && FormacionAcademicaPost["Id"] != nil {
			if FormacionAcademicaPost["Status"] != 400 {
				respuesta["FormacionAcademica"] = FormacionAcademicaPost
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuesta, nil)
			} else {
				logs.Error(errFormacion)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, FormacionAcademicaPost)
				return APIResponseDTO
			}
		} else {
			logs.Error(errFormacion)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errFormacion)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func GetUniversidadInfo(idUniversidad string) (APIResponseDTO requestresponse.APIResponse) {
	var universidad []map[string]interface{}
	var universidadTercero map[string]interface{}
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})

	re := regexp.MustCompile("[^0-9-]")
	idUniversidad = re.ReplaceAllString(idUniversidad, "")
	partes := strings.Split(idUniversidad, "-")
	numeroNit := partes[0]

	endpoit := "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + numeroNit

	//GET que asocia el nit con la universidad
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+endpoit, &universidad)
	if errNit == nil {
		if universidad != nil && fmt.Sprintf("%v", universidad[0]) != "map[]" {
			respuesta["NumeroIdentificacion"] = idUniversidad
			idUniversidad := universidad[0]["TerceroId"].(map[string]interface{})["Id"]
			//GET que trae la información de la universidad
			errUniversidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idUniversidad), &universidadTercero)
			if errUniversidad == nil && fmt.Sprintf("%v", universidadTercero["System"]) != "map[]" && universidadTercero["Id"] != nil {
				if universidadTercero["Status"] != "400" {
					respuesta["NombreCompleto"] = map[string]interface{}{
						"Id":             idUniversidad,
						"NombreCompleto": universidadTercero["NombreCompleto"],
					}

					var lugar map[string]interface{}
					//GET para traer los datos de la ubicación
					errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", universidadTercero["LugarOrigen"]), &lugar)
					if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
						if lugar["Status"] != "404" {
							respuesta["Ubicacion"] = map[string]interface{}{
								"Id":     lugar["PAIS"].(map[string]interface{})["Id"],
								"Nombre": lugar["PAIS"].(map[string]interface{})["Nombre"],
							}
						} else {
							logs.Error(errLugar)
							respuesta["Ubicacion"] = nil
							//c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errLugar.Error(), "Type": "error"}
							//c.Data["system"] = lugar
							//c.Abort("400")
						}
					} else {
						logs.Error(errLugar)
						respuesta["Ubicacion"] = nil
						//c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errLugar.Error(), "Type": "error"}
						//c.Data["system"] = lugar
						//c.Abort("400")
					}

					//GET para traer la dirección de la universidad (info_complementaria 54)
					var resultadoDireccion []map[string]interface{}
					errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:54,TerceroId:"+fmt.Sprintf("%.f", idUniversidad), &resultadoDireccion)
					if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
						if resultadoDireccion[0]["Status"] != "404" && resultadoDireccion[0]["Id"] != nil {
							// Unmarshall dato
							var direccionJson map[string]interface{}
							if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &direccionJson); err != nil {
								respuesta["Direccion"] = nil
							} else {
								respuesta["Direccion"] = direccionJson["address"]
							}
						} else {
							if resultadoDireccion[0]["Message"] == "Not found resource" {
								respuesta["Direccion"] = nil
								//c.Data["json"] = nil
							} else {
								logs.Error(resultadoDireccion)
								respuesta["Direccion"] = nil
								//c.Data["system"] = errDireccion
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(resultadoDireccion)
						respuesta["Direccion"] = nil
						//c.Data["system"] = resultadoDireccion
						//c.Abort("404")
					}

					// GET para traer el telefono de la universidad (info_complementaria 51)
					var resultadoTelefono []map[string]interface{}
					errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:51,TerceroId:"+fmt.Sprintf("%.f", idUniversidad), &resultadoTelefono)
					if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
						if resultadoTelefono[0]["Status"] != "404" && resultadoTelefono[0]["Id"] != nil {
							// Unmarshall dato
							var telefonoJson map[string]interface{}
							if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &telefonoJson); err != nil {
								respuesta["Telefono"] = nil
							} else {
								respuesta["Telefono"] = telefonoJson["telefono"]
							}
						} else {
							if resultadoTelefono[0]["Message"] == "Not found resource" {
								respuesta["Telefono"] = nil
								//c.Data["json"] = nil
							} else {
								logs.Error(resultadoTelefono)
								respuesta["Telefono"] = nil
								//c.Data["system"] = errTelefono
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(resultadoTelefono)
						respuesta["Telefono"] = nil
						//c.Data["system"] = resultadoTelefono
						//c.Abort("404")
					}

					// GET para traer el correo de la universidad (info_complementaria 53)
					var resultadoCorreo []map[string]interface{}
					errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:53,TerceroId:"+fmt.Sprintf("%.f", idUniversidad), &resultadoCorreo)
					if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
						if resultadoCorreo[0]["Status"] != "404" && resultadoCorreo[0]["Id"] != nil {
							// Unmarshall dato
							var correoJson map[string]interface{}
							if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJson); err != nil {
								respuesta["Correo"] = nil
							} else {
								respuesta["Correo"] = correoJson["email"]
							}
						} else {
							if resultadoCorreo[0]["Message"] == "Not found resource" {
								respuesta["Correo"] = nil
								//c.Data["json"] = nil
							} else {
								logs.Error(resultadoCorreo)
								respuesta["Correo"] = nil
								//c.Data["system"] = errCorreo
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(resultadoCorreo)
						respuesta["Correo"] = nil
						//c.Data["system"] = resultadoTelefono
						//c.Abort("404")
					}

					APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuesta, nil)

				} else {
					logs.Error(errUniversidad)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, universidadTercero, errUniversidad.Error())
					return APIResponseDTO
				}
			} else {
				logs.Error(errUniversidad)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, universidadTercero, errUniversidad.Error())
				return APIResponseDTO
			}
		} else {
			logs.Error(errNit)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, universidad, errNit.Error())
			return APIResponseDTO
		}
	} else {
		logs.Error(errNit)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, universidad, errNit.Error())
		return APIResponseDTO
	}

	return APIResponseDTO
}

func GetUniversidadNombre(nombre string) (APIResponseDTO requestresponse.APIResponse) {
	var universidades []map[string]interface{}
	NombresAux := strings.Split(nombre, " ")

	if len(NombresAux) == 1 {
		err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=NombreCompleto__contains:"+nombre+"&limit=0", &universidades)
		if err == nil {
			if universidades != nil && fmt.Sprintf("%v", universidades[0]) != "map[]" {
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, universidades, nil)
			} else {
				logs.Error(universidades)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, err)
				return APIResponseDTO
			}
		} else {
			logs.Error(universidades)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, err)
			return APIResponseDTO
		}
	} else if len(NombresAux) > 1 {
		err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=NombreCompleto__contains:"+NombresAux[0]+",NombreCompleto__contains:"+NombresAux[1]+"&limit=0", &universidades)
		if err == nil {
			if universidades != nil {
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, universidades, nil)
			} else {
				logs.Error(universidades)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, err)
				return APIResponseDTO
			}
		} else {
			logs.Error(universidades)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, err)
			return APIResponseDTO
		}
	}
	return APIResponseDTO
}

func ActualizarFormacionAcademica(idFormacion string, data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var Data []map[string]interface{}
	var Put map[string]interface{}
	var InfoAcademica map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var errorGetAll bool

	if err := json.Unmarshal(data, &InfoAcademica); err == nil {
		errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+idFormacion, &Data)
		if errData == nil {
			if Data != nil {
				Data[0]["Dato"] = "{\n    " +
					"\"ProgramaAcademico\": " + fmt.Sprintf("%v", InfoAcademica["ProgramaAcademicoId"]) + ",    " +
					"\"FechaInicio\": " + fmt.Sprintf("%q", InfoAcademica["FechaInicio"]) + ",    " +
					"\"FechaFin\": " + fmt.Sprintf("%q", InfoAcademica["FechaFinalizacion"]) + ",    " +
					"\"TituloTrabajoGrado\": " + fmt.Sprintf("%q", InfoAcademica["TituloTrabajoGrado"]) + ",    " +
					"\"DesTrabajoGrado\": " + fmt.Sprintf("%q", InfoAcademica["DescripcionTrabajoGrado"]) + ",    " +
					"\"DocumentoId\": " + fmt.Sprintf("%v", InfoAcademica["DocumentoId"]) + ",    " +
					"\"NitUniversidad\": " + fmt.Sprintf("%q", InfoAcademica["NitUniversidad"]) +
					"\n }"

				errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idFormacion, "PUT", &Put, Data[0])
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
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		return APIResponseDTO
	} else {
		return APIResponseDTO
	}
}

func GetFormacionAcademicaById(id string) (APIResponseDTO requestresponse.APIResponse) {
	var Data []map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var errorGetAll bool

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+id, &Data)
	if errData == nil {
		if Data != nil {
			var formacion map[string]interface{}
			resultadoAux := make(map[string]interface{})
			if err := json.Unmarshal([]byte(Data[0]["Dato"].(string)), &formacion); err == nil {
				resultadoAux["Id"] = Data[0]["Id"]
				resultadoAux["Nit"] = formacion["NitUniversidad"]
				resultadoAux["Documento"] = formacion["DocumentoId"]
				resultadoAux["DescripcionTrabajoGrado"] = formacion["DesTrabajoGrado"]
				resultadoAux["FechaInicio"] = formacion["FechaInicio"]
				resultadoAux["FechaFinalizacion"] = formacion["FechaFin"]
				resultadoAux["TituloTrabajoGrado"] = formacion["TituloTrabajoGrado"]
				NumProyecto := fmt.Sprintf("%v", formacion["ProgramaAcademico"])
				//GET para consultar el proyecto curricular

				var ProyectoV2 map[string]interface{}
				errProyectoV2 := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro?query=TipoParametroId__Id:60,Id:"+fmt.Sprintf("%v", NumProyecto)+"&limit=0", &ProyectoV2)
				if errProyectoV2 == nil && ProyectoV2["Status"] == "200" && fmt.Sprintf("%v", ProyectoV2["Data"]) != "[map[]]" {
					resultadoAux["ProgramaAcademico"] = map[string]interface{}{
						"Id":     NumProyecto,
						"Nombre": ProyectoV2["Data"].([]interface{})[0].(map[string]interface{})["Nombre"],
					}
				} else {
					var Proyecto []map[string]interface{}
					errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", NumProyecto)+"&limit=0", &Proyecto)
					if errProyecto == nil && fmt.Sprintf("%v", Proyecto[0]) != "map[]" && Proyecto[0]["Id"] != nil {
						if Proyecto[0]["Status"] != 404 {
							resultadoAux["ProgramaAcademico"] = Proyecto[0]
						} else {
							errorGetAll = true
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errProyecto.Error())
						}
					} else {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
					}
				}

				resultado = resultadoAux
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
			}
		}
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		return APIResponseDTO
	} else {
		return APIResponseDTO
	}
}

func GetFormacionAcademicaByIdTercero(idTercero string) (APIResponseDTO requestresponse.APIResponse) {
	var resultado []map[string]interface{}
	resultado = make([]map[string]interface{}, 0)
	var Data []map[string]interface{}
	var errorGetAll bool

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idTercero+",InfoComplementariaId__CodigoAbreviacion:FORM_ACADEMICA,Activo:true&limit=0&sortby=Id&order=asc", &Data)
	if errData == nil {
		if Data != nil && fmt.Sprintf("%v", Data) != "[map[]]" {
			var formacion map[string]interface{}
			for i := 0; i < len(Data); i++ {
				resultadoAux := make(map[string]interface{})
				if err := json.Unmarshal([]byte(Data[i]["Dato"].(string)), &formacion); err == nil {
					if formacion["ProgramaAcademico"] != "colegio" {
						resultadoAux["Id"] = Data[i]["Id"]
						resultadoAux["Nit"] = formacion["NitUniversidad"]
						resultadoAux["Documento"] = formacion["DocumentoId"]
						resultadoAux["FechaInicio"] = formacion["FechaInicio"]
						resultadoAux["FechaFinalizacion"] = formacion["FechaFin"]

						endpoit := "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + fmt.Sprintf("%v", formacion["NitUniversidad"])

						if strings.Contains(fmt.Sprintf("%v", formacion["NitUniversidad"]), "-") {
							var auxId = strings.Split(fmt.Sprintf("%v", formacion["NitUniversidad"]), "-")
							endpoit = "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + auxId[0] + ",DigitoVerificacion:" + auxId[1]
						}

						//GET para obtener el ID que relaciona las tablas tipo_documento y tercero
						var IdTercero []map[string]interface{}
						errIdTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+endpoit, &IdTercero)
						if errIdTercero == nil && fmt.Sprintf("%v", IdTercero[0]) != "map[]" && IdTercero[0]["Id"] != nil {
							if IdTercero[0]["Status"] != "404" {
								IdTerceroAux := IdTercero[0]["TerceroId"].(map[string]interface{})["Id"]

								// GET para traer el nombre de la universidad y el país
								var Tercero []map[string]interface{}
								errTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+fmt.Sprintf("%v", IdTerceroAux), &Tercero)
								if errTercero == nil && fmt.Sprintf("%v", Tercero[0]) != "map[]" && Tercero[0]["Id"] != nil {
									if Tercero[0]["Status"] != "404" {
										resultadoAux["NombreCompleto"] = Tercero[0]["NombreCompleto"]
										var lugar map[string]interface{}

										//GET para traer los datos de la ubicación
										errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", Tercero[0]["LugarOrigen"]), &lugar)
										if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
											if lugar["Status"] != "404" {
												resultadoAux["Ubicacion"] = lugar["PAIS"].(map[string]interface{})["Nombre"]
											} else {
												resultadoAux["Ubicacion"] = nil
												/* errorGetAll = true
												alertas = append(alertas, errLugar.Error())
												alerta.Code = "400"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta} */
											}
										} else {
											resultadoAux["Ubicacion"] = nil
											/* errorGetAll = true
											alertas = append(alertas, "No data found")
											alerta.Code = "404"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Response": alerta} */
										}
									} else {
										errorGetAll = true
										APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errTercero.Error())
									}
								} else {
									errorGetAll = true
									APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
								}
							} else {
								errorGetAll = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errIdTercero.Error())
							}
						} else {
							errorGetAll = true
							APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
						}

						NumProyecto := fmt.Sprintf("%v", formacion["ProgramaAcademico"])

						var ProyectoV2 map[string]interface{}
						errProyectoV2 := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro?query=TipoParametroId__Id:60,Id:"+fmt.Sprintf("%v", NumProyecto)+"&limit=0", &ProyectoV2)
						if errProyectoV2 == nil && ProyectoV2["Status"] == "200" && fmt.Sprintf("%v", ProyectoV2["Data"]) != "[map[]]" {
							resultadoAux["ProgramaAcademico"] = map[string]interface{}{
								"Id":     NumProyecto,
								"Nombre": ProyectoV2["Data"].([]interface{})[0].(map[string]interface{})["Nombre"],
							}
						} else {
							//GET para consultar el proyecto curricular Modo antiguo
							var Proyecto []map[string]interface{}
							errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", NumProyecto)+"&limit=0", &Proyecto)
							if errProyecto == nil && fmt.Sprintf("%v", Proyecto[0]) != "map[]" && Proyecto[0]["Id"] != nil {
								if Proyecto[0]["Status"] != "404" {
									resultadoAux["ProgramaAcademico"] = map[string]interface{}{
										"Id":     NumProyecto,
										"Nombre": Proyecto[0]["Nombre"],
									}
								} else {
									errorGetAll = true
									APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errProyecto.Error())
								}
							} else {
								errorGetAll = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
							}
						}

						resultado = append(resultado, resultadoAux)
					}
				} else {
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				}
			}
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, "No hay formación académica registrada")
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

func NuevoTercero(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroPost map[string]interface{}

	if err := json.Unmarshal(data, &tercero); err == nil {
		//beego.Info(tercero)
		TipoContribuyenteId := map[string]interface{}{
			"Id": 2,
		}
		guardarpersona := map[string]interface{}{
			"NombreCompleto":      tercero["NombreCompleto"],
			"Activo":              false,
			"LugarOrigen":         tercero["Pais"].(map[string]interface{})["Id"].(float64),
			"TipoContribuyenteId": TipoContribuyenteId, // Persona natural actualmente tiene ese id en el api
		}
		errPersona := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &terceroPost, guardarpersona)

		if errPersona == nil && fmt.Sprintf("%v", terceroPost) != "map[]" && terceroPost["Id"] != nil {
			if terceroPost["Status"] != 400 {
				beego.Info("tercero", terceroPost)
				idTerceroCreado := terceroPost["Id"]
				var identificacion map[string]interface{}

				TipoDocumentoId := map[string]interface{}{
					"Id": 7,
				}
				TerceroId := map[string]interface{}{
					"Id": idTerceroCreado,
				}
				TipoTerceroId := map[string]interface{}{
					"Id": tercero["TipoTrecero"].(map[string]interface{})["Id"].(float64),
				}
				identificaciontercero := map[string]interface{}{
					"Numero":             tercero["Nit"],
					"DigitoVerificacion": tercero["Verificacion"],
					"TipoDocumentoId":    TipoDocumentoId,
					"TerceroId":          TerceroId,
					"Activo":             true,
				}
				errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion", "POST", &identificacion, identificaciontercero)
				if errIdentificacion == nil && fmt.Sprintf("%v", identificacion) != "map[]" && identificacion["Id"] != nil {
					if identificacion["Status"] != 400 {
						//beego.Info(identificacion)
						estado := identificacion
						APIResponseDTO = requestresponse.APIResponseDTO(true, 200, estado, nil)

						var telefono map[string]interface{}
						var correo map[string]interface{}
						var direccion map[string]interface{}

						terceroTipoTercero := map[string]interface{}{
							"TerceroId":     TerceroId,
							"TipoTerceroId": TipoTerceroId,
						}

						errTipoTercero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &terceroTipoTercero, terceroTipoTercero)
						if errTipoTercero == nil && fmt.Sprintf("%v", terceroTipoTercero) != "map[]" && terceroTipoTercero["Id"] != nil {
							if terceroTipoTercero["Status"] != 400 {
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["TipoTerceroId"] = terceroTipoTercero["Id"]
								APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)

							} else {
								logs.Error(errTipoTercero)
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errTipoTercero)
								return APIResponseDTO
							}
						} else {
							logs.Error(errTipoTercero)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errTipoTercero)
							return APIResponseDTO
						}

						InfoComplementariaTelefono := map[string]interface{}{
							"Id": 51,
						}
						InfoComplementariaCorreo := map[string]interface{}{
							"Id": 53,
						}
						InfoComplementariaDireccion := map[string]interface{}{
							"Id": 54,
						}

						Telefono := map[string]interface{}{
							"telefono": tercero["Telefono"],
						}
						jsonTelefono, _ := json.Marshal(Telefono)

						Correo := map[string]interface{}{
							"email": tercero["Correo"],
						}
						jsonCorreo, _ := json.Marshal(Correo)

						Direccion := map[string]interface{}{
							"address": tercero["Direccion"],
						}
						jsonDireccion, _ := json.Marshal(Direccion)

						telefonoTercero := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaTelefono,
							"Activo":               true,
							"Dato":                 string(jsonTelefono),
						}
						correoTercero := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaCorreo,
							"Activo":               true,
							"Dato":                 string(jsonCorreo),
						}
						direccionTercero := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaDireccion,
							"Activo":               true,
							"Dato":                 string(jsonDireccion),
						}

						errGenero1 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &telefono, telefonoTercero)
						if errGenero1 == nil && fmt.Sprintf("%v", telefono) != "map[]" && telefono["Id"] != nil {
							//beego.Info(telefono)
							if telefono["Status"] != 400 {
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["Telefono"] = telefono["Id"]
								APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)

							} else {
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errGenero1)
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errGenero1)
								return APIResponseDTO
							}
						} else {
							logs.Error(errGenero1)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errGenero1)
							return APIResponseDTO
						}
						errGenero2 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &correo, correoTercero)
						//beego.Info("correo tercero", correo)
						if errGenero2 == nil && errGenero1 == nil && fmt.Sprintf("%v", correo) != "map[]" && correo["Id"] != nil {
							if correo["Status"] != 400 {
								//beego.Info(correo)
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["Correo"] = correo["Id"]
								APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)

							} else {
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errGenero2)
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errGenero2)
								return APIResponseDTO
							}
						} else {
							//beego.Info("error genero", errGenero2)
							logs.Error(errGenero2)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errGenero2)
							return APIResponseDTO
						}
						errGenero3 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &direccion, direccionTercero)
						if errGenero3 == nil && errGenero2 == nil && errGenero1 == nil && fmt.Sprintf("%v", direccion) != "map[]" && direccion["Id"] != nil {
							if direccion["Status"] != 400 {
								//beego.Info(direccion)
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["Direccion"] = direccion["Id"]
								APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)

							} else {
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errGenero3)
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errGenero3)
								return APIResponseDTO
							}
						} else {
							logs.Error(errGenero3)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errGenero3)
							return APIResponseDTO
						}
					} else {
						//Si pasa un error borra todo lo creado al momento del registro del documento de identidad
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error(errIdentificacion)
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errIdentificacion)
						return APIResponseDTO
					}
				} else {
					//beego.Info("error identificacion", errPersona)
					logs.Error(errIdentificacion)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errIdentificacion)
					return APIResponseDTO
				}
			} else {
				//beego.Info(errPersona)
				logs.Error(errPersona)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPersona)
				return APIResponseDTO
			}
		} else {
			logs.Error(errPersona)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPersona)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}

	return APIResponseDTO

}

func EliminarFormacion(idFormacion string) (APIResponseDTO requestresponse.APIResponse) {
	var Data []map[string]interface{}
	var Put map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var errorGetAll bool

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+idFormacion, &Data)
	if errData == nil {
		if Data != nil {
			Data[0]["Activo"] = false

			errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idFormacion, "PUT", &Put, Data[0])
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
	} else {
		return APIResponseDTO
	}
}
