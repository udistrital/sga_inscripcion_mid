package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_inscripcion_mid/helpers"
	"github.com/udistrital/utils_oas/requestresponse"
)

func PostAlertSolicitud(idTercero string, idTipoProduccionSrt string, data []byte) (APIResponseDTO requestresponse.APIResponse){
	idTipoProduccion, _ := strconv.Atoi(idTipoProduccionSrt)

	//resultado experiencia
	resultado := make(map[string]interface{})
	var SolicitudProduccion map[string]interface{}
	fmt.Println("Post Alert Solicitud")
	fmt.Println("Id Tercero: ", idTercero)
	fmt.Println("Id Tercero: ", idTipoProduccionSrt)

	if err := json.Unmarshal(data, &SolicitudProduccion); err == nil {
		if SolicitudProduccionAlert, errAlert := helpers.CheckCriteriaData(SolicitudProduccion, idTipoProduccion, idTercero); errAlert == nil {
			if SolicitudProduccionPut, errCoincidence := helpers.CheckCoincidenceProduction(SolicitudProduccionAlert, idTipoProduccion, idTercero); errCoincidence == nil {
				idStr := fmt.Sprintf("%v", SolicitudProduccionPut["Id"])
				fmt.Println(idStr)
				if resultadoPutSolicitudDocente, errPut := helpers.PutSolicitudDocente(SolicitudProduccionPut, idStr); errPut == nil {
					resultado = resultadoPutSolicitudDocente
					APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
					} else {
					logs.Error(errPut)
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPut)
					return APIResponseDTO
				}
			} else {
				logs.Error(errCoincidence)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errCoincidence)
				return APIResponseDTO
			}
		} else {
			logs.Error(errAlert)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errAlert)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func PutResultado(idProduccion string, data []byte)(APIResponseDTO requestresponse.APIResponse) {
	var SolicitudProduccion map[string]interface{}
	if err := json.Unmarshal(data, &SolicitudProduccion); err == nil {
		if SolicitudProduccionResult, errPuntaje := helpers.GenerateResult(SolicitudProduccion); errPuntaje == nil {
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, SolicitudProduccionResult, nil)
		} else {
			logs.Error(SolicitudProduccionResult)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errPuntaje)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func PostSolicitudEvaluacion(idSolicitud string, idSolicitudCoincidencia string, idTercero string, data []byte)(APIResponseDTO requestresponse.APIResponse){
	//resultado experiencia
	resultado := make(map[string]interface{})
	var SolicitudProduccion map[string]interface{}

	if err := json.Unmarshal(data, &SolicitudProduccion); err == nil {
		if SolicitudProduccionClone, errClone := helpers.GenerateEvaluationsCloning(SolicitudProduccion, idSolicitud, idSolicitudCoincidencia, idTercero); errClone == nil {
			if len(SolicitudProduccionClone) > 0 {
				resultado = SolicitudProduccion
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
					return APIResponseDTO
			}
		} else {
			logs.Error(errClone)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errClone)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}
	return APIResponseDTO
}

