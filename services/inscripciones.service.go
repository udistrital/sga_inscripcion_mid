package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_inscripcion_mid/helpers"
	"github.com/udistrital/sga_inscripcion_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/time_bogota"
	"golang.org/x/sync/errgroup"
)

func EstadoInscripcion(idPersona string, idPeriodo string) (APIResponseDTO requestresponse.APIResponse) {

	recibosResultado, err := helpers.VerificarRecibos(idPersona, idPeriodo)

	if err == "" {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, recibosResultado, nil)
	} else if err == "400" {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Bad request")
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
	}

	return APIResponseDTO
}

func InformacionFamiliar(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var InformacionFamiliar map[string]interface{}
	var TerceroFamiliarPost map[string]interface{}
	var FamiliarParentescoPost map[string]interface{}
	var InfoContactoPost map[string]interface{}

	if err := json.Unmarshal(data, &InformacionFamiliar); err == nil {
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
										APIResponseDTO = requestresponse.APIResponseDTO(true, 200, TerceroFamiliarPost, nil)
									} else {
										logs.Error(errFamiliarParentesco)
										APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errFamiliarParentesco)
										return APIResponseDTO
									}
								} else {
									//var resultado2 map[string]interface{}
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
									helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/%.f", FamiliarParentescoPost["Id"]), "DELETE", &resultado2, nil)
									helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/%.f", FamiliarParentescoPost["Id"]))
									logs.Error(errFamiliarParentesco)
									APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errFamiliarParentesco)
									return APIResponseDTO
								}
							}
						} else {
							//var resultado2 map[string]interface{}
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
							helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
							logs.Error(errFamiliarParentesco)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errFamiliarParentesco)
							return APIResponseDTO
						}
					} else {
						//var resultado2 map[string]interface{}
						//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
						helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
						logs.Error(errFamiliarParentesco)
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errFamiliarParentesco)
						return APIResponseDTO
					}

				} else {
					//var resultado2 map[string]interface{}
					//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
					helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
					logs.Error(errTerceroFamiliar)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errTerceroFamiliar)
					return APIResponseDTO
				}
			} else {
				logs.Error(errTerceroFamiliar)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errTerceroFamiliar)
				return APIResponseDTO
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	return APIResponseDTO
}

func Reintegro(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var Reintegro map[string]interface{}
	if err := json.Unmarshal(data, &Reintegro); err == nil {

		var resultadoReintegro map[string]interface{}
		errReintegro := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"tr_inscripcion/reintegro", "POST", &resultadoReintegro, Reintegro)
		if resultadoReintegro["Type"] == "error" || errReintegro != nil || resultadoReintegro["Status"] == "404" || resultadoReintegro["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoReintegro)
			return APIResponseDTO
		} else {
			fmt.Println("Reintegrro registrado")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, Reintegro, nil)
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func TransferenciaPost(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var Transferencia map[string]interface{}

	if err := json.Unmarshal(data, &Transferencia); err == nil {

		var resultadoTransferencia map[string]interface{}
		errTransferencia := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"tr_inscripcion/transferencia", "POST", &resultadoTransferencia, Transferencia)
		if resultadoTransferencia["Type"] == "error" || errTransferencia != nil || resultadoTransferencia["Status"] == "404" || resultadoTransferencia["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoTransferencia)
			return APIResponseDTO
		} else {
			fmt.Println("Transferencia registrada")
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, Transferencia, nil)
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func InfoIcfesColegio(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var InfoIcfesColegio map[string]interface{}
	resultado := []interface{}{}

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
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoInfoComeplementaria)
				return APIResponseDTO
			} else {
				fmt.Println("Info complementaria registrada", dato["InfoComplementariaId"])
				// alertas = append(alertas, Transferencia)
			}
		}

		var resultadoInscripcionPregrado map[string]interface{}
		errInscripcionPregrado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado", "POST", &resultadoInscripcionPregrado, InscripcionPregrado)
		if resultadoInscripcionPregrado["Type"] == "error" || errInscripcionPregrado != nil || resultadoInscripcionPregrado["Status"] == "404" || resultadoInscripcionPregrado["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoInscripcionPregrado)
			return APIResponseDTO
		} else {
			fmt.Println("Inscripcion registrada")
			resultado = append(resultado, InfoIcfesColegio)
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
				resultado = append(resultado, InfoIcfesColegio)
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoRegistroColegio)
				return APIResponseDTO
			}
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoRegistroColegio)
			return APIResponseDTO
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}

	if len(resultado) > 0 {
		return requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}
}

func PreinscripcionPost(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var Infopreinscripcion map[string]interface{}
	if err := json.Unmarshal(data, &Infopreinscripcion); err == nil {

		var InfoPreinscripcionTodas = Infopreinscripcion["DatosPreinscripcion"].([]interface{})
		for _, datoPreinscripcion := range InfoPreinscripcionTodas {
			var dato = datoPreinscripcion.(map[string]interface{})

			var resultadoPreinscripcion map[string]interface{}
			errPreinscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion", "POST", &resultadoPreinscripcion, dato)
			if resultadoPreinscripcion["Type"] == "error" || errPreinscripcion != nil || resultadoPreinscripcion["Status"] == "404" || resultadoPreinscripcion["Message"] != nil {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoPreinscripcion)
				return APIResponseDTO
			} else {
				fmt.Println("Preinscripcion registrada", dato)
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, InfoPreinscripcionTodas, nil)
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
	resultado := []interface{}{}
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
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoRegistroColegio)
			return APIResponseDTO
		} else {
			fmt.Println("Colegio registrado")
			resultado = append(resultado, resultadoRegistroColegio)
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
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoDirecionColegio)
			return APIResponseDTO
		} else {
			fmt.Println("Direccion Colegio registrado")
			resultado = append(resultado, resultadoDirecionColegio)
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
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoUbicacionColegio)
			return APIResponseDTO
		} else {
			fmt.Println("Ubicacion Colegio registrado")
			resultado = append(resultado, resultadoUbicacionColegio)

		}
		tipoColegioPost := map[string]interface{}{
			"TerceroId":     map[string]interface{}{"Id": IdColegio},
			"TipoTerceroId": map[string]interface{}{"Id": InformaciontipoColegio["TipoTerceroId"].(map[string]interface{})["Id"].(float64)},
			"Activo":        true,
		}

		var resultadoTipoColegio map[string]interface{}
		errRegistroTipoColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &resultadoTipoColegio, tipoColegioPost)
		if resultadoTipoColegio["Type"] == "error" || errRegistroTipoColegio != nil || resultadoTipoColegio["Status"] == "404" || resultadoTipoColegio["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoTipoColegio)
			return APIResponseDTO
		} else {
			fmt.Println("TipoColegio registrado")
			resultado = append(resultado, resultadoTipoColegio)

		}

		VerificarColegioPost := map[string]interface{}{
			"TerceroId":     map[string]interface{}{"Id": IdColegio},
			"TipoTerceroId": map[string]interface{}{"Id": 14},
			"Activo":        true,
		}

		var resultadoVerificarColegio map[string]interface{}
		errRegistroVerificarColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &resultadoVerificarColegio, VerificarColegioPost)
		if resultadoVerificarColegio["Type"] == "error" || errRegistroVerificarColegio != nil || resultadoVerificarColegio["Status"] == "404" || resultadoVerificarColegio["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoVerificarColegio)
			return APIResponseDTO
		} else {
			fmt.Println("Verificar registrado")
			resultado = append(resultado, resultadoVerificarColegio)
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
				resultado = append(resultado, InfoIcfesColegio)
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoRegistroColegioTercero)
				return APIResponseDTO
			}
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoRegistroColegioTercero)
			return APIResponseDTO
		}

		var resultadoInfoComeplementaria map[string]interface{}

		errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, InfoComplementariaTercero)
		if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoInfoComeplementaria)
			return APIResponseDTO
		} else {
			fmt.Println("Info complementaria registrada", InfoComplementariaTercero)
			// alertas = append(alertas, Transferencia)
		}

		var resultadoInscripcionPregrado map[string]interface{}
		errInscripcionPregrado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado", "POST", &resultadoInscripcionPregrado, InscripcionPregrado)
		if resultadoInscripcionPregrado["Type"] == "error" || errInscripcionPregrado != nil || resultadoInscripcionPregrado["Status"] == "404" || resultadoInscripcionPregrado["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoInscripcionPregrado)
			return APIResponseDTO
		} else {
			fmt.Println("Inscripcion registrada")
			resultado = append(resultado, InfoIcfesColegio)
		}

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}

	if len(resultado) > 0 {
		return requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}

}

func PutInfoComplementaria(data []byte) (APIResponseDTO requestresponse.APIResponse) {
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
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultadoInfoComeplementaria)
				return APIResponseDTO
			} else {
				fmt.Println("Info complementaria registrada", dato["InfoComplementariaId"])
				//APIResponseDTO = requestresponse.APIResponseDTO(true, 200, transferencia ,nil)
			}
		}

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func ConsultarEventos(idEvento string) (APIResponseDTO requestresponse.APIResponse) {
	// resultado datos complementarios persona
	var resultado []map[string]interface{}
	var EventosInscripcionMap []map[string]interface{}
	wge := new(errgroup.Group)
	var mutex sync.Mutex // Mutex para proteger el acceso a resultados

	erreVentos := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/?query=Activo:true,EventoPadreId:"+idEvento+"&limit=0", &EventosInscripcionMap)
	if erreVentos == nil && fmt.Sprintf("%v", EventosInscripcionMap[0]) != "[map[]]" {
		if EventosInscripcionMap[0]["Status"] != 404 {

			var Proyectos_academicos []map[string]interface{}
			var Proyectos_academicos_Get []map[string]interface{}
			wge.SetLimit(-1)
			for _, EventosInscripcion := range EventosInscripcionMap {
				EventosInscripcion = EventosInscripcion
				wge.Go(func() error {

					if len(EventosInscripcion) > 0 {
						proyectoacademico := EventosInscripcion["TipoEventoId"].(map[string]interface{})

						var ProyectosAcademicosConEvento map[string]interface{}

						erreproyectos := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%v", proyectoacademico["DependenciaId"]), &ProyectosAcademicosConEvento)
						if erreproyectos == nil && fmt.Sprintf("%v", ProyectosAcademicosConEvento) != "map[]" {
							if ProyectosAcademicosConEvento["Status"] != 404 {
								periodoevento := EventosInscripcion["PeriodoId"]
								fmt.Println(periodoevento)
								ProyectosAcademicosConEvento["PeriodoId"] = map[string]interface{}{"Id": periodoevento}
								Proyectos_academicos_Get = append(Proyectos_academicos_Get, ProyectosAcademicosConEvento)

							} else {
								if ProyectosAcademicosConEvento["Message"] == "Not found resource" {
									return errors.New("No data found")
								} else {
									logs.Error(ProyectosAcademicosConEvento)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									return erreproyectos
								}
							}
						} else {
							logs.Error(ProyectosAcademicosConEvento)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							return erreproyectos
						}
						mutex.Lock()
						Proyectos_academicos = append(Proyectos_academicos, proyectoacademico)
						mutex.Unlock()
					} else {
						return errors.New("No data found")
					}
					return nil
				})
			}
			//Si existe error, se realiza
			if err := wge.Wait(); err != nil {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
				return APIResponseDTO
			}
			resultado = Proyectos_academicos_Get
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)

		} else {
			if EventosInscripcionMap[0]["Message"] == "Not found resource" {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				return APIResponseDTO
			} else {
				logs.Error(EventosInscripcionMap)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, erreVentos)
				return APIResponseDTO
			}
		}
	} else {
		logs.Error(EventosInscripcionMap)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, erreVentos)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func InfoComplementariaTercero(data []byte) (APIResponseDTO requestresponse.APIResponse) {
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
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInfoComplementaria.Error())
			} else {
				respuestas = append(respuestas, resultadoInfoComeplementaria)
				inactivePosts = append(inactivePosts, resultadoInfoComeplementaria)
			}
			if algoFallo {
				break
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !algoFallo {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuestas, nil)
	} else {
		for _, disable := range inactivePosts {
			helpers.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + fmt.Sprintf("%.f", disable["Id"].(float64)))
		}
		return APIResponseDTO
	}
	return APIResponseDTO
}

func GetInfoCompTercero(idTercero string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado consulta
	resultado := map[string]interface{}{}
	// var resultado map[string]interface{}
	var errorGetAll bool

	// 41 = estrato
	IdEstrato, _ := helpers.IdInfoCompTercero("9", "ESTRATO")
	var resultadoEstrato []map[string]interface{}
	errEstratoResidencia := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdEstrato+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoEstrato)
	if errEstratoResidencia == nil && fmt.Sprintf("%v", resultadoEstrato[0]["System"]) != "map[]" {
		if resultadoEstrato[0]["Status"] != 404 && resultadoEstrato[0]["Id"] != nil {
			resultado["IdEstratoEnte"] = resultadoEstrato[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoEstrato[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["EstratoResidencia"] = nil
			} else {
				resultado["EstratoResidencia"] = estratoJson["value"]
			}
		} else {
			if resultadoEstrato[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errEstratoResidencia)
	}

	// 55 = codigo postal
	IdCodPostal, _ := helpers.IdInfoCompTercero("10", "CODIGO_POSTAL")
	var resultadoCodigoPostal []map[string]interface{}
	errCodigoPostal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCodPostal+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoCodigoPostal)
	if errCodigoPostal == nil && fmt.Sprintf("%v", resultadoCodigoPostal[0]["System"]) != "map[]" {
		if resultadoCodigoPostal[0]["Status"] != 404 && resultadoCodigoPostal[0]["Id"] != nil {
			resultado["IdCodigoEnte"] = resultadoCodigoPostal[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoCodigoPostal[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["CodigoPostal"] = nil
			} else {
				resultado["CodigoPostal"] = estratoJson["value"]
			}
		} else {
			if resultadoCodigoPostal[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errCodigoPostal)
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errCodigoPostal)
	}

	// 51 = telefono
	IdTelefono, _ := helpers.IdInfoCompTercero("10", "TELEFONO")
	var resultadoTelefono []map[string]interface{}
	errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdTelefono+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoTelefono)
	if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
		if resultadoTelefono[0]["Status"] != 404 && resultadoTelefono[0]["Id"] != nil {
			resultado["IdTelefonoEnte"] = resultadoTelefono[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["Telefono"] = nil
				resultado["TelefonoAlterno"] = nil
			} else {
				resultado["Telefono"] = estratoJson["principal"]
				resultado["TelefonoAlterno"] = estratoJson["alterno"]
			}
		} else {
			if resultadoTelefono[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errTelefono)
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errTelefono)
	}

	// 54 = direccion
	IdDireccion, _ := helpers.IdInfoCompTercero("10", "DIRECCIÓN")
	var resultadoDireccion []map[string]interface{}
	errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdDireccion+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoDireccion)
	if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
		if resultadoDireccion[0]["Status"] != 404 && resultadoDireccion[0]["Id"] != nil {
			resultado["IdLugarEnte"] = resultadoDireccion[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["PaisResidencia"] = nil
				resultado["DepartamentoResidencia"] = nil
				resultado["CiudadResidencia"] = nil
				resultado["DireccionResidencia"] = nil
			} else {
				resultado["PaisResidencia"] = estratoJson["country"]
				resultado["DepartamentoResidencia"] = estratoJson["department"]
				resultado["CiudadResidencia"] = estratoJson["city"]
				resultado["DireccionResidencia"] = estratoJson["address"]

			}
		} else {
			if resultadoDireccion[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDireccion)
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDireccion)
	}

	// Correo registro
	IdCorreo, _ := helpers.IdInfoCompTercero("10", "CORREO")
	var resultadoCorreo []map[string]interface{}
	errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCorreo+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoCorreo)
	if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
		if resultadoCorreo[0]["Status"] != 404 && resultadoCorreo[0]["Id"] != nil {
			resultado["IdCorreo"] = resultadoCorreo[0]["Id"]
			// unmarshall dato
			var correoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJson); err != nil {
				resultado["Correo"] = nil
			} else {
				resultado["Correo"] = correoJson["value"]
			}
		} else {
			if resultadoCorreo[0]["Message"] == "Not found resource" {
				/* //errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			} else {
				/* //errorGetAll = true
				alertas = append(alertas, errCorreo)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			}
		}
	} else {
		/* //errorGetAll = true
		alertas = append(alertas, errCorreo)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta} */
	}

	// Correo alterno
	IdCorreoAlterno, _ := helpers.IdInfoCompTercero("10", "CORREOALTER")
	var resultadoCorreoAlterno []map[string]interface{}
	errCorreoAlterno := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCorreoAlterno+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoCorreoAlterno)
	if errCorreoAlterno == nil && fmt.Sprintf("%v", resultadoCorreoAlterno[0]["System"]) != "map[]" {
		if resultadoCorreoAlterno[0]["Status"] != 404 && resultadoCorreoAlterno[0]["Id"] != nil {
			resultado["IdCorreoAlterno"] = resultadoCorreoAlterno[0]["Id"]
			// unmarshall dato
			var correoAlternoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoCorreoAlterno[0]["Dato"].(string)), &correoAlternoJson); err != nil {
				resultado["CorreoAlterno"] = nil
			} else {
				resultado["CorreoAlterno"] = correoAlternoJson["value"]
			}
		} else {
			if resultadoCorreoAlterno[0]["Message"] == "Not found resource" {
				/* //errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			} else {
				/* //errorGetAll = true
				alertas = append(alertas, errCorreoAlterno)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			}
		}
	} else {
		/* //errorGetAll = true
		alertas = append(alertas, errCorreoAlterno)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta} */
	}

	if !errorGetAll {
		return requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}
}

func ActualizarInfoContact(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var InfoContacto map[string]interface{}

	var algoFallo bool = false

	var revertPuts []map[string]interface{}
	var inactivePosts []map[string]interface{}

	var respuestas []interface{}

	if err := json.Unmarshal(data, &InfoContacto); err == nil {
		var InfoComplementariaTercero = InfoContacto["InfoComplementariaTercero"].([]interface{})

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var InfoComplementaria = datoInfoComplementaria.(map[string]interface{})

			var getInfoComp map[string]interface{}
			id := InfoComplementaria["Id"].(float64)
			if id > 0 {
				errGetInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", id), &getInfoComp)
				if errGetInfoComp == nil && getInfoComp["Status"] != "404" && getInfoComp["Status"] != "400" {
					putInfoComp := getInfoComp
					revertPuts = append(revertPuts, getInfoComp)
					putInfoComp["TerceroId"] = InfoComplementaria["TerceroId"]
					putInfoComp["InfoComplementariaId"] = InfoComplementaria["InfoComplementariaId"]
					putInfoComp["Dato"] = InfoComplementaria["Dato"].(string)
					putInfoComp["Activo"] = InfoComplementaria["Activo"]
					var resp map[string]interface{}
					errPutInfoComp := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", id), "PUT", &resp, putInfoComp)
					if errPutInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
						respuestas = append(respuestas, resp)
					} else {
						algoFallo = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPutInfoComp.Error())
					}
				} else {
					algoFallo = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
				}
			} else {
				var resp map[string]interface{}
				errPostInfoComp := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resp, InfoComplementaria)
				if errPostInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
					respuestas = append(respuestas, resp)
					inactivePosts = append(inactivePosts, resp)
				} else {
					algoFallo = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPostInfoComp.Error())
				}
			}
			if algoFallo {
				break
			}
		}
	} else {
		algoFallo = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !algoFallo {
		return requestresponse.APIResponseDTO(true, 200, respuestas, nil)
	} else {
		for _, revert := range revertPuts {
			var resp map[string]interface{}
			request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", revert["Id"].(float64)), "PUT", &resp, revert)
		}
		for _, disable := range inactivePosts {
			helpers.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + fmt.Sprintf("%.f", disable["Id"].(float64)))
		}
		return APIResponseDTO
	}

}

func GenerarInscripcion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	// var reciboVencido bool
	var SolicitudInscripcion map[string]interface{}
	var TipoParametro string
	var parametro map[string]interface{}
	var Valor map[string]interface{}
	var NuevoRecibo map[string]interface{}
	var inscripcionRealizada map[string]interface{}
	// var contadorRecibos int

	if err := json.Unmarshal(data, &SolicitudInscripcion); err == nil {
		objTransaccion := map[string]interface{}{
			"codigo":              SolicitudInscripcion["Id"].(float64),
			"nombre":              SolicitudInscripcion["Nombre"].(string),
			"apellido":            SolicitudInscripcion["Apellido"].(string),
			"correo":              SolicitudInscripcion["Correo"].(string),
			"proyecto":            SolicitudInscripcion["ProgramaAcademicoCodigo"].(float64),
			"tiporecibo":          15, // se define 15 por que es el id definido en el api de recibos para inscripcion
			"concepto":            "",
			"valorordinario":      0,
			"valorextraordinario": 0,
			"cuota":               1,
			"fechaordinario":      SolicitudInscripcion["FechaPago"].(string),
			"fechaextraordinario": SolicitudInscripcion["FechaPago"].(string),
			"aniopago":            SolicitudInscripcion["Year"].(float64),
			"perpago":             SolicitudInscripcion["Periodo"].(float64),
		}

		inscripcion := map[string]interface{}{
			"PersonaId":           SolicitudInscripcion["PersonaId"].(float64),
			"ProgramaAcademicoId": SolicitudInscripcion["ProgramaAcademicoId"].(float64),
			"ReciboInscripcion":   "",
			"PeriodoId":           SolicitudInscripcion["PeriodoId"].(float64),
			"AceptaTerminos":      true,
			"FechaAceptaTerminos": time.Now(),
			"Activo":              true,
			"EstadoInscripcionId": map[string]interface{}{"Id": 1},
			"TipoInscripcionId":   map[string]interface{}{"Id": SolicitudInscripcion["TipoInscripcionId"]},
		}

		if SolicitudInscripcion["Nivel"].(float64) == 1 {
			TipoParametro = "13"
			id_periodo := int(SolicitudInscripcion["PeriodoId"].(float64))
			credencial := helpers.GenerarCredencialInscripcionPregrado(float64(id_periodo))
			inscripcion["Credencial"] = credencial
		} else if SolicitudInscripcion["Nivel"].(float64) == 2 {
			TipoParametro = "12"
		}

		// persona_id := fmt.Sprintf("%d", int(SolicitudInscripcion["PersonaId"].(float64)))
		// id_periodo := fmt.Sprintf("%d", int(SolicitudInscripcion["PeriodoId"].(float64)))
		// //id_programa_academico := fmt.Sprintf("%d", int(SolicitudInscripcion["ProgramaAcademicoId"].(float64)))

		// recibosResultado, err := helpers.VerificarRecibos(persona_id, id_periodo)

		// if err == "" {
		// if inscripciones, ok := recibosResultado["Inscripciones"]; ok {
		// 	// Convertir la variable de tipo interface{} a un slice de mapas
		// 	inscripcionesMap, ok := inscripciones.([]map[string]interface{})
		// 	if len(inscripcionesMap) > 0 && ok {
		// 		for i := 0; i < len(inscripcionesMap); i++ {
		// 			if inscripcionesMap[i]["ProgramaAcademicoId"] != nil {
		// 				// id_programa_inscripciones := fmt.Sprintf("%d", int(inscripcionesMap[i]["ProgramaAcademicoId"].(float64)))
		// 				estado_recibo_inscripciones := inscripcionesMap[i]["Estado"].(string)
		// 				// if id_programa_inscripciones == id_programa_academico {
		// 				if estado_recibo_inscripciones == "Vencido" {
		// 					reciboVencido = true
		// 				} else {
		// 					reciboVencido = false
		// 				}
		// 				// } else {
		// 				// 	contadorRecibos++
		// 				// }
		// 			}
		// 		}
		// 		if contadorRecibos == len(inscripcionesMap) {
		// 			reciboVencido = true
		// 		}
		// 	}

		// }

		coincideCodigoSnies := false
		var proyectos []map[string]interface{}
		idProyecto := fmt.Sprintf("%.0f", SolicitudInscripcion["ProgramaAcademicoId"].(float64))
		errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/"+idProyecto, &proyectos)
		if errproyecto == nil {
			proyecto := proyectos[0]
			proyectoAcademico := proyecto["ProyectoAcademico"].(map[string]interface{})
			codigoSnies := proyectoAcademico["CodigoSnies"].(string)

			var HomologacionXML map[string]interface{}
			codigoProyecto := fmt.Sprintf("%.0f", SolicitudInscripcion["ProgramaAcademicoCodigo"].(float64))
			errHomologacion := request.GetJsonWSO2("http://"+beego.AppConfig.String("HomologacionDependenciaService")+"proyecto_acad_snies/"+codigoSnies, &HomologacionXML)
			resultadoHomologacion := HomologacionXML["proyecto_snies"].(map[string]interface{})
			if errHomologacion == nil && fmt.Sprintf("%v", resultadoHomologacion) != "map[]" {
				proyectosSnies := resultadoHomologacion["proyectos"].([]interface{})

				for _, proyectoSnies := range proyectosSnies {
					proyectoSnies := proyectoSnies.(map[string]interface{})
					codigoProyectoHomologacion := fmt.Sprintf("%.0f", proyectoSnies["codigo_proyecto"].(float64))
					coincideCodigoSnies = codigoProyectoHomologacion == codigoProyecto

					if coincideCodigoSnies {
						break
					}
				}
			}
		}

		//Verificar si existe un recibo vencido o es la primera vez que inscribe el postgrado
		// if reciboVencido || fmt.Sprintf("%v", recibosResultado) == "map[]" {
		if coincideCodigoSnies {
			errInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion", "POST", &inscripcionRealizada, inscripcion)
			if errInscripcion == nil && inscripcionRealizada["Status"] != "400" {
				var respCambio map[string]interface{}
				estado := helpers.GetEstadoInscripcion(inscripcionRealizada)
				inscripcionEvolucionEstado := helpers.GenerarInscripcionEvolucionEstado(int(inscripcionRealizada["Id"].(float64)), nil, helpers.IDStruct{Id: estado}, helpers.ObtenerTerceroInscripcion(SolicitudInscripcion))
				errorCambioEstado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_evolucion_estado", "POST", &respCambio, inscripcionEvolucionEstado)

				errParam := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=Activo:true,ParametroId.TipoParametroId.Id:2,ParametroId.CodigoAbreviacion:"+TipoParametro+",PeriodoId.Year:"+fmt.Sprintf("%v", objTransaccion["aniopago"])+",PeriodoId.CodigoAbreviacion:VG", &parametro)
				if errorCambioEstado == nil && respCambio["Status"] != "400" && errParam == nil && fmt.Sprintf("%v", parametro["Data"].([]interface{})[0]) != "map[]" {
					Dato := parametro["Data"].([]interface{})[0]
					if errJson := json.Unmarshal([]byte(Dato.(map[string]interface{})["Valor"].(string)), &Valor); errJson == nil {
						objTransaccion["valorordinario"] = Valor["Costo"].(float64)
						objTransaccion["valorextraordinario"] = Valor["Costo"].(float64)
						//objTransaccion["tiporecibo"] = Dato.(map[string]interface{})["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)
						objTransaccion["concepto"] = Dato.(map[string]interface{})["ParametroId"].(map[string]interface{})["Nombre"].(string)

						SolicitudRecibo := objTransaccion

						reciboSolicitud := httplib.Post("http://" + beego.AppConfig.String("GenerarReciboJbpmService") + "recibos_pago_proxy")
						reciboSolicitud.Header("Accept", "application/json")
						reciboSolicitud.Header("Content-Type", "application/json")
						reciboSolicitud.JSONBody(SolicitudRecibo)
						//errRecibo := request.SendJson("http://"+beego.AppConfig.String("GenerarReciboJbpmService")+"recibosPagoProxy", "POST", &NuevoRecibo, SolicitudRecibo)
						//fmt.Println("http://" + beego.AppConfig.String("GenerarReciboJbpmService") + "recibosPagoProxy")

						if errRecibo := reciboSolicitud.ToJSON(&NuevoRecibo); errRecibo == nil {
							inscripcionRealizada["ReciboInscripcion"] = fmt.Sprintf("%v/%v", NuevoRecibo["creaTransaccionResponse"].(map[string]interface{})["secuencia"], NuevoRecibo["creaTransaccionResponse"].(map[string]interface{})["anio"])
							var inscripcionUpdate map[string]interface{}
							errInscripcionUpdate := request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "PUT", &inscripcionUpdate, inscripcionRealizada)
							if errInscripcionUpdate == nil {
								APIResponseDTO = requestresponse.APIResponseDTO(true, 200, inscripcionUpdate, nil)

								fecha_actual := time.Now()
								dataEmail := map[string]interface{}{
									"dia":    fecha_actual.Day(),
									"mes":    utils.GetNombreMes(fecha_actual.Month()),
									"anio":   fecha_actual.Year(),
									"nombre": SolicitudInscripcion["Nombre"].(string) + " " + SolicitudInscripcion["Apellido"].(string),
									"estado": "inscripción solicitada",
								}
								fmt.Println(dataEmail)
								//utils.SendNotificationInscripcionSolicitud(dataEmail, objTransaccion["correo"].(string))
							} else {
								logs.Error(errInscripcionUpdate)
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInscripcionUpdate.Error())
							}
						} else {
							//var resDelete string
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
							helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
							logs.Error(errRecibo)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errRecibo.Error())
						}
					} else {
						//var resDelete string
						//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
						helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
						logs.Error(errJson)
						APIResponseDTO = requestresponse.APIResponseDTO(false, 403, nil, errJson.Error())
					}
				} else {
					//var resDelete string
					//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
					helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
					logs.Error(errParam)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errParam.Error())
				}

			} else {
				logs.Error(errInscripcion)
				APIResponseDTO = requestresponse.APIResponseDTO(true, 204, nil, errInscripcion.Error())
			}
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 409, nil, "Los códigos del proyecto curricular no coinciden, favor comunicarse con el administrador de sistema.")
		}
		// } else {
		// 	APIResponseDTO = requestresponse.APIResponseDTO(true, 204, nil, "Recipe already exist")
		// }

		// } else if err == "400" {
		// 	APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Bad request")
		// } else {
		// 	APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
		// }

	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 403, nil, err.Error())
	}

	return APIResponseDTO
}

func ActualizarEstadoMatriculado(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var actualizacionEstadoRequest map[string]interface{}
	var errorGetAll bool = false
	var resultado []map[string]interface{}

	if err := json.Unmarshal(data, &actualizacionEstadoRequest); err == nil {
		idTercero := actualizacionEstadoRequest["personaId"].(float64)
		idPeriodo := actualizacionEstadoRequest["periodoId"].(float64)
		idProyecto := actualizacionEstadoRequest["proyectoId"].(float64)
		var resultadoInscripcion []map[string]interface{}

		errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=Activo:true,PeriodoId:"+fmt.Sprintf("%.f", idPeriodo)+",PersonaId:"+fmt.Sprintf("%.f", idTercero)+"&sortby=Id&order=asc", &resultadoInscripcion)
		if errInscripcion == nil && fmt.Sprintf("%v", resultadoInscripcion[0]["System"]) != "map[]" {
			if resultadoInscripcion[0]["Status"] != 404 && resultadoInscripcion[0]["Id"] != nil {
				for _, inscripcion := range resultadoInscripcion {
					id := inscripcion["Id"].(float64)

					if tipoInscripcion, ok := inscripcion["TipoInscripcionId"].(map[string]interface{}); ok {
						if inscripcion["ProgramaAcademicoId"] == idProyecto {
							infoInscripcion := GenerarCuerpoActualizacionEstadoInscripcion(11, inscripcion, tipoInscripcion)
							var cuposInscripcion []map[string]interface{}

							if resInsc, errInsc := ActualizarInscripcion(infoInscripcion, id); errInsc == nil {
								resultado = append(resultado, resInsc)

								if resInscripcion, err := RecuperarProgramasPeriodoProyecto(idPeriodo, idProyecto); err == nil {
									cuposInscripcion = resInscripcion
								} else {
									errorGetAll = true
									APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
									return APIResponseDTO
								}

								cuposHabilitadosPrograma := cuposInscripcion[0]["CuposHabilitados"].(float64) - 1
								if tipoInscripcion, ok := cuposInscripcion[0]["TipoInscripcionId"].(map[string]interface{}); ok {
									id := cuposInscripcion[0]["Id"].(float64)
									infoCupoInscripcion := GenerarCuerpoActualizarCupoInscripcion(cuposInscripcion[0], cuposHabilitadosPrograma, tipoInscripcion)

									if resCupoInsc, errCupoInsc := ActualizarCupoInscripcion(infoCupoInscripcion, id); errCupoInsc == nil {
										resultado = append(resultado, resCupoInsc)
									} else {
										errorGetAll = true
										APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errCupoInsc)
									}
								}
							} else {
								errorGetAll = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInsc)
							}
						} else {
							infoInscripcion := GenerarCuerpoActualizacionEstadoInscripcion(4, inscripcion, tipoInscripcion)

							if resInsc, errInsc := ActualizarInscripcion(infoInscripcion, id); errInsc == nil {
								resultado = append(resultado, resInsc)
							} else {
								errorGetAll = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInsc)
							}
						}
					} else {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource"))
					}
				}
			} else {
				if resultadoInscripcion[0]["Message"] == "Not found resource" {
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource"))
				} else {
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource"))
				}
			}
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInscripcion)
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}
	return APIResponseDTO
}

func ActualizarInscripcion(infoComp map[string]interface{}, id float64) (map[string]interface{}, error) {
	var resp map[string]interface{}

	cambiaEstado := false
	nuevoEstado := helpers.GetEstadoInscripcion(infoComp)
	var estadoActual *int

	var resultadoInscripcion map[string]interface{}
	errorConsulta := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%.f", id), &resultadoInscripcion)
	if errorConsulta == nil {
		estadoActual = helpers.GetEstadoInscripcion(resultadoInscripcion)
		cambiaEstado = estadoActual != nuevoEstado
	}

	errPutInfoComp := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%.f", id), "PUT", &resp, infoComp)
	if errPutInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
		if cambiaEstado {
			var respCambio map[string]interface{}
			inscripcionEvolucionEstado := helpers.GenerarInscripcionEvolucionEstado(int(infoComp["Id"].(float64)), &helpers.IDStruct{Id: estadoActual}, helpers.IDStruct{Id: nuevoEstado}, helpers.ObtenerTerceroInscripcion(infoComp))
			errorCambioEstado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_evolucion_estado", "POST", &respCambio, inscripcionEvolucionEstado)
			if errorCambioEstado != nil {
				return resp, errPutInfoComp
			}
		}
		return resp, nil
	} else {
		return resp, errPutInfoComp
	}
}

func ActualizarEstadoInscripcion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var inscripcion map[string]interface{}
	if err := json.Unmarshal(data, &inscripcion); err == nil {
		if resInsc, errInsc := ActualizarInscripcion(inscripcion, inscripcion["Id"].(float64)); errInsc == nil {
			APIResponseDTO = requestresponse.APIResponseDTO(true, 400, resInsc, nil)
		} else {
			logs.Error(errInsc)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInsc.Error())
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	return APIResponseDTO
}

func ActualizarCupoInscripcion(infoComp map[string]interface{}, id float64) (map[string]interface{}, error) {
	var resp map[string]interface{}
	errPutInfoComp := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"cupo_inscripcion/"+fmt.Sprintf("%.f", id), "PUT", &resp, infoComp)
	if errPutInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
		return resp, nil
	} else {
		return resp, errPutInfoComp
	}
}

func GenerarCuerpoActualizacionEstadoInscripcion(nuevoEstado int, inscripcion map[string]interface{}, tipoInscripcion map[string]interface{}) map[string]interface{} {
	InfoEstadoInscripcionId := map[string]interface{}{
		"Id": nuevoEstado,
	}
	InfoTipoInscripcionId := map[string]interface{}{
		"Id": tipoInscripcion["Id"],
	}
	bodyInscripcion := map[string]interface{}{
		"PersonaId":           inscripcion["PersonaId"],
		"ProgramaAcademicoId": inscripcion["ProgramaAcademicoId"],
		"ReciboInscripcion":   inscripcion["ReciboInscripcion"],
		"PeriodoId":           inscripcion["PeriodoId"],
		"EnfasisId":           inscripcion["EnfasisId"],
		"AceptaTerminos":      inscripcion["AceptaTerminos"],
		"FechaAceptaTerminos": inscripcion["FechaAceptaTerminos"],
		"Activo":              true,
		"EstadoInscripcionId": InfoEstadoInscripcionId,
		"TipoInscripcionId":   InfoTipoInscripcionId,
		"NotaFinal":           inscripcion["NotaFinal"],
		"Credencial":          inscripcion["Credencial"],
		"Opcion":              inscripcion["Opcion"],
	}
	return bodyInscripcion
}

func GenerarCuerpoActualizarCupoInscripcion(CupoInscripcion map[string]interface{}, cuposHabilitados float64, tipoInscripcion map[string]interface{}) map[string]interface{} {
	InfoTipoInscripcionId := map[string]interface{}{
		"Id": tipoInscripcion["Id"],
	}
	bodyCupoInscripcion := map[string]interface{}{
		"Activo":              CupoInscripcion["Activo"],
		"CuposHabilitados":    cuposHabilitados,
		"CuposOpcionados":     CupoInscripcion["CuposOpcionados"],
		"PeriodoId":           CupoInscripcion["PeriodoId"],
		"ProgramaAcademicoId": CupoInscripcion["ProgramaAcademicoId"],
		"TipoInscripcionId":   InfoTipoInscripcionId,
		"CupoId":              CupoInscripcion["CupoId"],
	}
	return bodyCupoInscripcion
}

func ActualizarCupos(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var actualizacionCuposRequest map[string]interface{}
	var errorGetAll bool = false
	var resultado []map[string]interface{}

	if err := json.Unmarshal(data, &actualizacionCuposRequest); err == nil {
		idPeriodo := actualizacionCuposRequest["periodoId"].(float64)
		var fechaActualCiclo int
		var cicloActual string
		var cuposInscripcion []map[string]interface{}

		// RECUPERACIÓN DEL CICLO Y FECHA DEL CICLO ACTUAL
		if fecha, ciclo, err := RecuperarFechaCicloActual(idPeriodo); err == nil {
			fechaActualCiclo = fecha
			cicloActual = ciclo
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
			return APIResponseDTO
		}

		// RECUPERACIÓN DE LOS PROGRAMAS ACADEMICOS QUE TIENEN CUPOS EN CIERTO PERIODO
		if resInscripcion, err := RecuperarProgramasPeriodo(idPeriodo); err == nil {
			cuposInscripcion = resInscripcion
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
			return APIResponseDTO
		}

		for _, CupoInscripcion := range cuposInscripcion {
			idprograma := CupoInscripcion["ProgramaAcademicoId"].(float64)
			cuposHabilitadosPrograma := CupoInscripcion["CuposHabilitados"].(float64)
			var resultadoInscripcion []map[string]interface{}

			errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=Activo:true,PeriodoId:"+fmt.Sprintf("%.f", idPeriodo)+",ProgramaAcademicoId:"+fmt.Sprintf("%.f", idprograma)+",Opcion:"+cicloActual+"&sortby=Id&order=asc", &resultadoInscripcion)
			if errInscripcion == nil && fmt.Sprintf("%v", resultadoInscripcion[0]["System"]) != "map[]" {
				if resultadoInscripcion[0]["Status"] != 404 && resultadoInscripcion[0]["Id"] != nil {
					for _, inscripcion := range resultadoInscripcion {
						if estadoInscripcion, ok := inscripcion["EstadoInscripcionId"].(map[string]interface{}); ok {
							idEstadoInscripcion := estadoInscripcion["Id"].(float64)

							if idEstadoInscripcion == 2 || idEstadoInscripcion == 8 || idEstadoInscripcion == 10 {
								if tipoInscripcion, ok := inscripcion["TipoInscripcionId"].(map[string]interface{}); ok {
									id := inscripcion["Id"].(float64)
									infoInscripcion := GenerarCuerpoActualizacionEstadoInscripcion(4, inscripcion, tipoInscripcion)

									if resInsc, errInsc := ActualizarInscripcion(infoInscripcion, id); errInsc == nil {
										resultado = append(resultado, resInsc)
									} else {
										errorGetAll = true
										APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInsc)
									}
								}
							}
						}
					}

					// EN CASO DE SER LA PRIMERA FECHA SE CAMBIAN LOS OPCIONADOS A ADMITIDOS, SI QUEDAN CUPOS DISPONIBLES
					if fechaActualCiclo == 1 {
						var inscripcionesOpcionados []map[string]interface{}
						if resInscripcion, err := RecuperarInscripcionesOpcionadas(idPeriodo, idprograma); err == nil {
							inscripcionesOpcionados = resInscripcion
						} else {
							errorGetAll = true
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
							return APIResponseDTO
						}

						for _, inscripcionOpcionada := range inscripcionesOpcionados {
							if cuposHabilitadosPrograma == 0 {
								break
							}
							opcion := inscripcionOpcionada["Opcion"].(float64)
							if ciclo, err := strconv.ParseFloat(cicloActual, 64); err == nil {
								if opcion > ciclo {
									continue
								}
							} else {
								errorGetAll = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
								return APIResponseDTO
							}

							if tipoInscripcion, ok := inscripcionOpcionada["TipoInscripcionId"].(map[string]interface{}); ok {
								id := inscripcionOpcionada["Id"].(float64)
								infoInscripcion := GenerarCuerpoActualizacionEstadoInscripcion(2, inscripcionOpcionada, tipoInscripcion)

								if resInsc, errInsc := ActualizarInscripcion(infoInscripcion, id); errInsc == nil {
									resultado = append(resultado, resInsc)
									cuposHabilitadosPrograma = cuposHabilitadosPrograma - 1
								} else {
									errorGetAll = true
									APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInsc)
								}
							}
						}
					}
				} else {
					if resultadoInscripcion[0]["Message"] == "Not found resource" {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource"))
					} else {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource"))
					}
				}
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInscripcion)
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}
	return APIResponseDTO
}

// Función encargada de recuperar la fecha y el ciclo actual del proceso de admisión
func RecuperarFechaCicloActual(idPeriodo float64) (int, string, error) {
	var resultadoCalendario []map[string]interface{}
	var fechaActualCiclo int
	var cicloActual string
	fechaActual := time.Now()

	errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Nivel:1,Activo:true,PeriodoId:"+fmt.Sprintf("%.f", idPeriodo)+"&sortby=Id&order=asc", &resultadoCalendario)
	if errCalendario == nil && fmt.Sprintf("%v", resultadoCalendario[0]["System"]) != "map[]" {
		if resultadoCalendario[0]["Status"] != 404 && resultadoCalendario[0]["Id"] != nil {
			idCalendario := resultadoCalendario[0]["Id"].(float64)
			var resultadoTipoEvento []map[string]interface{}

			errTipoEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"tipo_evento?query=CodigoAbreviacion:CIAD,Activo:true,CalendarioID.Id:"+fmt.Sprintf("%.f", idCalendario)+"&sortby=Id&order=asc", &resultadoTipoEvento)
			if errTipoEvento == nil && fmt.Sprintf("%v", resultadoTipoEvento[0]["System"]) != "map[]" {
				if resultadoTipoEvento[0]["Status"] != 404 && resultadoTipoEvento[0]["Id"] != nil {
					idTipoEventoCiclos := resultadoTipoEvento[0]["Id"].(float64)
					var resultadoCalendarioEvento []map[string]interface{}

					errCalendarioEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=Activo:true,TipoEventoId.Id:"+fmt.Sprintf("%.f", idTipoEventoCiclos)+"&sortby=Id&order=asc", &resultadoCalendarioEvento)
					if errCalendarioEvento == nil && fmt.Sprintf("%v", resultadoCalendarioEvento[0]["System"]) != "map[]" {
						if resultadoCalendarioEvento[0]["Status"] != 404 && resultadoCalendarioEvento[0]["Id"] != nil {
							for _, evento := range resultadoCalendarioEvento {
								fechaInicioStr, ok1 := evento["FechaInicio"].(string)
								fechaFinStr, ok2 := evento["FechaFin"].(string)

								if ok1 && ok2 {
									fechaInicio, err1 := time.Parse(time.RFC3339, fechaInicioStr)
									fechaFin, err2 := time.Parse(time.RFC3339, fechaFinStr)
									if err1 == nil && err2 == nil {
										if fechaActual.After(fechaInicio) && fechaActual.Before(fechaFin) {
											if descripcion, ok := evento["Descripcion"].(string); ok {
												// AQUI SE SETEA EL CICLO ACTUAL
												cicloActual = descripcion
												idCicloActual := evento["Id"].(float64)
												var resultadoFechasCiclo []map[string]interface{}

												errFechasCiclo := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=Activo:true,EventoPadreId.Id:"+fmt.Sprintf("%.f", idCicloActual)+"&sortby=Id&order=asc", &resultadoFechasCiclo)
												if errFechasCiclo == nil && fmt.Sprintf("%v", resultadoFechasCiclo[0]["System"]) != "map[]" {
													if resultadoFechasCiclo[0]["Status"] != 404 && resultadoFechasCiclo[0]["Id"] != nil {
														for _, fecha := range resultadoFechasCiclo {
															fechaCicloInicioStr, okFecha1 := fecha["FechaInicio"].(string)
															fechaCicloFinStr, okFecha2 := fecha["FechaFin"].(string)

															if okFecha1 && okFecha2 {
																fechaInicioCiclo, errFecha1 := time.Parse(time.RFC3339, fechaCicloInicioStr)
																fechaFinCiclo, errFecha2 := time.Parse(time.RFC3339, fechaCicloFinStr)

																if errFecha1 == nil && errFecha2 == nil {
																	if fechaActual.After(fechaInicioCiclo) && fechaActual.Before(fechaFinCiclo) {
																		if tipoEvento, ok := fecha["TipoEventoId"].(map[string]interface{}); ok {
																			if codigoAbreviacion, ok := tipoEvento["CodigoAbreviacion"].(string); ok {
																				if strings.TrimSpace(codigoAbreviacion) == "PRAD" {
																					fechaActualCiclo = 1
																				} else if strings.TrimSpace(codigoAbreviacion) == "PROP" {
																					fechaActualCiclo = 2
																				}
																			} else {
																				return 0, "", fmt.Errorf("El valor de CodigoAbreviacion no es una cadena")
																			}
																		}
																	}
																} else {
																	return 0, "", fmt.Errorf("Error parsing dates")
																}
															} else {
																return 0, "", fmt.Errorf("Error: FechaInicio o FechaFin no es una cadena")
															}
														}
													} else {
														if resultadoFechasCiclo[0]["Message"] == "Not found resource" {
															return 0, "", fmt.Errorf("Not found resource")
														} else {
															return 0, "", fmt.Errorf("Not found resource")
														}
													}
												} else {
													return 0, "", errFechasCiclo
												}
											}
										}
									} else {
										return 0, "", fmt.Errorf("Error parsing dates")
									}
								} else {
									return 0, "", fmt.Errorf("Error: FechaInicio o FechaFin no es una cadena")
								}
							}
						} else {
							if resultadoCalendarioEvento[0]["Message"] == "Not found resource" {
								return 0, "", fmt.Errorf("Not found resource")
							} else {
								return 0, "", fmt.Errorf("Not found resource")
							}
						}
					} else {
						return 0, "", errCalendarioEvento
					}
				} else {
					if resultadoTipoEvento[0]["Message"] == "Not found resource" {
						return 0, "", fmt.Errorf("Not found resource")
					} else {
						return 0, "", fmt.Errorf("Not found resource")
					}
				}
			} else {
				return 0, "", errTipoEvento
			}
		} else {
			if resultadoCalendario[0]["Message"] == "Not found resource" {
				return 0, "", fmt.Errorf("Not found resource")
			} else {
				return 0, "", fmt.Errorf("Not found resource")
			}
		}
	} else {
		return 0, "", errCalendario
	}

	return fechaActualCiclo, cicloActual, nil
}

// Función encargada de recuperar las filas de la tabla cupo_inscripcion que esten relacionadas a un periodo
func RecuperarProgramasPeriodo(idPeriodo float64) ([]map[string]interface{}, error) {
	var resultadoCuposInscripcion []map[string]interface{}

	errCuposInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"cupo_inscripcion?query=Activo:true,PeriodoId:"+fmt.Sprintf("%.f", idPeriodo)+"&sortby=Id&order=asc", &resultadoCuposInscripcion)
	if errCuposInscripcion == nil && fmt.Sprintf("%v", resultadoCuposInscripcion[0]["System"]) != "map[]" {
		if resultadoCuposInscripcion[0]["Status"] != 404 && resultadoCuposInscripcion[0]["Id"] != nil {
			return resultadoCuposInscripcion, nil
		} else {
			if resultadoCuposInscripcion[0]["Message"] == "Not found resource" {
				return nil, fmt.Errorf("Not found resource")
			} else {
				return nil, fmt.Errorf("Not found resource")
			}
		}
	} else {
		return nil, errCuposInscripcion
	}
}

func RecuperarProgramasPeriodoProyecto(idPeriodo float64, idProgramaAcademico float64) ([]map[string]interface{}, error) {
	var resultadoCuposInscripcion []map[string]interface{}

	errCuposInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"cupo_inscripcion?query=Activo:true,ProgramaAcademicoId:"+fmt.Sprintf("%.f", idProgramaAcademico)+",PeriodoId:"+fmt.Sprintf("%.f", idPeriodo)+"&sortby=Id&order=asc", &resultadoCuposInscripcion)
	if errCuposInscripcion == nil && fmt.Sprintf("%v", resultadoCuposInscripcion[0]["System"]) != "map[]" {
		if resultadoCuposInscripcion[0]["Status"] != 404 && resultadoCuposInscripcion[0]["Id"] != nil {
			return resultadoCuposInscripcion, nil
		} else {
			if resultadoCuposInscripcion[0]["Message"] == "Not found resource" {
				return nil, fmt.Errorf("Not found resource")
			} else {
				return nil, fmt.Errorf("Not found resource")
			}
		}
	} else {
		return nil, errCuposInscripcion
	}
}

func RecuperarInscripcionesOpcionadas(idPeriodo float64, idprograma float64) ([]map[string]interface{}, error) {
	var resultadoInscripcion []map[string]interface{}

	errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=Activo:true,PeriodoId:"+fmt.Sprintf("%.f", idPeriodo)+",ProgramaAcademicoId:"+fmt.Sprintf("%.f", idprograma)+",EstadoInscripcionId.Id:3&sortby=Opcion&order=asc", &resultadoInscripcion)
	if errInscripcion == nil && fmt.Sprintf("%v", resultadoInscripcion[0]["System"]) != "map[]" {
		if resultadoInscripcion[0]["Status"] != 404 && resultadoInscripcion[0]["Id"] != nil {
			return resultadoInscripcion, nil
		} else {
			if resultadoInscripcion[0]["Message"] == "Not found resource" {
				return nil, fmt.Errorf("Not found resource")
				// errorGetAll = true
				// APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource"))
			} else {
				return nil, fmt.Errorf("Not found resource")
				// errorGetAll = true
				// APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource"))
			}
		}
	} else {
		return nil, errInscripcion
		// errorGetAll = true
		// APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errCuposInscripcion)
	}
}
