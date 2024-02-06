package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid_inscripcion/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/time_bogota"
)

func EstadoInscripcion(idPersona string, idPeriodo string) (APIResponseDTO requestresponse.APIResponse) {

	recibosResultado, err := helpers.VerificarRecibos(idPersona, idPeriodo)

	if err == "" {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, recibosResultado)
	} else if err == "400" {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400,nil , "Bad request")
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
	}

	return APIResponseDTO
}

func InformacionFamiliar(data []byte) (APIResponseDTO requestresponse.APIResponse){
	var InformacionFamiliar map[string]interface{}
	var TerceroFamiliarPost map[string]interface{}
	var FamiliarParentescoPost map[string]interface{}
	var InfoContactoPost map[string]interface{}

	if err := json.Unmarshal( data, &InformacionFamiliar); err == nil {
		InfoFamiliarAux := InformacionFamiliar["Familiares"].([]interface{})
		//InfoTercero := InformacionFamiliar["Tercero_Familiar"]

		for _, terceroAux := range InfoFamiliarAux {
			//Se añade primero el familiar a la tabla de terceros
			//fmt.Println(terceroAux)
			TerceroFamiliarAux := terceroAux.(map[string]interface{})["Familiar"].(map[string]interface{})["TerceroFamiliarId"]

			TerceroFamiliar := map[string]interface{}{
				"NombreCompleto":      TerceroFamiliarAux.(map[string]interface{})["NombreCompleto"],
				"Activo":              true,
				"TipoContribuyenteId": map[string]interface{}{"Id": TerceroFamiliarAux.(map[string]interface{})["TipoContribuyenteId"].(map[string]interface{})["Id"].(float64)},
			}
			errTerceroFamiliar := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &TerceroFamiliarPost, TerceroFamiliar)

			if errTerceroFamiliar == nil && fmt.Sprintf("%v", TerceroFamiliarPost) != "map[]" && TerceroFamiliarPost["Id"] != nil {
				if TerceroFamiliarPost["Status"] != 400 {
					// Se relaciona el tercero creado con el aspirante en la tabla tercero_familiar
					FamiliarParentesco := map[string]interface{}{
						"TerceroId":         map[string]interface{}{"Id": terceroAux.(map[string]interface{})["Familiar"].(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"].(float64)},
						"TerceroFamiliarId": map[string]interface{}{"Id": TerceroFamiliarPost["Id"]},
						"TipoParentescoId":  map[string]interface{}{"Id": terceroAux.(map[string]interface{})["Familiar"].(map[string]interface{})["TipoParentescoId"].(map[string]interface{})["Id"].(float64)},
						"Activo":            true,
					}
					errFamiliarParentesco := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar", "POST", &FamiliarParentescoPost, FamiliarParentesco)
					if errFamiliarParentesco == nil && fmt.Sprintf("%v", FamiliarParentescoPost) != "map[]" && FamiliarParentescoPost["Id"] != nil {
						if FamiliarParentescoPost["Status"] != 400 {
							//Se guarda la información del familiar en info_complementaria_tercero
							InfoComplementariaFamiliar := terceroAux.(map[string]interface{})["InformacionContacto"].([]interface{})
							for _, infoComplementaria := range InfoComplementariaFamiliar {
								infoContacto := map[string]interface{}{
									"TerceroId":            map[string]interface{}{"Id": TerceroFamiliarPost["Id"]},
									"InfoComplementariaId": map[string]interface{}{"Id": infoComplementaria.(map[string]interface{})["InfoComplementariaId"].(map[string]interface{})["Id"].(float64)},
									"Dato":                 infoComplementaria.(map[string]interface{})["Dato"],
									"Activo":               true,
								}
								errInfoContacto := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &InfoContactoPost, infoContacto)
								if errInfoContacto == nil && fmt.Sprintf("%v", InfoContactoPost) != "map[]" && InfoContactoPost["Id"] != nil {
									if InfoContactoPost["Status"] != 400 {
										APIResponseDTO = requestresponse.APIResponseDTO(true, 200, TerceroFamiliarPost)
									} else {
										logs.Error(errFamiliarParentesco)
										APIResponseDTO = requestresponse.APIResponseDTO(false, 400, TerceroFamiliarPost)
										return APIResponseDTO
									}
								} else {
									//var resultado2 map[string]interface{}
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
									helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/%.f", FamiliarParentescoPost["Id"]), "DELETE", &resultado2, nil)
									helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/%.f", FamiliarParentescoPost["Id"]))
									logs.Error(errFamiliarParentesco)
									APIResponseDTO = requestresponse.APIResponseDTO(false, 400, TerceroFamiliarPost)
									return APIResponseDTO
								}
							}
						} else {
							//var resultado2 map[string]interface{}
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
							helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
							logs.Error(errFamiliarParentesco)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, TerceroFamiliarPost)
							return APIResponseDTO
						}
					} else {
						//var resultado2 map[string]interface{}
						//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
						helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
						logs.Error(errFamiliarParentesco)
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, TerceroFamiliarPost)
						return APIResponseDTO
					}

				} else {
					//var resultado2 map[string]interface{}
					//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
					helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
					logs.Error(errTerceroFamiliar)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, TerceroFamiliarPost)
					return APIResponseDTO
				}
			} else {
				logs.Error(errTerceroFamiliar)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, TerceroFamiliarPost)
				return APIResponseDTO
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil,  err.Error())
	}
	return APIResponseDTO
}

func Reintegro(data []byte) (APIResponseDTO requestresponse.APIResponse){
	var Reintegro map[string]interface{}
	if err := json.Unmarshal(data, &Reintegro); err == nil {

		var resultadoReintegro map[string]interface{}
		errReintegro := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"tr_inscripcion/reintegro", "POST", &resultadoReintegro, Reintegro)
		if resultadoReintegro["Type"] == "error" || errReintegro != nil || resultadoReintegro["Status"] == "404" || resultadoReintegro["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoReintegro)
			return APIResponseDTO
		} else {
			fmt.Println("Reintegrro registrado")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, Reintegro)
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func TransferenciaPost(data []byte)(APIResponseDTO requestresponse.APIResponse){
	var Transferencia map[string]interface{}

	if err := json.Unmarshal(data, &Transferencia); err == nil {

		var resultadoTransferencia map[string]interface{}
		errTransferencia := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"tr_inscripcion/transferencia", "POST", &resultadoTransferencia, Transferencia)
		if resultadoTransferencia["Type"] == "error" || errTransferencia != nil || resultadoTransferencia["Status"] == "404" || resultadoTransferencia["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoTransferencia)
			return APIResponseDTO
		} else {
			fmt.Println("Transferencia registrada")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, Transferencia)
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func InfoIcfesColegio(data []byte) (APIResponseDTO requestresponse.APIResponse){
	var InfoIcfesColegio map[string]interface{}

	if err := json.Unmarshal(data, &InfoIcfesColegio); err == nil {

		var InscripcionPregrado = InfoIcfesColegio["InscripcionPregrado"].(map[string]interface{})
		var InfoComplementariaTercero = InfoIcfesColegio["InfoComplementariaTercero"].([]interface{})
		var InformacionColegio = InfoIcfesColegio["dataColegio"].(map[string]interface{})
		var Tercero = InfoIcfesColegio["Tercero"].(map[string]interface{})
		var date = time.Now()

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var dato = datoInfoComplementaria.(map[string]interface{})
			dato["FechaCreacion"] = date
			dato["FechaModificacion"] = date
			var resultadoInfoComeplementaria map[string]interface{}
			errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, dato)
			if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoInfoComeplementaria)
				return APIResponseDTO
			} else {
				fmt.Println("Info complementaria registrada", dato["InfoComplementariaId"])
				// alertas = append(alertas, Transferencia)
			}
		}

		var resultadoInscripcionPregrado map[string]interface{}
		errInscripcionPregrado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado", "POST", &resultadoInscripcionPregrado, InscripcionPregrado)
		if resultadoInscripcionPregrado["Type"] == "error" || errInscripcionPregrado != nil || resultadoInscripcionPregrado["Status"] == "404" || resultadoInscripcionPregrado["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoInscripcionPregrado)
			return APIResponseDTO
		} else {
			fmt.Println("Inscripcion registrada")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, InfoIcfesColegio)
		}

		// Registro de colegio a tercero
		NombrePrograma := fmt.Sprintf("%q", "colegio")
		FechaI := fmt.Sprintf("%q", date)
		colegioId, _ := json.Marshal(map[string]interface{}{"Id": InformacionColegio["Id"].(float64)})

		ColegioRegistro := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": Tercero["TerceroId"].(map[string]interface{})["Id"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 313},
			"Dato": "{\"ProgramaAcademico\": " + NombrePrograma + ",    " +
				"\"FechaInicio\": " + FechaI + ",    " +
				"\"NitUniversidad\": " + string(colegioId) + "}",
			"Activo": true,
		}

		var resultadoRegistroColegio map[string]interface{}

		errRegistroColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &resultadoRegistroColegio, ColegioRegistro)
		if errRegistroColegio == nil && fmt.Sprintf("%v", resultadoRegistroColegio["System"]) != "map[]" && resultadoRegistroColegio["Id"] != nil {
			if resultadoRegistroColegio["Status"] != 400 {
				fmt.Println("Colegio registrado")
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, InfoIcfesColegio)
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoRegistroColegio)
				return APIResponseDTO
			}
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoRegistroColegio)
			return APIResponseDTO
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func PreinscripcionPost(data []byte) (APIResponseDTO requestresponse.APIResponse){
	var Infopreinscripcion map[string]interface{}
	if err := json.Unmarshal(data, &Infopreinscripcion); err == nil {

		var InfoPreinscripcionTodas = Infopreinscripcion["DatosPreinscripcion"].([]interface{})
		for _, datoPreinscripcion := range InfoPreinscripcionTodas {
			var dato = datoPreinscripcion.(map[string]interface{})

			var resultadoPreinscripcion map[string]interface{}
			errPreinscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion", "POST", &resultadoPreinscripcion, dato)
			if resultadoPreinscripcion["Type"] == "error" || errPreinscripcion != nil || resultadoPreinscripcion["Status"] == "404" || resultadoPreinscripcion["Message"] != nil {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoPreinscripcion)
				return APIResponseDTO
			} else {
				fmt.Println("Preinscripcion registrada", dato)
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, InfoPreinscripcionTodas)
			}
		}

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func InfoNuevoColegioIcfes(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var InfoIcfesColegio map[string]interface{}
	var IdColegio float64
	if err := json.Unmarshal(data, &InfoIcfesColegio); err == nil {

		var InscripcionPregrado = InfoIcfesColegio["InscripcionPregrado"].(map[string]interface{})
		var InfoComplementariaTercero = InfoIcfesColegio["InfoComplementariaTercero"].(map[string]interface{})
		var InformacionColegio = InfoIcfesColegio["TerceroColegio"].(map[string]interface{})
		var InformacionDireccionColegio = InfoIcfesColegio["DireccionColegio"].(map[string]interface{})
		var InformacionUbicacionColegio = InfoIcfesColegio["UbicacionColegio"].(map[string]interface{})
		var InformaciontipoColegio = InfoIcfesColegio["TipoColegio"].(map[string]interface{})
		var Tercero = InfoIcfesColegio["Tercero"].(map[string]interface{})
		var date = time.Now()

		var resultadoRegistroColegio map[string]interface{}
		errRegistroColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &resultadoRegistroColegio, InformacionColegio)
		if resultadoRegistroColegio["Type"] == "error" || errRegistroColegio != nil || resultadoRegistroColegio["Status"] == "404" || resultadoRegistroColegio["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoRegistroColegio)
			return APIResponseDTO
		} else {
			fmt.Println("Colegio registrado")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultadoRegistroColegio)
			IdColegio = resultadoRegistroColegio["Id"].(float64)
			fmt.Println(IdColegio)
		}
		DireccionColegioPost := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": IdColegio},
			"InfoComplementariaId": map[string]interface{}{"Id": InformacionDireccionColegio["InfoComplementariaId"].(map[string]interface{})["Id"].(float64)},
			"Dato":                 InformacionDireccionColegio["Dato"],
			"Activo":               true,
		}

		var resultadoDirecionColegio map[string]interface{}
		errRegistroDirecionColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoDirecionColegio, DireccionColegioPost)
		if resultadoDirecionColegio["Type"] == "error" || errRegistroDirecionColegio != nil || resultadoDirecionColegio["Status"] == "404" || resultadoDirecionColegio["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoDirecionColegio)
				return APIResponseDTO
		} else {
			fmt.Println("Direccion Colegio registrado")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultadoDirecionColegio)

		}
		UbicacionColegioPost := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": IdColegio},
			"InfoComplementariaId": map[string]interface{}{"Id": InformacionUbicacionColegio["InfoComplementariaId"].(map[string]interface{})["Id"].(float64)},
			"Dato":                 InformacionUbicacionColegio["Dato"],
			"Activo":               true,
		}
		var resultadoUbicacionColegio map[string]interface{}
		errRegistroUbicacionColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoUbicacionColegio, UbicacionColegioPost)
		if resultadoUbicacionColegio["Type"] == "error" || errRegistroUbicacionColegio != nil || resultadoUbicacionColegio["Status"] == "404" || resultadoUbicacionColegio["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoUbicacionColegio)
				return APIResponseDTO
		} else {
			fmt.Println("Ubicacion Colegio registrado")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultadoUbicacionColegio)

		}
		tipoColegioPost := map[string]interface{}{
			"TerceroId":     map[string]interface{}{"Id": IdColegio},
			"TipoTerceroId": map[string]interface{}{"Id": InformaciontipoColegio["TipoTerceroId"].(map[string]interface{})["Id"].(float64)},
			"Activo":        true,
		}

		var resultadoTipoColegio map[string]interface{}
		errRegistroTipoColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &resultadoTipoColegio, tipoColegioPost)
		if resultadoTipoColegio["Type"] == "error" || errRegistroTipoColegio != nil || resultadoTipoColegio["Status"] == "404" || resultadoTipoColegio["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoTipoColegio)
			return APIResponseDTO
		} else {
			fmt.Println("TipoColegio registrado")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultadoTipoColegio)

		}

		VerificarColegioPost := map[string]interface{}{
			"TerceroId":     map[string]interface{}{"Id": IdColegio},
			"TipoTerceroId": map[string]interface{}{"Id": 14},
			"Activo":        true,
		}

		var resultadoVerificarColegio map[string]interface{}
		errRegistroVerificarColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &resultadoVerificarColegio, VerificarColegioPost)
		if resultadoVerificarColegio["Type"] == "error" || errRegistroVerificarColegio != nil || resultadoVerificarColegio["Status"] == "404" || resultadoVerificarColegio["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoVerificarColegio)
			return APIResponseDTO
		} else {
			fmt.Println("Verificar registrado")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultadoVerificarColegio)

		}
		// Registro de colegio a tercero

		// Registro de colegio a tercero
		NombrePrograma := fmt.Sprintf("%q", "colegio")
		FechaI := fmt.Sprintf("%q", date)
		colegioId, _ := json.Marshal(map[string]interface{}{"Id": IdColegio})

		ColegioRegistro := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": Tercero["TerceroId"].(map[string]interface{})["Id"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 313},
			"Dato": "{\"ProgramaAcademico\": " + NombrePrograma + ",    " +
				"\"FechaInicio\": " + FechaI + ",    " +
				"\"NitUniversidad\": " + string(colegioId) + "}",
			"Activo": true,
		}

		var resultadoRegistroColegioTercero map[string]interface{}

		errRegistroColegioTercero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &resultadoRegistroColegioTercero, ColegioRegistro)
		if errRegistroColegioTercero == nil && fmt.Sprintf("%v", resultadoRegistroColegioTercero["System"]) != "map[]" && resultadoRegistroColegioTercero["Id"] != nil {
			if resultadoRegistroColegioTercero["Status"] != 400 {
				fmt.Println("Colegio Tercero registrado")
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, InfoIcfesColegio)
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoRegistroColegioTercero)
				return APIResponseDTO
			}
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoRegistroColegioTercero)
				return APIResponseDTO
		}

		var resultadoInfoComeplementaria map[string]interface{}

		errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, InfoComplementariaTercero)
		if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoInfoComeplementaria)
			return APIResponseDTO
		} else {
			fmt.Println("Info complementaria registrada", InfoComplementariaTercero)
			// alertas = append(alertas, Transferencia)
		}

		var resultadoInscripcionPregrado map[string]interface{}
		errInscripcionPregrado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado", "POST", &resultadoInscripcionPregrado, InscripcionPregrado)
		if resultadoInscripcionPregrado["Type"] == "error" || errInscripcionPregrado != nil || resultadoInscripcionPregrado["Status"] == "404" || resultadoInscripcionPregrado["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoInscripcionPregrado)
			return APIResponseDTO
		} else {
			fmt.Println("Inscripcion registrada")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, InfoIcfesColegio)
		}

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func PutInfoComplementaria(data []byte) (APIResponseDTO requestresponse.APIResponse){
	var InfoComplementariaUniversidad map[string]interface{}
	if err := json.Unmarshal(data, &InfoComplementariaUniversidad); err == nil {

		var InfoComplementariaTercero = InfoComplementariaUniversidad["InfoComplementariaTercero"].([]interface{})
		var date = time.Now()

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var dato = datoInfoComplementaria.(map[string]interface{})
			dato["FechaCreacion"] = date
			dato["FechaModificacion"] = date
			var resultadoInfoComeplementaria map[string]interface{}
			errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, dato)
			if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, resultadoInfoComeplementaria)
				return APIResponseDTO
			} else {
				fmt.Println("Info complementaria registrada", dato["InfoComplementariaId"])
				// alertas = append(alertas, Transferencia)
			}
		}

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func ConsultarEventos(idEvento string) (APIResponseDTO requestresponse.APIResponse) {
		// resultado datos complementarios persona
		var resultado []map[string]interface{}
		var EventosInscripcion []map[string]interface{}
	
		erreVentos := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/?query=Activo:true,EventoPadreId:"+idEvento+"&limit=0", &EventosInscripcion)
		if erreVentos == nil && fmt.Sprintf("%v", EventosInscripcion[0]) != "map[]" {
			if EventosInscripcion[0]["Status"] != 404 {
	
				var Proyectos_academicos []map[string]interface{}
				var Proyectos_academicos_Get []map[string]interface{}
				for i := 0; i < len(EventosInscripcion); i++ {
					if len(EventosInscripcion) > 0 {
						proyectoacademico := EventosInscripcion[i]["TipoEventoId"].(map[string]interface{})
	
						var ProyectosAcademicosConEvento map[string]interface{}
	
						erreproyectos := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%v", proyectoacademico["DependenciaId"]), &ProyectosAcademicosConEvento)
						if erreproyectos == nil && fmt.Sprintf("%v", ProyectosAcademicosConEvento) != "map[]" {
							if ProyectosAcademicosConEvento["Status"] != 404 {
								periodoevento := EventosInscripcion[i]["PeriodoId"]
								fmt.Println(periodoevento)
								ProyectosAcademicosConEvento["PeriodoId"] = map[string]interface{}{"Id": periodoevento}
								Proyectos_academicos_Get = append(Proyectos_academicos_Get, ProyectosAcademicosConEvento)
	
							} else {
								if ProyectosAcademicosConEvento["Message"] == "Not found resource" {
									APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil)
								} else {
									logs.Error(ProyectosAcademicosConEvento)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									APIResponseDTO = requestresponse.APIResponseDTO(false, 404, erreproyectos)
									return APIResponseDTO
								}
							}
						} else {
							logs.Error(ProyectosAcademicosConEvento)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							APIResponseDTO = requestresponse.APIResponseDTO(false, 404, erreproyectos)
							return APIResponseDTO
						}
	
						Proyectos_academicos = append(Proyectos_academicos, proyectoacademico)
	
					}
				}
				resultado = Proyectos_academicos_Get
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado)
	
			} else {
				if EventosInscripcion[0]["Message"] == "Not found resource" {
					APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil)
				} else {
					logs.Error(EventosInscripcion)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, erreVentos)
					return APIResponseDTO
				}
			}
		} else {
			logs.Error(EventosInscripcion)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			APIResponseDTO = requestresponse.APIResponseDTO(false, 404, erreVentos)
			return APIResponseDTO
		}

		return APIResponseDTO
}

func InfoComplementariaTercero(data []byte) (APIResponseDTO requestresponse.APIResponse){
	var InfoComplementaria map[string]interface{}

	var algoFallo bool = false

	var inactivePosts []map[string]interface{}

	var respuestas []interface{}

	if err := json.Unmarshal(data, &InfoComplementaria); err == nil {

		var InfoComplementariaTercero = InfoComplementaria["InfoComplementariaTercero"].([]interface{})
		var date = time_bogota.TiempoBogotaFormato()

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var dato = datoInfoComplementaria.(map[string]interface{})
			dato["FechaCreacion"] = date
			dato["FechaModificacion"] = date
			var resultadoInfoComeplementaria map[string]interface{}
			errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, dato)
			if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Status"] == "400" || resultadoInfoComeplementaria["Message"] != nil {
				algoFallo = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400,  nil, errInfoComplementaria.Error())
			} else {
				respuestas = append(respuestas, resultadoInfoComeplementaria)
				inactivePosts = append(inactivePosts, resultadoInfoComeplementaria)
			}
			if algoFallo {
				break
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400,  nil, err.Error())
	}

	if !algoFallo {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200,  respuestas)
	} else {
		for _, disable := range inactivePosts {
			helpers.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + fmt.Sprintf("%.f", disable["Id"].(float64)))
		}
		return APIResponseDTO
	}
	return APIResponseDTO
}

