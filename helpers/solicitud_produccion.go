package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	//"github.com/agnivade/levenshtein"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// CheckCriteriaData is...
func CheckCriteriaData(SolicitudProduccion map[string]interface{}, idTipoProduccion int, idTercero string) (result map[string]interface{}, outputError interface{}) {
	var producciones []map[string]interface{}
	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			var ProduccionAcademica map[string]interface{}
			ProduccionAcademica = SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})
			var coincidences int
			var isbnCoincidences int
			var numRegisterCoincidences int
			var issnVolNumCoincidences int
			var eventCoincidences int
			var numAnnualProductions int
			var accumulatedPoints int
			var categoryLast int
			var isDurationAccepted bool
			var rangeAccepted int
			isDurationAccepted = true
			for _, produccion := range producciones {
				distance := CheckTitle(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion)
				if distance < 6 {
					coincidences++
				}

				if idTipoProduccion == 1 {
					if !checkLastChangeCategory(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion, idTercero) {
						categoryLast++
					}
				}
				if idTipoProduccion == 2 {
					accumulatedPoints += checkGradePoints(produccion, idTipoProduccion, idTercero)
				}
				if idTipoProduccion == 2 {
					if checkRageGrade(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
						rangeAccepted++
					}
				}
				if idTipoProduccion == 3 || idTipoProduccion == 4 || idTipoProduccion == 5 {
					if checkISSNVolNumber(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion) {
						issnVolNumCoincidences++
					}
				}
				if idTipoProduccion == 6 || idTipoProduccion == 7 || idTipoProduccion == 8 {
					if checkISBN(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion) {
						isbnCoincidences++
					}
				}
				if idTipoProduccion == 11 || idTipoProduccion == 12 {
					if checkRegisterNumber(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
						numRegisterCoincidences++
					}
				}
				if idTipoProduccion == 13 || idTipoProduccion == 14 {
					if checkEventName(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
						eventCoincidences++
					}
				}
				if idTipoProduccion >= 13 && idTipoProduccion != 18 {
					if checkAnnualProductionNumber(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
						numAnnualProductions++
					}
				}
			}
			if idTipoProduccion == 18 {
				isDurationAccepted = checkDurationPostDoctorado(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}))
			}
			coincidences--
			numAnnualProductions--
			isbnCoincidences--
			numRegisterCoincidences--
			issnVolNumCoincidences--
			eventCoincidences--
			isAccumulatedPass := checkMaxGradePoints(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), accumulatedPoints)
			generateAlerts(SolicitudProduccion, coincidences, numAnnualProductions, isAccumulatedPass, isbnCoincidences, numRegisterCoincidences, issnVolNumCoincidences, eventCoincidences, isDurationAccepted, rangeAccepted, categoryLast, idTipoProduccion)
			return SolicitudProduccion, nil
		}
	} else {
		logs.Error(producciones)
		return nil, errProduccion
	}
	return SolicitudProduccion, nil
}

// CheckTitle is ...
func CheckTitle(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}) (result int) {
	//TODO: Ajustar a como debe quedar este metodo
	//distance := levenshtein.ComputeDistance(fmt.Sprintf("%v", ProduccionAcademicaNew["Titulo"]), fmt.Sprintf("%v", ProduccionAcademicaRegister["Titulo"]))
	distance := 0
	return distance
}

func checkLastChangeCategory(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int, idTercero string) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	idSubTipoProduccionNewSrt := fmt.Sprintf("%v", ProduccionAcademicaNew["SubtipoProduccionId"].(map[string]interface{})["Id"])
	idSubTipoProduccionNew, _ := strconv.Atoi(idSubTipoProduccionNewSrt)

	idProduccionStr := fmt.Sprintf("%v", ProduccionAcademicaRegister["Id"])
	idProduccion, _ := strconv.Atoi(idProduccionStr)

	if idTipoProduccion == idTipoProduccionRegister {
		var solicitudes []map[string]interface{}
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/inactive/"+idTercero, &solicitudes)
		if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
			if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
				for _, solicitud := range solicitudes {
					type Reference struct{ Id int }
					var reference Reference
					json.Unmarshal([]byte(fmt.Sprintf("%v", solicitud["Referencia"])), &reference)
					if reference.Id == idProduccion {

						EvolucionEstadoList := solicitud["EvolucionEstado"].([]interface{})
						EvolucionEstado := EvolucionEstadoList[len(EvolucionEstadoList)-1].(map[string]interface{})

						dateNew, _ := time.Parse("2006-01-02", string([]rune(fmt.Sprintf("%v", ProduccionAcademicaNew["Fecha"]))[0:10]))
						dateRegister, _ := time.Parse("2006-01-02", string([]rune(fmt.Sprintf("%v", EvolucionEstado["FechaModificacion"]))[0:10]))
						resultDate := dateNew.Sub(dateRegister)

						if idSubTipoProduccionNew == 1 && (17532-resultDate.Hours()) > 0 {
							return false
						} else if idSubTipoProduccionNew == 2 && (26304-resultDate.Hours()) > 0 {
							return false
						} else if idSubTipoProduccionNew == 3 && (35064-resultDate.Hours()) > 0 {
							return false
						}
					}
				}
			}
		} else {
			return true
		}
	}
	return true
}

func checkAnnualProductionNumber(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	if idTipoProduccion != 16 {
		idSubTipoProduccionNewSrt := fmt.Sprintf("%v", ProduccionAcademicaNew["SubtipoProduccionId"].(map[string]interface{})["Id"])
		idSubTipoProduccionNew, _ := strconv.Atoi(idSubTipoProduccionNewSrt)
		idSubTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["Id"])
		idSubTipoProduccionRegister, _ := strconv.Atoi(idSubTipoProduccionRegisterSrt)
		if idSubTipoProduccionNew == idSubTipoProduccionRegister {
			yearNew := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaNew["Fecha"]))[0:4])
			yearRegister := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaRegister["Fecha"]))[0:4])
			if yearNew == yearRegister {
				return true
			}
		}
	} else {
		idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
		idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
		if idTipoProduccion == idTipoProduccionRegister {
			yearNew := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaNew["FechaCreacion"]))[0:4])
			yearRegister := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaRegister["Metadatos"].([]interface{})[0].(map[string]interface{})["FechaCreacion"]))[0:4])
			if yearNew == yearRegister {
				return true
			}
		}
	}
	return false
}

func checkRageGrade(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	if idTipoProduccionRegister == idTipoProduccion {
		idTipoProduccionNewSrt := fmt.Sprintf("%v", ProduccionAcademicaNew["SubtipoProduccionId"].(map[string]interface{})["Id"])
		idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["Id"])
		if idTipoProduccionRegisterSrt > idTipoProduccionNewSrt {
			return false
		}
	}
	return true
}

func checkEventName(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var eventNew string
	var eventRegister string
	var dateNew string
	var dateRegister string
	if idTipoProduccionRegister == idTipoProduccion {
		dateNew = string([]rune(fmt.Sprintf("%v", ProduccionAcademicaNew["ProduccionAcademica"].(map[string]interface{})["Fecha"]))[0:10])
		dateRegister = string([]rune(fmt.Sprintf("%v", ProduccionAcademicaRegister["Fecha"]))[0:10])
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 181 || tipoMetadatoID == 196 || tipoMetadatoID == 210 || tipoMetadatoID == 225 {
				eventNew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 181 || tipoMetadatoID == 196 || tipoMetadatoID == 210 || tipoMetadatoID == 225 {
				eventRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if eventNew == eventRegister && dateNew == dateRegister {
			return true
		}
	}
	return false
}

func checkISBN(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var ISBNnew string
	var ISBNregister string
	if idTipoProduccionRegister == 6 || idTipoProduccionRegister == 7 || idTipoProduccionRegister == 8 {
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 72 || tipoMetadatoID == 83 || tipoMetadatoID == 92 || tipoMetadatoID == 101 || tipoMetadatoID == 114 || tipoMetadatoID == 126 || tipoMetadatoID == 138 {
				ISBNnew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 72 || tipoMetadatoID == 83 || tipoMetadatoID == 92 || tipoMetadatoID == 101 || tipoMetadatoID == 114 || tipoMetadatoID == 126 || tipoMetadatoID == 138 {
				ISBNregister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if ISBNnew == ISBNregister {
			return true
		}
	}
	return false
}

func checkISSNVolNumber(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var ISSNnew string
	var ISSNregister string
	var volumeNew string
	var volumeRegister string
	var numberNew string
	var numberRegister string
	if idTipoProduccionRegister == 3 || idTipoProduccionRegister == 4 || idTipoProduccionRegister == 5 {
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 42 || tipoMetadatoID == 52 || tipoMetadatoID == 62 {
				ISSNnew = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 43 || tipoMetadatoID == 53 || tipoMetadatoID == 63 {
				volumeNew = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 46 || tipoMetadatoID == 56 || tipoMetadatoID == 66 {
				numberNew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 42 || tipoMetadatoID == 52 || tipoMetadatoID == 62 {
				ISSNregister = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 43 || tipoMetadatoID == 53 || tipoMetadatoID == 63 {
				volumeRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 46 || tipoMetadatoID == 56 || tipoMetadatoID == 66 {
				numberRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if ISSNnew == ISSNregister && volumeNew == volumeRegister && numberNew == numberRegister {
			return true
		}
	}
	return false
}

func checkRegisterNumber(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var nrNew string
	var nrRegister string
	if idTipoProduccionRegister == idTipoProduccion {
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 163 || tipoMetadatoID == 166 || tipoMetadatoID == 169 {
				nrNew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 163 || tipoMetadatoID == 166 || tipoMetadatoID == 169 {
				nrRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if nrNew == nrRegister {
			return true
		}
	}
	return false
}

func checkDurationPostDoctorado(ProduccionAcademicaNew map[string]interface{}) (result bool) {
	for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
		metadato := metadatoTemp.(map[string]interface{})
		metadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
		metadatoValor, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["Valor"]))
		if metadatoID == 257 && metadatoValor < 9 {
			return false
		}
	}
	return true
}

func checkGradePoints(ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int, idTercero string) (result int) {
	idProduccionStr := fmt.Sprintf("%v", ProduccionAcademicaRegister["Id"])
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idProduccion, _ := strconv.Atoi(idProduccionStr)
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var points int
	points = 0
	if idTipoProduccion == idTipoProduccionRegister {
		var solicitudes []map[string]interface{}
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/inactive/"+idTercero, &solicitudes)
		if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
			if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
				for _, solicitud := range solicitudes {
					type Reference struct{ Id int }
					var reference Reference
					json.Unmarshal([]byte(fmt.Sprintf("%v", solicitud["Referencia"])), &reference)
					if reference.Id == idProduccion && fmt.Sprintf("%v", solicitud["Resultado"]) != "" {
						type Result struct{ Puntaje int }
						var result Result
						json.Unmarshal([]byte(fmt.Sprintf("%v", solicitud["Resultado"])), &result)
						points += result.Puntaje
					}
				}
				return points
			}
		} else {
			return 0
		}
	}
	return 0
}

func checkMaxGradePoints(ProduccionAcademicaNew map[string]interface{}, accumulatedPoints int) (result bool) {
	idSubtipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaNew["SubtipoProduccionId"].(map[string]interface{})["Id"])
	idSubtipoProduccion, _ := strconv.Atoi(idSubtipoProduccionRegisterSrt)
	if idSubtipoProduccion == 4 && (140-accumulatedPoints) < 20 {
		return false
	} else if idSubtipoProduccion == 5 && (140-accumulatedPoints) < 40 {
		return false
	} else if idSubtipoProduccion == 6 && (140-accumulatedPoints) < 80 {
		return false
	}
	return true
}

func generateAlerts(SolicitudDocente map[string]interface{}, coincidences int, numAnnualProductions int, isAccumulatedPass bool, isbnCoincidences int, numRegisterCoincidences int, issnVolNumCoincidences int, eventCoincidences int, isDurationAccepted bool, rangeAccepted int, categoryLast int, idTipoProduccion int) {
	coincidencesSrt := strconv.Itoa(coincidences)
	var observaciones []interface{}
	var tipoObservacionData map[string]interface{}
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tipo_observacion/?query=Id:2", &tipoObservacionData)
	if errSolicitud == nil && fmt.Sprintf("%v", tipoObservacionData["System"]) != "map[]" {
		if tipoObservacionData["Status"] != 404 && tipoObservacionData["Data"] != nil {
			var tipoObservacion interface{}
			tipoObservacion = tipoObservacionData["Data"].([]interface{})[0]
			if coincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_numero_coincidencias" + coincidencesSrt,
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if eventCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_evento",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if categoryLast > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_ultima_categoria",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if isbnCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_isbn",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if issnVolNumCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_issn_volumen_numero",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if numRegisterCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_numero_registro",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if numAnnualProductions > 0 {
				switch idTipoProduccion {
				case 13, 14, 16, 17, 19:
					if numAnnualProductions > 5 {
						observaciones = append(observaciones, map[string]interface{}{
							"Titulo":            "alerta.titulo",
							"Valor":             "alerta.alerta_numero_produccion_anual_5",
							"TipoObservacionId": &tipoObservacion,
							"TerceroId":         0,
						})
					}
				case 15, 20:
					if numAnnualProductions > 3 {
						observaciones = append(observaciones, map[string]interface{}{
							"Titulo":            "alerta.titulo",
							"Valor":             "alerta.alerta_numero_produccion_anual_3",
							"TipoObservacionId": &tipoObservacion,
							"TerceroId":         0,
						})
					}
				default:
					fmt.Println("No entro a ninguno de los caso")
				}
			}
			if !isAccumulatedPass {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_puntos_grados",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if !isDurationAccepted {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_duracion",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if rangeAccepted > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_rango",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			SolicitudDocente["Observaciones"] = observaciones
		}
	}
}

// GenerateResult is ...
func GenerateResult(SolicitudProduccion map[string]interface{}) (result map[string]interface{}, outputError interface{}) {
	produccionAcademica := SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})
	subTipoProduccionID := produccionAcademica["SubtipoProduccionId"].(map[string]interface{})
	idSubtipo := subTipoProduccionID["Id"]
	idSubtipoStr := fmt.Sprintf("%v", idSubtipo)
	idSubtipoInt, _ := strconv.Atoi(idSubtipoStr)
	Metadatos := produccionAcademica["Metadatos"].([]interface{})
	var valor int
	valor = 1
	var autores float64
	autores = 0
	if idSubtipoInt == 4 || idSubtipoInt == 5 || idSubtipoInt == 6 {
		valor, autores = findGradePoints(SolicitudProduccion, idSubtipoInt)
	} else {
		valor, autores = findCategoryPoints(Metadatos)
	}

	if SolicitudProduccionResult, errPuntaje := addResult(SolicitudProduccion, idSubtipoStr, valor, autores); errPuntaje == nil {
		return SolicitudProduccionResult, nil
	} else {
		logs.Error(SolicitudProduccion)
		return nil, errPuntaje
	}
}

func findCategoryPoints(Metadatos []interface{}) (valorNum int, autoresNum float64) {
	var autores float64
	autores = 0
	var valor int
	valor = 1
	for _, metaDatotemp := range Metadatos {
		metaDato := metaDatotemp.(map[string]interface{})
		metaDatoSubtipo := metaDato["MetadatoSubtipoProduccionId"].(map[string]interface{})
		tipoMetadatoID := metaDatoSubtipo["TipoMetadatoId"].(map[string]interface{})
		idTipoMetadato := tipoMetadatoID["Id"]
		idTipoMetadatoStr := fmt.Sprintf("%v", idTipoMetadato)
		idSubtipoInt, _ := strconv.Atoi(idTipoMetadatoStr)
		if idSubtipoInt == 38 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			valor, _ = strconv.Atoi(numTipoMetadatoStr)
		} else if idSubtipoInt == 43 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			valor, _ = strconv.Atoi(numTipoMetadatoStr)
		} else if idSubtipoInt == 44 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			valor, _ = strconv.Atoi(numTipoMetadatoStr)
		}
		if idSubtipoInt == 21 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			autores, _ = strconv.ParseFloat(numTipoMetadatoStr, 64)
		}
	}
	return valor, autores
}

func findGradePoints(SolicitudProduccion map[string]interface{}, idSubtipoInt int) (valorNum int, autores float64) {
	var producciones []map[string]interface{}
	var tercero string
	solicitantes := SolicitudProduccion["Solicitantes"].([]interface{})
	for _, solicitantestemp := range solicitantes {
		solicitante := solicitantestemp.(map[string]interface{})
		tercero = fmt.Sprintf("%v", solicitante["TerceroId"])
	}
	idTercero := fmt.Sprintf("%v", tercero)
	valorNum = 1
	autores = 0
	var numEspecializacion int
	var numMaestria int
	var numDoctorado int
	numDoctorado = 0
	numEspecializacion = 0
	numMaestria = 0
	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				subtipoIdStr := fmt.Sprintf("%v", produccion["SubtipoProduccionId"].(map[string]interface{})["Id"])
				subtipoId, _ := strconv.Atoi(subtipoIdStr)
				if subtipoId == 4 {
					numEspecializacion++
				} else if subtipoId == 5 {
					numMaestria++
				} else if subtipoId == 6 {
					numDoctorado++
				}

			}
		}
	}
	fmt.Println(numMaestria)
	fmt.Println(numEspecializacion)
	if idSubtipoInt == 4 {
		if numEspecializacion >= 2 {
			//editar este valor a -1 para ajustar cuentas
			valorNum = -1
		} else if numEspecializacion == 1 {
			valorNum = 2
		} else if numEspecializacion == 0 {
			valorNum = 1
		}
	} else if idSubtipoInt == 5 {
		if numMaestria == 0 && numEspecializacion <= 1 {
			valorNum = 1
		} else if numMaestria == 0 && numEspecializacion >= 2 {
			//editar if numEspecializacion >=2 a ==2
			valorNum = 2
		} else if numMaestria == 1 && numEspecializacion == 0 {
			// editar if numMaestria >=1 a ==1
			valorNum = 3
		} else {
			//cambiar a 0 el valornum
			valorNum = -1
		}
	} else if idSubtipoInt == 6 {
		if numMaestria >= 1 && numDoctorado == 0 {
			valorNum = 1
		} else if numMaestria == 0 && numDoctorado == 0 {
			valorNum = 2
		} else {
			valorNum = -1
		}
	}

	return valorNum, autores
}

func addResult(SolicitudProduccion map[string]interface{}, idSubtipoStr string, valor int, autores float64) (result map[string]interface{}, outputError interface{}) {
	var resultado float64
	var puntajes []map[string]interface{}
	if valor == (-1) {
		resultado = 0.0
		resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
		SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
	} else {
		errPuntaje := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/puntaje_subtipo_produccion/?query=SubTipoProduccionId:"+idSubtipoStr+"&sortby=Id&order=asc", &puntajes)
		if errPuntaje == nil && fmt.Sprintf("%v", puntajes[0]["System"]) != "map[]" {
			if puntajes[0]["Status"] != 404 && puntajes[0]["Id"] != nil {

				Puntajes := puntajes[valor-1]

				type Caracteristica struct {
					Puntaje string
				}
				var caracteristica Caracteristica
				json.Unmarshal([]byte(fmt.Sprintf("%v", Puntajes["Caracteristicas"])), &caracteristica)
				puntajeStr := caracteristica.Puntaje
				puntajeStrF := strings.ReplaceAll(puntajeStr, ",", ".")
				puntajeInt, _ := strconv.ParseFloat(puntajeStrF, 64)

				if autores <= 3 && autores > 0 {
					resultado = puntajeInt
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				} else if autores > 3 && autores <= 5 {
					resultado = (puntajeInt / 2)
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				} else if autores > 5 {
					resultado = (puntajeInt / autores)
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				} else {
					resultado = puntajeInt
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				}

				return SolicitudProduccion, nil

			}
		} else {
			logs.Error(puntajes)
			return nil, errPuntaje
		}
	}

	return SolicitudProduccion, nil
}

// CheckCoincidenceProduction is...
func CheckCoincidenceProduction(SolicitudProduccion map[string]interface{}, idTipoProduccion int, idTercero string) (result map[string]interface{}, outputError interface{}) {
	var idSolicitudesList []float64
	var solicitudes []map[string]interface{}
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/inactive/", &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
			var produccionActual map[string]interface{}
			produccionActual = SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})

			for _, solicitud := range solicitudes {
				if fmt.Sprintf("%v", solicitud["Solicitantes"].([]interface{})[0].(map[string]interface{})["TerceroId"]) != idTercero {
					type Reference struct{ Id int }
					var reference Reference
					json.Unmarshal([]byte(fmt.Sprintf("%v", solicitud["Referencia"])), &reference)
					if produccionList, errProduccion := GetOneProduccionAcademica(fmt.Sprintf("%v", reference.Id)); errProduccion == nil {
						produccion := produccionList[0].(map[string]interface{})

						if fmt.Sprintf("%v", produccion["SubtipoProduccionId"].(map[string]interface{})["Id"]) == fmt.Sprintf("%v", produccionActual["ProduccionAcademica"].(map[string]interface{})["SubtipoProduccionId"].(map[string]interface{})["Id"]) {
							distance := CheckTitle(produccionActual["ProduccionAcademica"].(map[string]interface{}), produccion)
							if distance < 3 {
								idSolicitudesList = append(idSolicitudesList, solicitud["Id"].(float64))
							}
						}
					} else {
						logs.Error(produccionList)
						return nil, errProduccion
					}
				}
			}

			generateAlertCoincidences(SolicitudProduccion, idSolicitudesList)
			return SolicitudProduccion, nil
		}
	} else {
		logs.Error(solicitudes)
		return nil, errSolicitud
	}
	return SolicitudProduccion, nil
}

func generateAlertCoincidences(SolicitudDocente map[string]interface{}, idCoincidences []float64) {
	var observaciones []interface{}
	var idList string
	var tipoObservacionData map[string]interface{}

	if len(idCoincidences) > 0 {

		for _, id := range idCoincidences {
			idList += fmt.Sprintf("%v", id) + ","
		}

		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tipo_observacion/?query=Id:4", &tipoObservacionData)
		if errSolicitud == nil && fmt.Sprintf("%v", tipoObservacionData["System"]) != "map[]" {
			if tipoObservacionData["Status"] != 404 && tipoObservacionData["Data"] != nil {

				var tipoObservacion interface{}
				tipoObservacion = tipoObservacionData["Data"].([]interface{})[0]

				if SolicitudDocente["Observaciones"] != nil {
					observaciones = SolicitudDocente["Observaciones"].([]interface{})
				}

				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             idList,
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
				SolicitudDocente["Observaciones"] = observaciones
			}
		}
	}
}

// GenerateEvaluationsCloning is ...
func GenerateEvaluationsCloning(SolicitudProduccion map[string]interface{}, idSolicitud string, idSolicitudCoincidencia string, idTerceroSrt string) (result []map[string]interface{}, outputError interface{}) {
	idTercero, _ := strconv.Atoi(idTerceroSrt)
	var solicitudesEvaluaciones []map[string]interface{}
	var resultado []map[string]interface{}

	errEvaluaciones := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/solicitud/?limit=0&query=SolicitudPadreId:"+idSolicitudCoincidencia, &solicitudesEvaluaciones)
	if errEvaluaciones == nil && fmt.Sprintf("%v", solicitudesEvaluaciones[0]["System"]) != "map[]" {
		if solicitudesEvaluaciones[0]["Status"] != 404 && solicitudesEvaluaciones[0]["Id"] != nil {
			for _, evaluacion := range solicitudesEvaluaciones {
				if evaluacion["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"].(float64) == 13 {

					var evaluadores []interface{}
					errEvaluadores := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/solicitante/?query=SolicitudId:"+fmt.Sprintf("%v", evaluacion["Id"]), &evaluadores)
					if errEvaluadores == nil {
						SolicitudEvaluacion := make(map[string]interface{})
						SolicitudEvaluacion["Evaluacion"] = map[string]interface{}{
							"Autores":               evaluadores,
							"EstadoTipoSolicitudId": evaluacion["EstadoTipoSolicitudId"],
							"Referencia":            evaluacion["Referencia"],
							"Resultado":             evaluacion["Resultado"],
							"TerceroId":             idTercero,
							"SolicitudPadreId":      SolicitudProduccion,
						}
						if solicitudPost, errPost := PostSolicitudDocente(SolicitudEvaluacion["Evaluacion"].(map[string]interface{})); errPost == nil {
							resultado = append(resultado, solicitudPost)
						} else {
							logs.Error(solicitudPost)
							return nil, errPost
						}
					} else {
						logs.Error(evaluadores)
						return nil, errEvaluadores
					}
				}
			}
			return resultado, nil
		}
	} else {
		logs.Error(solicitudesEvaluaciones)
		return nil, errEvaluaciones
	}

	return resultado, nil
}

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
