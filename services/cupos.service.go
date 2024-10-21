package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/time_bogota"
	"golang.org/x/sync/errgroup"
)

// Funcion para recibir todos los cupos para una inscripcion
func GetAllCuposInscripcion(periodo string, proyecto string, tipo string) (APIResponseDTO requestresponse.APIResponse) {

	var cupo []map[string]interface{}
	var listado []map[string]interface{}
	//Definición del group para las gorutines
	wge := new(errgroup.Group)
	var mutex sync.Mutex // Mutex para proteger el acceso a resultados

	errCupos := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"cupo_inscripcion?query=Activo:true,PeriodoId:"+periodo+",ProyectoAcademicoId:"+proyecto+",TipoInscripcionId.Id:"+tipo+",Activo:true&limit=0", &cupo)

	if errCupos != nil {
		return requestresponse.APIResponseDTO(false, 400, nil, "Error en la consulta de cupos")
	}

	if cupo == nil || len(cupo[0]) == 0 {
		return requestresponse.APIResponseDTO(true, 200, nil, "No se encontraron asignaciones de cupos")
	}

	wge.SetLimit(-1)
	for _, c := range cupo {
		c := c
		wge.Go(func() error {
			var cupoContenido = make(map[string]interface{})
			tipoInscripcionId := c["TipoInscripcionId"].(map[string]interface{})
			idIns := tipoInscripcionId["Id"].(float64)
			nombreIns := tipoInscripcionId["Nombre"].(string)
			cupoContenido["Activo"] = c["Activo"]
			cupoContenido["CuposHabilitados"] = c["CuposHabilitados"]
			cupoContenido["CuposOpcionados"] = c["CuposOpcionados"]
			cupoContenido["CuposDisponibles"] = c["CuposDisponibles"]
			cupoContenido["PeriodoId"] = c["PeriodoId"]
			cupoContenido["ProyectoAcademicoId"] = c["ProyectoAcademicoId"]
			cupoContenido["FechaCreacion"] = c["FechaCreacion"]
			cupoContenido["CupoId"] = c["CupoId"]
			cupoContenido["Id"] = c["Id"]
			cupoContenido["TipoInscripcionId"] = idIns
			cupoContenido["NombreInscripcion"] = nombreIns
			idcupo := c["CupoId"].(float64)

			var tipocupo map[string]interface{}
			errtipocupo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"/parametro?query=TipoParametroId__Id:87,Id:"+fmt.Sprintf("%v", idcupo)+"&limit=0", &tipocupo)
			if errtipocupo == nil && tipocupo["Status"] == "200" && fmt.Sprintf("%v", tipocupo["Data"]) != "[map[]]" {
				cupoContenido["Nombre"] = tipocupo["Data"].([]interface{})[0].(map[string]interface{})["Nombre"]
				cupoContenido["Descripcion"] = tipocupo["Data"].([]interface{})[0].(map[string]interface{})["Descripcion"]
			}

			mutex.Lock()
			listado = append(listado, cupoContenido)
			mutex.Unlock()

			return errtipocupo
		})
	}
	//Si existe error, se realiza
	if err := wge.Wait(); err != nil {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, listado)
	}

	return APIResponseDTO
}

func UpdateCuposInscripcion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	//var body map[string]interface{}

	var cupoActualizado map[string]interface{}
	var date = time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(data, &cupoActualizado); err == nil {
		fmt.Println("Put")
		idcupo := cupoActualizado["Id"].(float64)
		dataActualizada := map[string]interface{}{
			"Activo":              cupoActualizado["Activo"],
			"FechaCreacion":       cupoActualizado["FechaCreacion"],
			"FechaModificacion":   date,
			"Id":                  cupoActualizado["Id"],
			"CuposHabilitados":    cupoActualizado["CuposHabilitados"],
			"CuposOpcionados":     cupoActualizado["CuposOpcionados"],
			"PeriodoId":           cupoActualizado["PeriodoId"],
			"ProyectoAcademicoId": cupoActualizado["ProyectoAcademicoId"],
			"CupoId":              cupoActualizado["CupoId"],
			"TipoInscripcionId":   cupoActualizado["TipoInscripcionId"],
		}
		cupoActualizado = dataActualizada
		errActualizarCupo := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/cupo_inscripcion/"+fmt.Sprintf("%.f", idcupo), "PUT", &cupoActualizado, dataActualizada)
		if errActualizarCupo == nil {
			return requestresponse.APIResponseDTO(false, 200, cupoActualizado, dataActualizada)
		} else {
			return requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar datos JSON")
		}

	} else {
		return requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar datos JSON")
	}
}
func GetAllDocCupos() (APIResponseDTO requestresponse.APIResponse) {
	fmt.Println("GetAll")
	var docCupo []map[string]interface{}

	var listado []map[string]interface{}
	errCupos := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+fmt.Sprintf("/documento_cupo?query=Activo:true&limit=0"), &docCupo)
	if errCupos == nil {

		for _, c := range docCupo {
			var cupoContenido = make(map[string]interface{})
			cupoInscripcionId := c["CupoInscripcionId"].(map[string]interface{})
			idCupo := cupoInscripcionId["Id"].(float64)
			idTipoCupo := cupoInscripcionId["CupoId"].(float64)
			cupoContenido["Activo"] = c["Activo"]
			cupoContenido["Uid"] = c["Uid"]
			cupoContenido["CuposOpcionados"] = c["CuposOpcionados"]
			cupoContenido["Comentario"] = c["Comentario"]
			cupoContenido["FechaCreacion"] = c["FechaCreacion"]
			cupoContenido["Id"] = c["Id"]
			cupoContenido["CupoInscripcionId"] = idCupo

			var tipocupo map[string]interface{}
			errtipocupo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"/parametro?query=TipoParametroId__Id:87,Id:"+fmt.Sprintf("%v", idTipoCupo)+"&limit=0", &tipocupo)
			//fmt.Println(ProyectoV2["Data"])
			if errtipocupo == nil && tipocupo["Status"] == "200" && fmt.Sprintf("%v", tipocupo["Data"]) != "[map[]]" {
				cupoContenido["Cupo"] = tipocupo["Data"].([]interface{})[0].(map[string]interface{})["Nombre"]
			} else {
			}

			listado = append(listado, cupoContenido)
		}

		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, listado)

	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errCupos.Error())
	}
	return APIResponseDTO
}
func PostDocCupos(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	//Almacena los comentarios con su documento por nuxeo (falta implementacion)
	var nuevoComentario map[string]interface{}
	var cupo map[string]interface{}
	//respuesta a la petición
	var respuesta map[string]interface{}
	//timestamp
	date := time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(data, &nuevoComentario); err == nil {
		idCupoIns := nuevoComentario["IdCupoIns"]

		errCupos := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+fmt.Sprintf("/cupo_inscripcion/")+fmt.Sprintf("%.f", idCupoIns), &cupo)
		if errCupos == nil {

			dataComentario := map[string]interface{}{
				"Activo":            true,
				"FechaCreacion":     date,
				"FechaModificacion": date,
				"Comentario":        nuevoComentario["Comentario"],
				"Uid":               nuevoComentario["Uid"],
				"CupoInscripcionId": cupo,
				"Id":                nuevoComentario["Id"],
			}

			nuevoComentario = dataComentario
			errNoticia := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/documento_cupo/", "POST", &nuevoComentario, dataComentario)
			if errNoticia == nil {

				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nuevoComentario)
				return APIResponseDTO
			} else {
				models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"/cupo_inscripcion/%.f", nuevoComentario["Data"].(map[string]interface{})["Id"].(float64)))
			}

			APIResponseDTO = requestresponse.APIResponseDTO(true, 500, respuesta, nuevoComentario)
			return APIResponseDTO
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(true, 500, respuesta, nuevoComentario)
			return APIResponseDTO
		}

	}
	APIResponseDTO = requestresponse.APIResponseDTO(true, 500, respuesta, nuevoComentario)
	return APIResponseDTO
}

func PostCuposInscripcion(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var registros interface{}
	var respuesta map[string]interface{}
	var errores []string

	if err := json.Unmarshal(data, &registros); err == nil {
		date := time_bogota.TiempoBogotaFormato()

		switch registros := registros.(type) {
		case map[string]interface{}:
			if registros["Id"] != nil {
				idCupo := registros["Id"].(float64)
				cupoActualizado := map[string]interface{}{
					"Activo":              true,
					"FechaCreacion":       date,
					"FechaModificacion":   date,
					"CuposHabilitados":    registros["CuposHabilitados"],
					"CuposOpcionados":     registros["CuposOpcionados"],
					"CuposDisponibles":    registros["CuposDisponibles"],
					"PeriodoId":           registros["PeriodoId"],
					"ProyectoAcademicoId": registros["ProyectoAcademicoId"],
					"CupoId":              registros["CupoId"],
					"TipoInscripcionId": map[string]interface{}{
						"Id": registros["TipoInscripcionId"],
					},
				}

				errActualizarCupo := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"cupo_inscripcion/"+fmt.Sprintf("%.f", idCupo), "PUT", &cupoActualizado, respuesta)

				if errActualizarCupo != nil {
					errores = append(errores, errActualizarCupo.Error())
				}
			}
		case []interface{}:
			for _, registro := range registros {
				if registroMap, ok := registro.(map[string]interface{}); ok {
					if registroMap["Id"] != nil {
						idCupo := registroMap["Id"].(float64)
						cupoActualizado := map[string]interface{}{
							"Activo":              true,
							"FechaCreacion":       date,
							"FechaModificacion":   date,
							"CuposHabilitados":    registroMap["CuposHabilitados"],
							"CuposOpcionados":     registroMap["CuposOpcionados"],
							"CuposDisponibles":    registroMap["CuposDisponibles"],
							"PeriodoId":           registroMap["PeriodoId"],
							"ProyectoAcademicoId": registroMap["ProyectoAcademicoId"],
							"CupoId":              registroMap["CupoId"],
							"TipoInscripcionId": map[string]interface{}{
								"Id": registroMap["TipoInscripcionId"],
							},
						}
						errActualizarCupo := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/cupo_inscripcion/"+fmt.Sprintf("%.f", idCupo), "PUT", &respuesta, cupoActualizado)

						if errActualizarCupo != nil {
							errores = append(errores, errActualizarCupo.Error())

						}
					} else {

						cupoNuevo := map[string]interface{}{

							"Activo":              true,
							"FechaCreacion":       date,
							"FechaModificacion":   date,
							"CuposHabilitados":    registroMap["CuposHabilitados"],
							"CuposOpcionados":     registroMap["CuposOpcionados"],
							"CuposDisponibles":    registroMap["CuposDisponibles"],
							"PeriodoId":           registroMap["PeriodoId"],
							"ProyectoAcademicoId": registroMap["ProyectoAcademicoId"],
							"CupoId":              registroMap["CupoId"],
							"TipoInscripcionId": map[string]interface{}{
								"Id": registroMap["TipoInscripcionId"],
							},
						}

						errActualizarCupo := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/cupo_inscripcion/", "POST", &respuesta, cupoNuevo)
						if errActualizarCupo != nil {
							errores = append(errores, errActualizarCupo.Error())
						} else {

							dataComentario := map[string]interface{}{

								"Activo":            true,
								"FechaCreacion":     date,
								"FechaModificacion": date,
								"Comentario":        registroMap["Comentario"],
								"Uid":               registroMap["Enlace"],
								"CupoInscripcionId": map[string]interface{}{
									"Id": respuesta["Id"],
								},
							}
							errDocumentoCupo := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/documento_cupo/", "POST", &respuesta, dataComentario)
							fmt.Println(respuesta)
							if errDocumentoCupo != nil {
								errores = append(errores, errDocumentoCupo.Error())
							}
						}

					}
				}
			}
			return requestresponse.APIResponseDTO(true, 200, respuesta)
		default:
			return requestresponse.APIResponseDTO(false, 400, nil, "Formato de datos no válido")
		}
	} else {
		fmt.Println("Error al decodificar datos JSON:", err)
		return requestresponse.APIResponseDTO(false, 400, nil, "Error al decodificar datos JSON: "+err.Error())
	}

	return requestresponse.APIResponseDTO(true, 200, respuesta, respuesta)
}
