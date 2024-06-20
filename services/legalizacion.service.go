package services

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/helpers"
	"github.com/udistrital/sga_inscripcion_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func CrearInfolegalizacionMatricula(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var legalizacionMatriculaRequest map[string]interface{}

	var algoFallo bool = false

	var respuestas []interface{}

	if err := json.Unmarshal(data, &legalizacionMatriculaRequest); err == nil {
		var infoSocioeconomicaRequest = map[string]interface{}{}
		var infoSocioeconomicaCosteaRequest = map[string]interface{}{}

		infoSocioeconomicaRequest["DireccionResidencia"] = legalizacionMatriculaRequest["DireccionResidencia"]
		infoSocioeconomicaRequest["Localidad"] = legalizacionMatriculaRequest["Localidad"]
		infoSocioeconomicaRequest["ColegioGraduado"] = legalizacionMatriculaRequest["ColegioGraduado"]
		infoSocioeconomicaRequest["PensionMensual11"] = legalizacionMatriculaRequest["PensionMensual11"]
		infoSocioeconomicaRequest["PensionMensualSM11"] = legalizacionMatriculaRequest["PensionMensualSM11"]
		infoSocioeconomicaRequest["NucleoFamiliar"] = legalizacionMatriculaRequest["NucleoFamiliar"]
		infoSocioeconomicaRequest["SituacionLaboral"] = legalizacionMatriculaRequest["SituacionLaboral"]
		infoSocioeconomicaRequest["SoporteDiploma"] = legalizacionMatriculaRequest["SoporteDiploma"]
		infoSocioeconomicaRequest["SoportePension"] = legalizacionMatriculaRequest["SoportePension"]
		infoSocioeconomicaRequest["SoporteNucleo"] = legalizacionMatriculaRequest["SoporteNucleo"]
		infoSocioeconomicaRequest["SoporteDocumental"] = legalizacionMatriculaRequest["SoporteDocumental"]

		infoSocioeconomicaCosteaRequest["DireccionResidenciaCostea"] = legalizacionMatriculaRequest["DireccionResidenciaCostea"]
		infoSocioeconomicaCosteaRequest["EstratoCostea"] = legalizacionMatriculaRequest["EstratoCostea"]
		infoSocioeconomicaCosteaRequest["UbicacionResidenciaCostea"] = legalizacionMatriculaRequest["UbicacionResidenciaCostea"]
		infoSocioeconomicaCosteaRequest["SoporteEstratoCostea"] = legalizacionMatriculaRequest["SoporteEstratoCostea"]
		infoSocioeconomicaCosteaRequest["IngresosCostea"] = legalizacionMatriculaRequest["IngresosCostea"]
		infoSocioeconomicaCosteaRequest["IngresosCosteaSM"] = legalizacionMatriculaRequest["IngresosCosteaSM"]
		infoSocioeconomicaCosteaRequest["SoporteIngresosCostea"] = legalizacionMatriculaRequest["SoporteIngresosCostea"]

		// FORMATEO ARCHIVOS INFO PERSONAL

		if newString, errMap := map2StringFieldStudyPlan(infoSocioeconomicaRequest, "SoporteDiploma"); errMap == nil {
			if newString != "" {
				infoSocioeconomicaRequest["SoporteDiploma"] = newString
			}
		}

		if newString, errMap := map2StringFieldStudyPlan(infoSocioeconomicaRequest, "SoportePension"); errMap == nil {
			if newString != "" {
				infoSocioeconomicaRequest["SoportePension"] = newString
			}
		}

		if newString, errMap := map2StringFieldStudyPlan(infoSocioeconomicaRequest, "SoporteNucleo"); errMap == nil {
			if newString != "" {
				infoSocioeconomicaRequest["SoporteNucleo"] = newString
			}
		}

		if newString, errMap := map2StringFieldStudyPlan(infoSocioeconomicaRequest, "SoporteDocumental"); errMap == nil {
			if newString != "" {
				infoSocioeconomicaRequest["SoporteDocumental"] = newString
			}
		}

		ok, value := exists("SoporteSituacionLaboral", legalizacionMatriculaRequest)
		if ok {
			infoSocioeconomicaRequest["SoporteSituacionLaboral"] = legalizacionMatriculaRequest["SoporteSituacionLaboral"]

			if newString, errMap := map2StringFieldStudyPlan(infoSocioeconomicaRequest, "SoporteSituacionLaboral"); errMap == nil {
				if newString != "" {
					infoSocioeconomicaRequest["SoporteSituacionLaboral"] = newString
				}
			}
			fmt.Println(value)
		}

		// FORMATEO ARCHIVOS INFO COSTEA

		if newString, errMap := map2StringFieldStudyPlan(infoSocioeconomicaCosteaRequest, "SoporteEstratoCostea"); errMap == nil {
			if newString != "" {
				infoSocioeconomicaCosteaRequest["SoporteEstratoCostea"] = newString
			}
		}

		if newString, errMap := map2StringFieldStudyPlan(infoSocioeconomicaCosteaRequest, "SoporteIngresosCostea"); errMap == nil {
			if newString != "" {
				infoSocioeconomicaCosteaRequest["SoporteIngresosCostea"] = newString
			}
		}

		if resLegalizacion, errLegalizacion := createLegalizacionMatricula(legalizacionMatriculaRequest["TerceroId"], infoSocioeconomicaRequest, infoSocioeconomicaCosteaRequest); errLegalizacion == nil {
			respuestas = append(respuestas, resLegalizacion)
		} else {
			algoFallo = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errLegalizacion.Error())
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !algoFallo {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, respuestas, nil)
	} else {
		return APIResponseDTO
	}
	return APIResponseDTO
}

func map2StringFieldStudyPlan(body map[string]any, fieldName string) (string, error) {
	if reflect.TypeOf(body[fieldName]).Kind() == reflect.Map {
		if stringNew, errMS := utils.Map2String(body[fieldName].(map[string]interface{})); errMS == nil {
			return stringNew, nil
		} else {
			return "", errMS
		}
	} else {
		return "", nil
	}
}

func exists(key string, m map[string]interface{}) (bool, interface{}) {
	val, ok := m[key]
	return ok, val
}

func createLegalizacionMatricula(terceroId interface{}, infoSocioeconomicaRequest map[string]interface{}, infoSocioeconomicaCosteaRequest map[string]interface{}) ([][]map[string]interface{}, error) {
	var infoLegalizacion [][]map[string]interface{}
	var errorInfoLegalizacion [][]map[string]interface{}

	if resInfoPer, errPlanPer := createInfoSocioEcomonomicaPersonal(infoSocioeconomicaRequest, terceroId); errPlanPer == nil {
		infoLegalizacion = append(infoLegalizacion, resInfoPer)
	} else {
		errorInfoLegalizacion = append(errorInfoLegalizacion, resInfoPer)
		return errorInfoLegalizacion, errPlanPer
	}

	if resInfoCos, errPlanCos := createInfoSocioEcomonomicaCostea(infoSocioeconomicaCosteaRequest, terceroId); errPlanCos == nil {
		infoLegalizacion = append(infoLegalizacion, resInfoCos)
	} else {
		errorInfoLegalizacion = append(errorInfoLegalizacion, resInfoCos)
		return errorInfoLegalizacion, errPlanCos
	}
	return infoLegalizacion, nil
}

func createInfoSocioEcomonomicaPersonal(infoSocioeconomicaBody map[string]interface{}, terceroId interface{}) ([]map[string]interface{}, error) {
	idsInfoComp := map[string]interface{}{
		"ColegioGraduado":         590,
		"DireccionResidencia":     588,
		"Localidad":               589,
		"PensionMensual11":        592,
		"PensionMensualSM11":      593,
		"NucleoFamiliar":          595,
		"SituacionLaboral":        597,
		"SoporteDiploma":          591,
		"SoportePension":          594,
		"SoporteNucleo":           596,
		"SoporteDocumental":       599,
		"SoporteSituacionLaboral": 598,
	}
	var allResInfoComp []map[string]interface{}
	var errorInfoComp []map[string]interface{}
	for key, value := range infoSocioeconomicaBody {
		TerceroId := map[string]interface{}{
			"Id": terceroId,
		}
		InfoComplementariaId := map[string]interface{}{
			"Id": idsInfoComp[key],
		}
		Dato := map[string]interface{}{
			"dato": value,
		}
		jsonDato, _ := json.Marshal(Dato)
		infoComp := map[string]interface{}{
			"TerceroId":            TerceroId,
			"InfoComplementariaId": InfoComplementariaId,
			"Activo":               true,
			"Dato":                 string(jsonDato),
		}

		if resInfoComp, errPlan := createInfoComplementaria(infoComp); errPlan == nil {
			allResInfoComp = append(allResInfoComp, resInfoComp)
		} else {
			errorInfoComp = append(errorInfoComp, resInfoComp)
			return errorInfoComp, errPlan
		}
	}
	return allResInfoComp, nil
}

func createInfoSocioEcomonomicaCostea(infoSocioeconomicaBody map[string]interface{}, terceroId interface{}) ([]map[string]interface{}, error) {
	idsInfoComp := map[string]interface{}{
		"DireccionResidenciaCostea": 600,
		"EstratoCostea":             601,
		"UbicacionResidenciaCostea": 602,
		"SoporteEstratoCostea":      603,
		"IngresosCostea":            604,
		"IngresosCosteaSM":          605,
		"SoporteIngresosCostea":     606,
	}
	var allResInfoComp []map[string]interface{}
	var errorInfoComp []map[string]interface{}
	for key, value := range infoSocioeconomicaBody {
		TerceroId := map[string]interface{}{
			"Id": terceroId,
		}
		InfoComplementariaId := map[string]interface{}{
			"Id": idsInfoComp[key],
		}
		Dato := map[string]interface{}{
			"dato": value,
		}
		jsonDato, _ := json.Marshal(Dato)
		infoComp := map[string]interface{}{
			"TerceroId":            TerceroId,
			"InfoComplementariaId": InfoComplementariaId,
			"Activo":               true,
			"Dato":                 string(jsonDato),
		}

		if resInfoComp, errPlan := createInfoComplementaria(infoComp); errPlan == nil {
			allResInfoComp = append(allResInfoComp, resInfoComp)
		} else {
			errorInfoComp = append(errorInfoComp, resInfoComp)
			return errorInfoComp, errPlan
		}
	}
	return allResInfoComp, nil
}

func createInfoComplementaria(infoCompBody map[string]interface{}) (map[string]interface{}, error) {
	var newInfoComp map[string]interface{}
	urlInfoComp := "http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero"

	if errNewPlan := request.SendJson(urlInfoComp, "POST", &newInfoComp, infoCompBody); errNewPlan == nil {
		return newInfoComp, nil
	} else {
		return newInfoComp, fmt.Errorf("TercerosService Error creando información complementaria de un tercero")
	}
}

func GetInfoLegalizacionTercero(idTercero string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado consulta
	resultado := map[string]interface{}{}
	var errorGetAll bool
	count := 0

	// Recuperación del estrato de la persona que costea
	IdEstratoCostea, _ := helpers.IdInfoCompTercero("1654", "est-c")
	var resultadoEstratoCostea []map[string]interface{}
	errEstratoCostea := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdEstratoCostea+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoEstratoCostea)
	if errEstratoCostea == nil && fmt.Sprintf("%v", resultadoEstratoCostea[0]["System"]) != "map[]" {
		if resultadoEstratoCostea[0]["Status"] != 404 && resultadoEstratoCostea[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoEstratoCostea[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["estratoCostea"] = nil
			} else {
				resultado["estratoCostea"] = direccionJson["dato"]
			}
		} else {
			if resultadoEstratoCostea[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				count += 1
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errEstratoCostea)
	}

	// Recuperación de la ubicación residencia de quien costea
	IdUbicacionResidenciaCostea, _ := helpers.IdInfoCompTercero("1654", "ubi-rc")
	var resultadoUbicacionResidenciaCostea []map[string]interface{}
	errUbicacionResidenciaCostea := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdUbicacionResidenciaCostea+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoUbicacionResidenciaCostea)
	if errUbicacionResidenciaCostea == nil && fmt.Sprintf("%v", resultadoUbicacionResidenciaCostea[0]["System"]) != "map[]" {
		if resultadoUbicacionResidenciaCostea[0]["Status"] != 404 && resultadoUbicacionResidenciaCostea[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoUbicacionResidenciaCostea[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["ubicacionResidenciaCostea"] = nil
			} else {
				resultado["ubicacionResidenciaCostea"] = direccionJson["dato"]
			}
		} else {
			if resultadoUbicacionResidenciaCostea[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				count += 1
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errUbicacionResidenciaCostea)
	}

	// Recuperación de los ingresos del año anterior en salarios minimos de quien costea
	IdIngresosSMCostea, _ := helpers.IdInfoCompTercero("1654", "ing-sm-ac")
	var resultadoIngresosSMCostea []map[string]interface{}
	errIngresosSMCostea := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdIngresosSMCostea+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoIngresosSMCostea)
	if errIngresosSMCostea == nil && fmt.Sprintf("%v", resultadoIngresosSMCostea[0]["System"]) != "map[]" {
		if resultadoIngresosSMCostea[0]["Status"] != 404 && resultadoIngresosSMCostea[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoIngresosSMCostea[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["ingresosSMCostea"] = nil
			} else {
				resultado["ingresosSMCostea"] = direccionJson["dato"]
			}
		} else {
			if resultadoIngresosSMCostea[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				count += 1
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errIngresosSMCostea)
	}

	// Recuperación de la localidad
	IdLocalidad, _ := helpers.IdInfoCompTercero("9", "localidad")
	var resultadoLocalidad []map[string]interface{}
	errLocalidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdLocalidad+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoLocalidad)
	if errLocalidad == nil && fmt.Sprintf("%v", resultadoLocalidad[0]["System"]) != "map[]" {
		if resultadoLocalidad[0]["Status"] != 404 && resultadoLocalidad[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoLocalidad[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["localidad"] = nil
			} else {
				resultado["localidad"] = direccionJson["dato"]
			}
		} else {
			if resultadoLocalidad[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				count += 1
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errLocalidad)
	}

	// Recuperación de la pensión pagada en grado 11 en salarios minimos
	IdPensionSM11, _ := helpers.IdInfoCompTercero("9", "pens-11-sm")
	var resultadoPensionSM11 []map[string]interface{}
	errPensionSM11 := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdPensionSM11+",TerceroId:"+idTercero+"&sortby=Id&order=desc&limit=1", &resultadoPensionSM11)
	if errPensionSM11 == nil && fmt.Sprintf("%v", resultadoPensionSM11[0]["System"]) != "map[]" {
		if resultadoPensionSM11[0]["Status"] != 404 && resultadoPensionSM11[0]["Id"] != nil {
			// unmarshall dato
			var direccionJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoPensionSM11[0]["Dato"].(string)), &direccionJson); err != nil {
				resultado["pensionSM11"] = nil
			} else {
				resultado["pensionSM11"] = direccionJson["dato"]
			}
		} else {
			if resultadoPensionSM11[0]["Message"] == "Not found resource" {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			} else {
				count += 1
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
			}
		}
	} else {
		errorGetAll = true
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errPensionSM11)
	}

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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
				count += 1
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
		if count == 18 {
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, "No existe legalizacion", nil)
		}
		return APIResponseDTO
	}
}

func ActualizarInfolegalizacionMatricula(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var legalizacionMatriculaRequest map[string]interface{}
	var errorGetAll bool = false
	var resultado []map[string]interface{}

	if err := json.Unmarshal(data, &legalizacionMatriculaRequest); err == nil {
		idTercero := legalizacionMatriculaRequest["personaId"].(float64)

		IdEstratoCostea, _ := helpers.IdInfoCompTercero("1654", "est-c")
		IdUbicacionResidenciaCostea, _ := helpers.IdInfoCompTercero("1654", "ubi-rc")
		IdIngresosSMCostea, _ := helpers.IdInfoCompTercero("1654", "ing-sm-ac")
		IdPensionSM11, _ := helpers.IdInfoCompTercero("9", "pens-11-sm")
		IdNucleoFam, _ := helpers.IdInfoCompTercero("9", "nuc-f")
		IdSituacionLab, _ := helpers.IdInfoCompTercero("9", "sit-l")
		IdIngresosCostea, _ := helpers.IdInfoCompTercero("1654", "ing-ac")
		IdPension11, _ := helpers.IdInfoCompTercero("9", "pens-11")

		idsInfoCompActualizar := map[string]interface{}{
			"estratoCostea":              IdEstratoCostea,
			"ubicacionResidencia":        IdUbicacionResidenciaCostea,
			"valorMatriculaUltimoAnio":   IdPension11,
			"valorMatriculaUltimoAnioSM": IdPensionSM11,
			"ingresosBrutosFam":          IdIngresosCostea,
			"ingresosBrutosFamSM":        IdIngresosSMCostea,
			"nucleoFamiliar":             IdNucleoFam,
			"situacionLaboral":           IdSituacionLab,
		}

		for key, value := range idsInfoCompActualizar {
			fmt.Println("Key:", key, "Value:", value)

			if stringValue, ok := value.(string); ok {
				if resInfoComp, errInfoComp := ObtenerInfoComplementariaActual(idTercero, stringValue, legalizacionMatriculaRequest, key); errInfoComp == nil {
					fmt.Println("ENTRADA PRINCIPAL OK:", resInfoComp)
					resultado = append(resultado, resInfoComp)
				} else {
					fmt.Println("ENTRADA PRINCIPAL MAL:", resInfoComp, errInfoComp)
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInfoComp.Error())
					return APIResponseDTO
				}
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}
	return APIResponseDTO
}

func ObtenerInfoComplementariaActual(idTercero float64, IdInfoComp string, legalizacionMatriculaRequest map[string]interface{}, stringInfoComp string) (map[string]interface{}, error) {
	var resultadoInfoComp []map[string]interface{}
	errInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdInfoComp+",TerceroId:"+fmt.Sprintf("%.f", idTercero)+"&sortby=Id&order=desc&limit=1", &resultadoInfoComp)
	if errInfoComp == nil && fmt.Sprintf("%v", resultadoInfoComp[0]["System"]) != "map[]" {
		if resultadoInfoComp[0]["Status"] != 404 && resultadoInfoComp[0]["Id"] != nil {
			id := resultadoInfoComp[0]["Id"].(float64)

			if infoComplementariaIdMap, ok := resultadoInfoComp[0]["InfoComplementariaId"].(map[string]interface{}); ok {

				TerceroId := map[string]interface{}{
					"Id": idTercero,
				}
				InfoComplementariaId := map[string]interface{}{
					"Id": infoComplementariaIdMap["Id"],
				}
				Dato := map[string]interface{}{
					"dato": legalizacionMatriculaRequest[stringInfoComp],
				}
				jsonDato, _ := json.Marshal(Dato)
				infoComp := map[string]interface{}{
					"TerceroId":            TerceroId,
					"InfoComplementariaId": InfoComplementariaId,
					"Activo":               true,
					"Dato":                 string(jsonDato),
				}

				if resInfoComp, errInfoComp := ActualizarInfoComplementaria(infoComp, id); errInfoComp == nil {
					return resInfoComp, nil
				} else {
					return resInfoComp, errInfoComp
				}
			} else {
				return infoComplementariaIdMap, fmt.Errorf("Not found resource in Info Complemenaria ID")
			}
		} else {
			if resultadoInfoComp[0]["Message"] == "Not found resource" {
				return resultadoInfoComp[0], fmt.Errorf("Not found resource")
			} else {
				return resultadoInfoComp[0], fmt.Errorf("Not found resource")
			}
		}
	} else {
		return resultadoInfoComp[0], errInfoComp
	}
}

func ActualizarInfoComplementaria(infoComp map[string]interface{}, id float64) (map[string]interface{}, error) {
	var resp map[string]interface{}
	errPutInfoComp := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", id), "PUT", &resp, infoComp)
	if errPutInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
		return resp, nil
	} else {
		return resp, errPutInfoComp
	}
}

func ActualizarInfoSoportelegalizacionMatricula(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var legalizacionMatriculaRequest map[string]interface{}
	var errorGetAll bool = false
	var resultado []map[string]interface{}

	if err := json.Unmarshal(data, &legalizacionMatriculaRequest); err == nil {
		idTercero := legalizacionMatriculaRequest["terceroId"].(float64)

		arvivosActualizar, ok := legalizacionMatriculaRequest["arvivosActualizar"].([]interface{})
		if ok {
			for _, item := range arvivosActualizar {
				archivo, ok := item.(map[string]interface{})
				if ok {
					idInfoComp := archivo["infoComplementariId"].(float64)
					dato := archivo["dato"]
					var resultadoInfoComp []map[string]interface{}
					fmt.Println("//**************************************//")
					fmt.Println("INFO COMPLEMENTARIA ID Y EL DATO")
					fmt.Println(idInfoComp, dato)

					if resInscripcion, err := recuperarInfoComplementariaByTerceroInfoCompId(idTercero, idInfoComp); err == nil {
						resultadoInfoComp = resInscripcion
					} else {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
						return APIResponseDTO
					}

					idInfo := resultadoInfoComp[0]["Id"].(float64)

					if infoComplementariaIdMap, ok := resultadoInfoComp[0]["InfoComplementariaId"].(map[string]interface{}); ok {

						TerceroId := map[string]interface{}{
							"Id": idTercero,
						}
						InfoComplementariaId := map[string]interface{}{
							"Id": infoComplementariaIdMap["Id"],
						}
						Dato := map[string]interface{}{
							"dato": dato,
						}
						jsonDato, _ := json.Marshal(Dato)
						infoComp := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaId,
							"Activo":               true,
							"Dato":                 string(jsonDato),
						}
						fmt.Println("INFO COMPLEMENTARIA ACTUALIZADA")
						fmt.Println(infoComp)

						if resInfoComp, errInfoComp := ActualizarInfoComplementaria(infoComp, idInfo); errInfoComp == nil {
							resultado = append(resultado, resInfoComp)
						} else {
							errorGetAll = true
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errInfoComp)
						}
					} else {
						errorGetAll = true
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Not found resource in Info Complemenaria ID"))
					}
				} else {
					errorGetAll = true
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Error al convertir el elemento a tipo map"))
				}
			}
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, fmt.Errorf("Archivos a actualizar esperaba ser tipo []interface{}"))
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
	} else {
		return APIResponseDTO
	}
	return APIResponseDTO
}

func recuperarInfoComplementariaByTerceroInfoCompId(idTercero float64, IdInfoComp float64) ([]map[string]interface{}, error) {
	var resultadoInfoComp []map[string]interface{}

	errInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+fmt.Sprintf("%.f", IdInfoComp)+",TerceroId:"+fmt.Sprintf("%.f", idTercero)+"&sortby=Id&order=desc&limit=1", &resultadoInfoComp)
	if errInfoComp == nil && fmt.Sprintf("%v", resultadoInfoComp[0]["System"]) != "map[]" {
		if resultadoInfoComp[0]["Status"] != 404 && resultadoInfoComp[0]["Id"] != nil {
			return resultadoInfoComp, nil
		} else {
			if resultadoInfoComp[0]["Message"] == "Not found resource" {
				return nil, fmt.Errorf("Not found resource")
			} else {
				return nil, fmt.Errorf("Not found resource")
			}
		}
	} else {
		return nil, errInfoComp
	}
}
