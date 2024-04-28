package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetInfoLegalizacionTercero(idTercero string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado consulta
	resultado := map[string]interface{}{}
	var errorGetAll bool

	// Recuperación de la dirección
	IdDirResidencia, _ := helpers.IdInfoCompTercero("9", "dir-residencia")
	var resultadoDirResidencia []map[string]interface{}
	errDirResidencia := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdDirResidencia+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoDirResidencia)
	if errDirResidencia == nil && fmt.Sprintf("%v", resultadoDirResidencia[0]["System"]) != "map[]" {
		if resultadoDirResidencia[0]["Status"] != 404 && resultadoDirResidencia[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoDirResidencia[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["direccionResidencia"] = nil
			} else {
				resultado["direccionResidencia"] = direccionJson["dato"]
			}
		} else {
			if resultadoDirResidencia[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDirResidencia)
	}

	// Recuperación del colegio donde se gradúo
	IdColGraduado, _ := helpers.IdInfoCompTercero("9", "col-g")
	var resultadoColGraduado []map[string]interface{}
	errColGraduado := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdColGraduado+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoColGraduado)
	if errColGraduado == nil && fmt.Sprintf("%v", resultadoColGraduado[0]["System"]) != "map[]" {
		if resultadoColGraduado[0]["Status"] != 404 && resultadoColGraduado[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoColGraduado[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["colegioGraduado"] = nil
			} else {
				resultado["colegioGraduado"] = direccionJson["dato"]
			}
		} else {
			if resultadoColGraduado[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errColGraduado)
	}

	// Recuperación del soporte del colegio donde se gradúo
	IdSoporteCol, _ := helpers.IdInfoCompTercero("9", "sop-col-g")
	var resultadoSoporteCol []map[string]interface{}
	errSoporteCol := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSoporteCol+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSoporteCol)
	if errSoporteCol == nil && fmt.Sprintf("%v", resultadoSoporteCol[0]["System"]) != "map[]" {
		if resultadoSoporteCol[0]["Status"] != 404 && resultadoSoporteCol[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSoporteCol[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["soporteColegio"] = nil
			} else {
				resultado["soporteColegio"] = direccionJson["dato"]
			}
		} else {
			if resultadoSoporteCol[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporteCol)
	}

	// Recuperación de la pensión mensual pagada en grado 11
	IdPension11, _ := helpers.IdInfoCompTercero("9", "pens-11")
	var resultadoPension11 []map[string]interface{}
	errPension11 := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdPension11+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoPension11)
	if errPension11 == nil && fmt.Sprintf("%v", resultadoPension11[0]["System"]) != "map[]" {
		if resultadoPension11[0]["Status"] != 404 && resultadoPension11[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoPension11[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["pensionGrado11"] = nil
			} else {
				resultado["pensionGrado11"] = direccionJson["dato"]
			}
		} else {
			if resultadoPension11[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errPension11)
	}

	// Recuperación del soporte de pensión mensual pagada en grado 11
	IdSoportePension11, _ := helpers.IdInfoCompTercero("9", "sop-pens-11")
	var resultadoSoportePension11 []map[string]interface{}
	errSoportePension11 := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSoportePension11+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSoportePension11)
	if errSoportePension11 == nil && fmt.Sprintf("%v", resultadoSoportePension11[0]["System"]) != "map[]" {
		if resultadoSoportePension11[0]["Status"] != 404 && resultadoSoportePension11[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSoportePension11[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["soportePensionGrado11"] = nil
			} else {
				resultado["soportePensionGrado11"] = direccionJson["dato"]
			}
		} else {
			if resultadoSoportePension11[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoportePension11)
	}

	// Recuperación del nucleo familiar
	IdNucleoFam, _ := helpers.IdInfoCompTercero("9", "nuc-f")
	var resultadoNucleoFam []map[string]interface{}
	errNucleoFam := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdNucleoFam+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoNucleoFam)
	if errNucleoFam == nil && fmt.Sprintf("%v", resultadoNucleoFam[0]["System"]) != "map[]" {
		if resultadoNucleoFam[0]["Status"] != 404 && resultadoNucleoFam[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoNucleoFam[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["nucleoFamiliar"] = nil
			} else {
				resultado["nucleoFamiliar"] = direccionJson["dato"]
			}
		} else {
			if resultadoNucleoFam[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errNucleoFam)
	}

	// Recuperación del soporte de nucleo familiar
	IdSoporteNucleoFam, _ := helpers.IdInfoCompTercero("9", "sop-nuc-f")
	var resultadoSoporteNucleoFam []map[string]interface{}
	errSoporteNucleoFam := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSoporteNucleoFam+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSoporteNucleoFam)
	if errSoporteNucleoFam == nil && fmt.Sprintf("%v", resultadoSoporteNucleoFam[0]["System"]) != "map[]" {
		if resultadoSoporteNucleoFam[0]["Status"] != 404 && resultadoSoporteNucleoFam[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSoporteNucleoFam[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["soporteNucleoFamiliar"] = nil
			} else {
				resultado["soporteNucleoFamiliar"] = direccionJson["dato"]
			}
		} else {
			if resultadoSoporteNucleoFam[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporteNucleoFam)
	}

	// Recuperación de la situación laboral
	IdSituacionLab, _ := helpers.IdInfoCompTercero("9", "sit-l")
	var resultadoSituacionLab []map[string]interface{}
	errSituacionLab := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSituacionLab+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSituacionLab)
	if errSituacionLab == nil && fmt.Sprintf("%v", resultadoSituacionLab[0]["System"]) != "map[]" {
		if resultadoSituacionLab[0]["Status"] != 404 && resultadoSituacionLab[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSituacionLab[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["situacionLaboral"] = nil
			} else {
				resultado["situacionLaboral"] = direccionJson["dato"]
			}
		} else {
			if resultadoSituacionLab[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSituacionLab)
	}

	// Recuperación del soporte de la situación laboral
	IdSoporteSituacionLab, _ := helpers.IdInfoCompTercero("9", "sop-sit-l")
	var resultadoSoporteSituacionLab []map[string]interface{}
	errSoporteSituacionLab := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSoporteSituacionLab+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSoporteSituacionLab)
	if errSoporteSituacionLab == nil && fmt.Sprintf("%v", resultadoSoporteSituacionLab[0]["System"]) != "map[]" {
		if resultadoSoporteSituacionLab[0]["Status"] != 404 && resultadoSoporteSituacionLab[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSoporteSituacionLab[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["soporteSituacionLaboral"] = nil
			} else {
				resultado["soporteSituacionLaboral"] = direccionJson["dato"]
			}
		} else {
			if resultadoSoporteSituacionLab[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporteSituacionLab)
	}

	// Recuperación de la dirección de residencia de la persona que costea
	IdDireccionCostea, _ := helpers.IdInfoCompTercero("1654", "Dir-rc")
	var resultadoDireccionCostea []map[string]interface{}
	errDireccionCostea := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdDireccionCostea+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoDireccionCostea)
	if errDireccionCostea == nil && fmt.Sprintf("%v", resultadoDireccionCostea[0]["System"]) != "map[]" {
		if resultadoDireccionCostea[0]["Status"] != 404 && resultadoDireccionCostea[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoDireccionCostea[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["direccionCostea"] = nil
			} else {
				resultado["direccionCostea"] = direccionJson["dato"]
			}
		} else {
			if resultadoDireccionCostea[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errDireccionCostea)
	}

	// Recuperación del soporte de estrato de residencia de la persona que costea
	IdSoporteEstratoCostea, _ := helpers.IdInfoCompTercero("1654", "sop-est-c")
	var resultadoSoporteEstratoCostea []map[string]interface{}
	errSoporteEstratoCostea := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSoporteEstratoCostea+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSoporteEstratoCostea)
	if errSoporteEstratoCostea == nil && fmt.Sprintf("%v", resultadoSoporteEstratoCostea[0]["System"]) != "map[]" {
		if resultadoSoporteEstratoCostea[0]["Status"] != 404 && resultadoSoporteEstratoCostea[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSoporteEstratoCostea[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["soporteEstratoCostea"] = nil
			} else {
				resultado["soporteEstratoCostea"] = direccionJson["dato"]
			}
		} else {
			if resultadoSoporteEstratoCostea[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporteEstratoCostea)
	}

	// Recuperación de los ingresos obtenidos el año anterior de la persona que costea
	IdIngresosCostea, _ := helpers.IdInfoCompTercero("1654", "ing-ac")
	var resultadoIngresosCostea []map[string]interface{}
	errIngresosCostea := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdIngresosCostea+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoIngresosCostea)
	if errIngresosCostea == nil && fmt.Sprintf("%v", resultadoIngresosCostea[0]["System"]) != "map[]" {
		if resultadoIngresosCostea[0]["Status"] != 404 && resultadoIngresosCostea[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoIngresosCostea[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["ingresosCostea"] = nil
			} else {
				resultado["ingresosCostea"] = direccionJson["dato"]
			}
		} else {
			if resultadoIngresosCostea[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errIngresosCostea)
	}

	// Recuperación del soporte de los ingresos obtenidos el año anterior de la persona que costea
	IdSoporteIngresosCostea, _ := helpers.IdInfoCompTercero("1654", "sop-ing-ac")
	var resultadoSoporteIngresosCostea []map[string]interface{}
	errSoporteIngresosCostea := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSoporteIngresosCostea+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSoporteIngresosCostea)
	if errSoporteIngresosCostea == nil && fmt.Sprintf("%v", resultadoSoporteIngresosCostea[0]["System"]) != "map[]" {
		if resultadoSoporteIngresosCostea[0]["Status"] != 404 && resultadoSoporteIngresosCostea[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSoporteIngresosCostea[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["soporteIngresosCostea"] = nil
			} else {
				resultado["soporteIngresosCostea"] = direccionJson["dato"]
			}
		} else {
			if resultadoSoporteIngresosCostea[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporteIngresosCostea)
	}

	// Recuperación del soporte general
	IdSoporteGeneral, _ := helpers.IdInfoCompTercero("9", "sop-gral")
	var resultadoSoporteGeneral []map[string]interface{}
	errSoporteGeneral := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdSoporteGeneral+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoSoporteGeneral)
	if errSoporteGeneral == nil && fmt.Sprintf("%v", resultadoSoporteGeneral[0]["System"]) != "map[]" {
		if resultadoSoporteGeneral[0]["Status"] != 404 && resultadoSoporteGeneral[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoSoporteGeneral[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["soporteGeneral"] = nil
			} else {
				resultado["soporteGeneral"] = direccionJson["dato"]
			}
		} else {
			if resultadoSoporteGeneral[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errSoporteGeneral)
	}

	// Retorno
	if !errorGetAll {
		return requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}
}
