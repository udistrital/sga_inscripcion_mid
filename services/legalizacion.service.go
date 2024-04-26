package services

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/astaxie/beego"
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
