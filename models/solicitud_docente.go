package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
	"strconv"
	"time"
)

// PostSolicitudDocente is ...
func PostSolicitudDocente(SolicitudDocente map[string]interface{}) (result map[string]interface{}, outputError interface{}) {
	date := time_bogota.TiempoBogotaFormato()
	var resultado map[string]interface{}

	SolicitudDocentePost := make(map[string]interface{})
	//formatdata.JsonPrint(SolicitudDocente)
	SolicitudDocentePost["Solicitud"] = map[string]interface{}{
		"Resultado":             SolicitudDocente["Resultado"],
		"Referencia":            SolicitudDocente["Referencia"],
		"FechaRadicacion":       date,
		"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
		"Activo":                true,
		"FechaCreacion":         date,
		"FechaModificacion":     date,
		"SolicitudPadreId":      SolicitudDocente["SolicitudPadreId"],
	}

	var terceroID interface{}
	var solicitantes []map[string]interface{}
	for _, solicitanteTemp := range SolicitudDocente["Autores"].([]interface{}) {
		solicitante := solicitanteTemp.(map[string]interface{})
		terceroID = solicitante["Persona"]
		solicitantes = append(solicitantes, map[string]interface{}{
			"TerceroId":         solicitante["Persona"],
			"SolicitudId":       map[string]interface{}{"Id": 0},
			"Activo":            true,
			"FechaCreacion":     date,
			"FechaModificacion": date,
		})
	}
	if len(solicitantes) == 0 {
		solicitantes = append(solicitantes, map[string]interface{}{})
	}
	SolicitudDocentePost["Solicitantes"] = solicitantes

	if terceroID == nil {
		terceroID = SolicitudDocente["TerceroId"]
	}

	fmt.Println("TerceroID: ")
	fmt.Println(terceroID)

	var solicitudesEvolucionEstado []map[string]interface{}
	solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{
		"TerceroId":             terceroID,
		"SolicitudId":           map[string]interface{}{"Id": 0},
		"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
		"FechaLimite":           CalcularFecha(SolicitudDocente["EstadoTipoSolicitudId"].(map[string]interface{})),
		"Activo":                true,
		"FechaCreacion":         date,
		"FechaModificacion":     date,
	})

	SolicitudDocentePost["EvolucionesEstado"] = solicitudesEvolucionEstado
	SolicitudDocentePost["Observaciones"] = nil

	fmt.Println("paso!")

	var resultadoSolicitudDocente map[string]interface{}
	errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud", "POST", &resultadoSolicitudDocente, SolicitudDocentePost)
	if errSolicitud == nil && fmt.Sprintf("%v", resultadoSolicitudDocente["System"]) != "map[]" && resultadoSolicitudDocente["Solicitud"] != nil {
		if resultadoSolicitudDocente["Status"] != 400 {
			resultado = resultadoSolicitudDocente
			return resultado, nil
		}
	} else {
		logs.Error(errSolicitud)
		return nil, errSolicitud
	}
	return resultado, nil
}

// PutSolicitudDocente is ...
func PutSolicitudDocente(SolicitudDocente map[string]interface{}, idStr string) (result map[string]interface{}, outputError interface{}) {
	date := time_bogota.TiempoBogotaFormato()
	//resultado experiencia
	var resultado map[string]interface{}
	SolicitudDocentePut := make(map[string]interface{})
	fechaRadicacion := time_bogota.TiempoCorreccionFormato(fmt.Sprintf("%v", SolicitudDocente["FechaRadicacion"]))
	yesterday, _ := strconv.Atoi(fmt.Sprintf("%v", SolicitudDocente["EstadoTipoSolicitudId"].(map[string]interface{})["Id"]))
	if yesterday == 1 {
		SolicitudDocentePut["Solicitud"] = map[string]interface{}{
			"Resultado":             SolicitudDocente["Resultado"],
			"Referencia":            SolicitudDocente["Referencia"],
			"FechaRadicacion":       date,
			"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
			"FechaModificacion":     date,
		}
	} else {
		SolicitudDocentePut["Solicitud"] = map[string]interface{}{
			"Resultado":             SolicitudDocente["Resultado"],
			"Referencia":            SolicitudDocente["Referencia"],
			"FechaRadicacion":       fechaRadicacion,
			"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
			"FechaModificacion":     date,
		}
	}
	var EstadoTipoSolicitudID interface{}
	for _, evolucionEstadoTemp := range SolicitudDocente["EvolucionEstado"].([]interface{}) {
		evolucionEstado := evolucionEstadoTemp.(map[string]interface{})
		EstadoTipoSolicitudID = evolucionEstado["EstadoTipoSolicitudId"]
	}

	var solicitudesEvolucionEstado []map[string]interface{}
	solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{
		"TerceroId":                     SolicitudDocente["TerceroId"],
		"SolicitudId":                   map[string]interface{}{"Id": 0},
		"EstadoTipoSolicitudId":         SolicitudDocente["EstadoTipoSolicitudId"],
		"EstadoTipoSolicitudIdAnterior": EstadoTipoSolicitudID,
		"Activo":                        true,
		"FechaLimite":                   CalcularFecha(SolicitudDocente["EstadoTipoSolicitudId"].(map[string]interface{})),
		"FechaCreacion":                 date,
		"FechaModificacion":             date,
	})

	var observaciones []map[string]interface{}
	for _, observacionTemp := range SolicitudDocente["Observaciones"].([]interface{}) {
		observacion := observacionTemp.(map[string]interface{})
		if observacion["Id"] == nil && observacion["Titulo"] != nil {
			observaciones = append(observaciones, map[string]interface{}{
				"TipoObservacionId": observacion["TipoObservacionId"],
				"SolicitudId":       map[string]interface{}{"Id": 0},
				"TerceroId":         observacion["TerceroId"],
				"Titulo":            observacion["Titulo"],
				"Valor":             observacion["Valor"],
				"FechaCreacion":     date,
				"FechaModificacion": date,
				"Activo":            true,
			})
		} else if observacion["Id"] != nil {
			observaciones = append(observaciones, map[string]interface{}{
				"Id":                observacion["Id"],
				"TipoObservacionId": observacion["TipoObservacionId"],
				"SolicitudId":       observacion["SolicitudId"],
				"TerceroId":         observacion["TerceroId"],
				"Titulo":            observacion["Titulo"],
				"Valor":             observacion["Valor"],
				"Activo":            true,
			})
		}
	}
	if len(observaciones) == 0 {
		observaciones = append(observaciones, map[string]interface{}{})
	}

	SolicitudDocentePut["Solicitantes"] = nil
	SolicitudDocentePut["EvolucionesEstado"] = solicitudesEvolucionEstado
	SolicitudDocentePut["Observaciones"] = observaciones

	var resultadoSolicitudDocente map[string]interface{}
	errSolicitudPut := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/"+idStr, "PUT", &resultadoSolicitudDocente, SolicitudDocentePut)
	if errSolicitudPut == nil && fmt.Sprintf("%v", resultadoSolicitudDocente["System"]) != "map[]" {
		if resultadoSolicitudDocente["Status"] != 400 {
			resultado = SolicitudDocente
			return resultado, nil
		}
	} else {
		logs.Error(errSolicitudPut)
		return nil, errSolicitudPut
	}
	return resultado, nil
}

// CalcularFecha is ...
func CalcularFecha(EstadoTipoSolicitud map[string]interface{}) (result string) {
	numDias, _ := strconv.Atoi(fmt.Sprintf("%v", EstadoTipoSolicitud["NumeroDias"]))
	var tiempoBogota time.Time
	tiempoBogota = time.Now()

	tiempoBogota = tiempoBogota.AddDate(0, 0, numDias)

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		fmt.Println(err)
	}
	tiempoBogota = tiempoBogota.In(loc)

	var tiempoBogotaStr = tiempoBogota.Format(time.RFC3339Nano)
	return tiempoBogotaStr
}
