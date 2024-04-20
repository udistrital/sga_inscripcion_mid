package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_inscripcion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/time_bogota"

	"github.com/k0kubun/pp"
)

func ProduccionAcademicaPost(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	//resultado experiencia
	var resultado map[string]interface{}
	var produccionAcademica map[string]interface{}

	date := time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(data, &produccionAcademica); err == nil {
		produccionAcademicaPost := make(map[string]interface{})
		produccionAcademicaPost["ProduccionAcademica"] = map[string]interface{}{
			"Titulo":              produccionAcademica["Titulo"],
			"Resumen":             produccionAcademica["Resumen"],
			"Fecha":               produccionAcademica["Fecha"],
			"SubtipoProduccionId": produccionAcademica["SubtipoProduccionId"],
			"Activo":              true,
			"FechaCreacion":       date,
			"FechaModificacion":   date,
		}

		var autores []map[string]interface{}
		for _, autorTemp := range produccionAcademica["Autores"].([]interface{}) {
			autor := autorTemp.(map[string]interface{})
			autores = append(autores, map[string]interface{}{
				"Persona":                 autor["PersonaId"],
				"EstadoAutorProduccionId": autor["EstadoAutorProduccionId"],
				"ProduccionAcademicaId":   map[string]interface{}{"Id": 0},
				"Activo":                  true,
				"FechaCreacion":           date,
				"FechaModificacion":       date,
			})
		}
		produccionAcademicaPost["Autores"] = autores

		var metadatos []map[string]interface{}
		for _, metadatoTemp := range produccionAcademica["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			metadatos = append(metadatos, map[string]interface{}{
				"Valor": fmt.Sprintf("%v", metadato["Valor"]),
				// "MetadatoSubtipoProduccionId": metadato["MetadatoSubtipoProduccionId"],
				"MetadatoSubtipoProduccionId": map[string]interface{}{"Id": metadato["MetadatoSubtipoProduccionId"]},
				"ProduccionAcademicaId":       map[string]interface{}{"Id": 0},
				"Activo":                      true,
				"FechaCreacion":               date,
				"FechaModificacion":           date,
			})
		}
		produccionAcademicaPost["Metadatos"] = metadatos
		var resultadoProduccionAcademica map[string]interface{}
		errProduccion := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica", "POST", &resultadoProduccionAcademica, produccionAcademicaPost)
		if errProduccion == nil && fmt.Sprintf("%v", resultadoProduccionAcademica["System"]) != "map[]" && resultadoProduccionAcademica["ProduccionAcademica"] != nil {
			if resultadoProduccionAcademica["Status"] != 400 {
				resultado = resultadoProduccionAcademica
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
			} else {
				logs.Error(errProduccion)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errProduccion)
				return APIResponseDTO
			}
		} else {
			logs.Error(errProduccion)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errProduccion)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func EstadoAutorProduccion(idAutor string, data []byte) (APIResponseDTO requestresponse.APIResponse) {
	//resultado experiencia
	var resultado map[string]interface{}
	var dataPut map[string]interface{}

	if err := json.Unmarshal(data, &dataPut); err == nil {
		fmt.Println("data put", dataPut)
		var acepta = dataPut["acepta"].(bool)
		var AutorProduccionAcademica = dataPut["AutorProduccionAcademica"].(map[string]interface{})
		if acepta {
			(AutorProduccionAcademica["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 2
		} else {
			(AutorProduccionAcademica["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 4
		}
		var resultadoAutor map[string]interface{}
		errAutor := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/autor_produccion_academica/"+idAutor, "PUT", &resultadoAutor, AutorProduccionAcademica)
		pp.Println(resultadoAutor)
		if errAutor == nil && fmt.Sprintf("%v", resultadoAutor["System"]) != "map[]" && resultadoAutor["Id"] != nil {
			if resultadoAutor["Status"] != 400 {
				resultado = AutorProduccionAcademica
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
			} else {
				logs.Error(errAutor)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errAutor)
				return APIResponseDTO
			}
		} else {
			logs.Error(errAutor)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errAutor)
			return APIResponseDTO
		}

	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func ProduccionAcademicaPut(idProduccion string, data []byte) (APIResponseDTO requestresponse.APIResponse) {

	date := time_bogota.TiempoBogotaFormato()

	//resultado experiencia
	var resultado map[string]interface{}
	//produccion academica
	var produccionAcademica map[string]interface{}
	if err := json.Unmarshal(data, &produccionAcademica); err == nil {
		produccionAcademicaPut := make(map[string]interface{})
		produccionAcademicaPut["ProduccionAcademica"] = map[string]interface{}{
			"Titulo":              produccionAcademica["Titulo"],
			"Resumen":             produccionAcademica["Resumen"],
			"Fecha":               produccionAcademica["Fecha"],
			"SubtipoProduccionId": produccionAcademica["SubtipoProduccionId"],
			"FechaModificacion":   date,
		}

		var metadatos []map[string]interface{}
		for _, metadatoTemp := range produccionAcademica["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			metadatos = append(metadatos, map[string]interface{}{
				"Valor":                       fmt.Sprintf("%v", metadato["Valor"]),
				"MetadatoSubtipoProduccionId": map[string]interface{}{"Id": metadato["MetadatoSubtipoProduccionId"]},
				"Activo":                      true,
				"FechaModificacion":           date,
			})
		}

		produccionAcademicaPut["Autores"] = nil
		produccionAcademicaPut["Metadatos"] = metadatos

		var resultadoProduccionAcademica map[string]interface{}

		errProduccion := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idProduccion, "PUT", &resultadoProduccionAcademica, produccionAcademicaPut)
		if errProduccion == nil && fmt.Sprintf("%v", resultadoProduccionAcademica["System"]) != "map[]" {

			if resultadoProduccionAcademica["Status"] != 400 {
				resultado = produccionAcademica
				APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
			} else {
				logs.Error(errProduccion)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errProduccion)
				return APIResponseDTO
			}
		} else {
			logs.Error(errProduccion)
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errProduccion)
			return APIResponseDTO
		}
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func GetProduccionById(idProduccion string) (APIResponseDTO requestresponse.APIResponse) {

	//resultado experiencia
	var resultadoGetProduccion []interface{}
	if resultado, err := helpers.GetOneProduccionAcademica(idProduccion); err == nil {
		resultadoGetProduccion = resultado
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultadoGetProduccion, nil)
	} else {
		logs.Error(err)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err)
		return APIResponseDTO
	}

	return APIResponseDTO
}

func GetAllProducciones() (APIResponseDTO requestresponse.APIResponse) {
	fmt.Println("Consultando todas las producciones")
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var producciones []map[string]interface{}

	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/?limit=0", &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				autores := produccion["Autores"].([]interface{})
				for _, autorTemp := range autores {
					autor := autorTemp.(map[string]interface{})
					produccion["EstadoEnteAutorId"] = autor
					//cargar nombre del autor
					var autorProduccion map[string]interface{}

					errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
					if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
						if autorProduccion["Status"] != 404 {
							autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
						} else {
							if autorProduccion["Message"] == "Not found resource" {
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Not found resource")
								return APIResponseDTO
							} else {
								logs.Error(autorProduccion)
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errAutor)
								return APIResponseDTO
							}
						}
					} else {
						logs.Error(autorProduccion)
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errAutor)
						return APIResponseDTO
					}
				}
			}
			resultado = producciones
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		} else {
			if producciones[0]["Message"] == "Not found resource" {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
				return APIResponseDTO
			} else {
				logs.Error(producciones, errProduccion)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errProduccion)
				return APIResponseDTO
			}
		}
	} else {
		logs.Error(producciones)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errProduccion)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func GetIdProduccion(idTercero string) (APIResponseDTO requestresponse.APIResponse) {
	var resultado []map[string]interface{}
	var producciones []map[string]interface{}
	var errorGetAll bool

	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"tr_produccion_academica/"+idTercero, &producciones)
	fmt.Println("//////////// ProduccionAcademicaService() Err: ", errProduccion, "Resp: ", producciones)
	if fmt.Sprintf("%v", producciones) != "" || fmt.Sprintf("%v", producciones) != "[map[]]" {
		if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
			if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
				for _, produccion := range producciones {
					autores := produccion["Autores"].([]interface{})
					for _, autorTemp := range autores {
						autor := autorTemp.(map[string]interface{})
						if fmt.Sprintf("%v", autor["Persona"]) == fmt.Sprintf("%v", idTercero) {
							produccion["EstadoEnteAutorId"] = autor
						}
						//cargar nombre del autor
						var autorProduccion map[string]interface{}

						errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
						fmt.Println("//////////// TercerosService() Err: ", errAutor, "Resp: ", autorProduccion)
						if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
							if autorProduccion["Status"] != 404 {
								autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
							} else {
								errorGetAll = true
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data founf")
							}
						} else {
							errorGetAll = true
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errAutor.Error())
						}
					}
				}
				resultado = producciones
			} else {
				errorGetAll = true
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "No data found")
			}
		} else {
			errorGetAll = true
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errProduccion.Error())
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
	}

	if !errorGetAll {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
	}
	return APIResponseDTO
}

func GetProduccion(idTercero string) (APIResponseDTO requestresponse.APIResponse) {
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var producciones []map[string]interface{}

	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				autores := produccion["Autores"].([]interface{})
				for _, autorTemp := range autores {
					autor := autorTemp.(map[string]interface{})
					if fmt.Sprintf("%v", autor["Persona"]) == fmt.Sprintf("%v", idTercero) {
						produccion["EstadoEnteAutorId"] = autor
					}
					//cargar nombre del autor
					var autorProduccion map[string]interface{}

					errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
					if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
						if autorProduccion["Status"] != 404 {
							autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
						} else {
							if autorProduccion["Message"] == "Not found resource" {
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not found resource")
								return APIResponseDTO
							} else {
								logs.Error(autorProduccion)
								APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errAutor)
								return APIResponseDTO
							}
						}
					} else {
						logs.Error(autorProduccion)
						APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errAutor)
						return APIResponseDTO
					}
				}
			}
			resultado = producciones
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, resultado, nil)
		} else {
			if producciones[0]["Message"] == "Not found resource" {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, "Not dound resource")
				return APIResponseDTO
			} else {
				logs.Error(producciones)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errProduccion)
				return APIResponseDTO
			}
		}
	} else {
		logs.Error(producciones)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil, errProduccion)
		return APIResponseDTO
	}
	return APIResponseDTO
}

func DeleteProduccion(idProduccion string) (APIResponseDTO requestresponse.APIResponse) {

	//resultados eliminacion
	var borrado map[string]interface{}

	errDelete := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idProduccion, "DELETE", &borrado, nil)
	//borradoOk := models.SetInactivo("http://" + beego.AppConfig.String("ProduccionAcademicaService") + "/tr_produccion_academica/" + idStr)

	if errDelete == nil {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, map[string]interface{}{"ProduccionAcademica": idProduccion}, nil)
	} else {
		logs.Error("Failed deleting tr_produccion_academica/" + idProduccion)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "Failed deleting tr_produccion_academica/"+idProduccion)
		return APIResponseDTO
	}
	return APIResponseDTO
}
