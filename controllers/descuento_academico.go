package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// DescuentoController ...
type DescuentoController struct {
	beego.Controller
}

// URLMapping ...
func (c *DescuentoController) URLMapping() {
	c.Mapping("PostDescuentoAcademico", c.PostDescuentoAcademico)
	c.Mapping("PutDescuentoAcademico", c.PutDescuentoAcademico)
	c.Mapping("GetDescuentoAcademico", c.GetDescuentoAcademico)
	c.Mapping("GetDescuentoAcademicoByPersona", c.GetDescuentoAcademicoByPersona)
	c.Mapping("GetDescuentoByPersonaPeriodoDependencia", c.GetDescuentoByPersonaPeriodoDependencia)
	c.Mapping("GetDescuentoAcademicoByDependenciaID", c.GetDescuentoAcademicoByDependenciaID)
}

// PostDescuentoAcademico ...
// @Title PostDescuentoAcademico
// @Description Agregar Descuento Academico
// @Param	body		body 	{}	true		"body Agregar Descuento Academico content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *DescuentoController) PostDescuentoAcademico() {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var solicitud map[string]interface{}
	var solicitudPost map[string]interface{}
	var tipoDescuento []map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err == nil {
		IDTipoDescuento := fmt.Sprintf("%v", solicitud["DescuentosDependenciaId"].(map[string]interface{})["Id"])
		IDDependencia := fmt.Sprintf("%v", solicitud["DescuentosDependenciaId"].(map[string]interface{})["Dependencia"])
		errDescuentosDependencia := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia?query=TipoDescuentoId__Id:"+IDTipoDescuento+",DependenciaId:"+IDDependencia, &tipoDescuento)
		if errDescuentosDependencia == nil && fmt.Sprintf("%v", tipoDescuento[0]["System"]) != "map[]" {

			// DescuentosDependenciaID := map[string]interface{}{
			// 	"Activo":          solicitud["DescuentosDependenciaId"].(map[string]interface{})["Activo"],
			// 	"DependenciaId":   solicitud["DescuentosDependenciaId"].(map[string]interface{})["Dependencia"],
			// 	"PeriodoId":       solicitud["DescuentosDependenciaId"].(map[string]interface{})["Periodo"],
			// 	"TipoDescuentoId": tipoDescuento,
			// }

			solicituddescuento := map[string]interface{}{
				"Id":                      0,
				"TerceroId":               solicitud["PersonaId"],
				"Estado":                  "Por aprobar",
				"PeriodoId":               solicitud["PeriodoId"],
				"Activo":                  true,
				"DescuentosDependenciaId": tipoDescuento[0],
			}
			formatdata.JsonPrint(solicituddescuento)
			// fmt.Println(solicituddescuento)

			errSolicitud := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento", "POST", &solicitudPost, solicituddescuento)
			if errSolicitud == nil && fmt.Sprintf("%v", solicitudPost["System"]) != "map[]" && solicitudPost["Id"] != nil {
				if solicitudPost["Status"] != 400 {
					//soporte de descuento
					var soporte map[string]interface{}

					soportedescuento := map[string]interface{}{
						"SolicitudDescuentoId": solicitudPost,
						"Activo":               true,
						"DocumentoId":          solicitud["DocumentoId"],
					}

					errSoporte := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento", "POST", &soporte, soportedescuento)
					if errSoporte == nil && fmt.Sprintf("%v", soporte["System"]) != "map[]" && soporte["Id"] != nil {
						if soporte["Status"] != 400 {
							resultado = map[string]interface{}{"Id": solicitudPost["Id"], "PersonaId": solicitudPost["PersonaId"], "Estado": solicitudPost["Estado"], "PeriodoId": solicitudPost["PeriodoId"], "DescuentosDependenciaId": solicitudPost["DescuentosDependenciaId"]}
							resultado["DocumentoId"] = soporte["DocumentoId"]
							c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}

						} else {
							//resultado solicitud de descuento
							var resultado2 map[string]interface{}
							request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/%.f", solicitudPost["Id"]), "DELETE", &resultado2, nil)
							logs.Error(errSoporte)
							//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
							c.Data["message"] = "Error service PostDescuentoAcademico: " + soporte["Body"].(string)
							c.Abort("400")
						}
					} else {
						logs.Error(errSoporte)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["message"] = "Error service PostDescuentoAcademico: " + soporte["Body"].(string)
						c.Abort("400")
					}
				} else {
					logs.Error(errSolicitud)
					//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
					c.Data["message"] = "Error service PostDescuentoAcademico: " + solicitudPost["Body"].(string)
					c.Abort("400")
				}
			} else {
				logs.Error(errSolicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["message"] = "Error service PostDescuentoAcademico: " + solicitudPost["Body"].(string)
				c.Abort("400")
			}
		} else {
			logs.Error(errDescuentosDependencia)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["message"] = "Error service PostDescuentoAcademico: " + errDescuentosDependencia.Error()
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["message"] = "Error service PostDescuentoAcademico: " + err.Error()
		c.Abort("400")
	}

	c.ServeJSON()
}

// PutDescuentoAcademico ...
// @Title PutDescuentoAcademico
// @Description Modificar Descuento Academico
// @Param	id	path 	int	true		"el id de la solicitud de descuento a modificar"
// @Param	body		body 	{}	true		"body Modificar Descuento Academico content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *DescuentoController) PutDescuentoAcademico() {
	idStr := c.Ctx.Input.Param(":id")
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var solicitud map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err == nil {
		//soporte de descuento
		var soporte []map[string]interface{}
		var soportePut map[string]interface{}

		errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=Activo:true,SolicitudDescuentoId:"+idStr, &soporte)
		if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
			if soporte[0]["Status"] != 404 {
				soporte[0]["DocumentoId"] = solicitud["DocumentoId"]

				errSoportePut := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/"+
					fmt.Sprintf("%v", soporte[0]["Id"]), "PUT", &soportePut, soporte[0])
				if errSoportePut == nil && fmt.Sprintf("%v", soportePut["System"]) != "map[]" && soportePut["Id"] != nil {
					if soportePut["Status"] != 400 {
						resultado = solicitud
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}

					} else {
						logs.Error(errSoportePut)
						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
						c.Data["message"] = "Error service PutDescuentoAcademico: " + soportePut["Body"].(string)
						c.Abort("400")
					}
				} else {
					logs.Error(errSoportePut)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["message"] = "Error service PutDescuentoAcademico: " + soportePut["Body"].(string)
					c.Abort("400")
				}

			} else {
				if soporte[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(soporte)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["message"] = "Error service PutDescuentoAcademico: " + errSoporte.Error()
					c.Abort("404")
				}
			}
		} else {
			logs.Error(soporte)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["message"] = "Error service PutDescuentoAcademico: " + errSoporte.Error()
			c.Abort("404")
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["message"] = "Error service PutDescuentoAcademico: " + err.Error()
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetDescuentoAcademico ...
// @Title GetDescuentoAcademico
// @Description consultar Descuento Academico por userid
// @Param	PersonaId		query 	int	true		"Id de la persona"
// @Param	SolicitudId		query 	int	true		"Id de la solicitud"
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *DescuentoController) GetDescuentoAcademico() {
	//Id de la persona
	idStr := c.GetString("PersonaId")
	fmt.Println("el id es: ", idStr)
	//Id de la solicitud
	idSolitudDes := c.GetString("SolicitudId")
	fmt.Println("el idSolitudDes es: ", idSolitudDes)
	//resultado consulta
	var resultado map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/?query=TerceroId:"+idStr+",Id:"+idSolitudDes+"&fields=Id,TerceroId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {
			resultado = solicitud[0]

			//resultado descuento dependencia
			var descuento map[string]interface{}
			errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+fmt.Sprintf("%v", solicitud[0]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
			if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
				if descuento["Status"] != 404 {
					//resultado tipo descuento
					var tipo map[string]interface{}
					errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
					if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
						if tipo["Status"] != 404 {
							descuento["TipoDescuentoId"] = tipo
							resultado["DescuentosDependenciaId"] = descuento

							//resultado soporte descuento
							var soporte []map[string]interface{}
							errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+idSolitudDes+"&fields=DocumentoId", &soporte)
							if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
								if soporte[0]["Status"] != 404 {
									//fmt.Println("el resultado de los documentos es: ", resultado4)
									resultado["DocumentoId"] = soporte[0]["DocumentoId"]
									c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
								} else {
									if soporte[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(soporte)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["message"] = "Error service GetDescuentoAcademico: " + errSoporte.Error()
										c.Abort("404")
									}
								}
							} else {
								logs.Error(soporte)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["message"] = "Error service GetDescuentoAcademico: " + errSoporte.Error()
								c.Abort("404")
							}
						} else {
							if tipo["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(tipo)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["message"] = "Error service GetDescuentoAcademico: " + errTipo.Error()
								c.Abort("404")
							}
						}
					} else {
						logs.Error(tipo)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["message"] = "Error service GetDescuentoAcademico: " + errTipo.Error()
						c.Abort("404")
					}
				} else {
					if descuento["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(descuento)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["message"] = "Error service GetDescuentoAcademico: " + errDescuento.Error()
						c.Abort("404")
					}
				}
			} else {
				logs.Error(descuento)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["message"] = "Error service GetDescuentoAcademico: " + errDescuento.Error()
				c.Abort("404")
			}
		} else {
			if solicitud[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(solicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["message"] = "Error service GetDescuentoAcademico: " + errSolicitud.Error()
				c.Abort("404")
			}
		}
	} else {
		logs.Error(solicitud)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["message"] = "Error service GetDescuentoAcademico: " + errSolicitud.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetDescuentoAcademicoByDependenciaID ...
// @Title GetDescuentoAcademicoByDependenciaID
// @Description consultar Descuento Academico por DependenciaId
// @Param	dependencia_id		path 	int	true		"DependenciaId"
// @Success 200 {}
// @Failure 404 not found resource
// @router /descuentoAcademicoByID/:dependencia_id [get]
func (c *DescuentoController) GetDescuentoAcademicoByDependenciaID() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":dependencia_id")
	//resultado consulta
	var resultados []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}
	//var alerta models.Alert
	var errorGetAll bool
	//alertas := append([]interface{}{"Data:"})

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia?limit=0&query=Activo:true,DependenciaId:"+idStr, &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {

			for u := 0; u < len(solicitud); u++ {
				var tipoDescuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", solicitud[u]["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipoDescuento)
				if errDescuento == nil && fmt.Sprintf("%v", tipoDescuento["System"]) != "map[]" {
					resultados = append(resultados, tipoDescuento)
				} else {
					errorGetAll = true

					logs.Error(errDescuento)
					c.Data["message"] = "Error service GetDescuentoAcademicoByDependenciaID: " + errDescuento.Error()
					c.Abort("400")
				}
			}
		} else {
			errorGetAll = true

			logs.Error("No data found")
			c.Data["message"] = "Error service GetDescuentoAcademicoByDependenciaID: " + "No data found"
			c.Abort("404")
		}
	} else {
		errorGetAll = true

		logs.Error(errSolicitud)
		c.Data["message"] = "Error service GetDescuentoAcademicoByDependenciaID: " + errSolicitud.Error()
		c.Abort("400")
	}
	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultados}
	}

	c.ServeJSON()

}

// GetDescuentoAcademicoByPersona ...
// @Title GetDescuentoAcademicoByPersona
// @Description consultar Descuento Academico por userid
// @Param	persona_id		path 	int	true		"Id de la persona"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:persona_id [get]
func (c *DescuentoController) GetDescuentoAcademicoByPersona() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":persona_id")
	fmt.Println("El id es: ", idStr)
	//resultado solicitud descuento
	var resultado []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/?query=PersonaId:"+idStr+"&fields=Id,PersonaId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {

			for u := 0; u < len(solicitud); u++ {
				//resultado solicitud descuento
				var descuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+
					fmt.Sprintf("%v", solicitud[u]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
				if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
					if descuento["Status"] != 404 {
						//resultado tipo descuento
						var tipo map[string]interface{}
						errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
						if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
							if tipo["Status"] != 404 {
								descuento["TipoDescuentoId"] = tipo
								solicitud[u]["DescuentosDependenciaId"] = descuento

								//resultado soporte descuento
								var soporte []map[string]interface{}
								errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+fmt.Sprintf("%v", solicitud[u]["Id"])+"&fields=DocumentoId", &soporte)
								if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
									if soporte[0]["Status"] != 404 {
										//fmt.Println("el resultado de los documentos es: ", resultado4)
										solicitud[u]["DocumentoId"] = soporte[0]["DocumentoId"]
									} else {
										if soporte[0]["Message"] == "Not found resource" {
											c.Data["json"] = nil
										} else {
											logs.Error(soporte)
											//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
											c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errSoporte.Error()
											c.Abort("404")
										}
									}
								} else {
									logs.Error(soporte)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errSoporte.Error()
									c.Abort("404")
								}
							} else {
								if tipo["Message"] == "Not found resource" {
									c.Data["json"] = nil
								} else {
									logs.Error(tipo)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errTipo.Error()
									c.Abort("404")
								}
							}
						} else {
							logs.Error(tipo)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errTipo.Error()
							c.Abort("404")
						}
					} else {
						if descuento["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(descuento)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errDescuento.Error()
							c.Abort("404")
						}
					}
				} else {
					logs.Error(descuento)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errDescuento.Error()
					c.Abort("404")
				}
			}
			resultado = solicitud
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
		} else {
			if solicitud[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(solicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errSolicitud.Error()
				c.Abort("404")
			}
		}
	} else {
		logs.Error(solicitud)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["message"] = "Error service GetDescuentoAcademicoByPersona: " + errSolicitud.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetDescuentoByPersonaPeriodoDependencia ...
// @Title GetDescuentoByPersonaPeriodoDependencia
// @Description consultar Descuento Academico por userid
// @Param	PersonaId		query 	int	true		"Id de la persona"
// @Param	DependenciaId		query 	int	true		"Id de la dependencia"
// @Param	PeriodoId		query 	int	true		"Id del periodo académico"
// @Success 200 {}
// @Failure 404 not found resource
// @router /descuentopersonaperiododependencia/ [get]
func (c *DescuentoController) GetDescuentoByPersonaPeriodoDependencia() {
	//Captura de parámetros
	idPersona := c.GetString("PersonaId")
	idDependencia := c.GetString("DependenciaId")
	idPeriodo := c.GetString("PeriodoId")
	//resultado solicitud descuento
	var resultado []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}
	//var alerta models.Alert
	var errorGetAll bool
	//alertas := append([]interface{}{"Data:"})

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento?query=Activo:true,TerceroId:"+idPersona+",PeriodoId:"+idPeriodo+",DescuentosDependenciaId.DependenciaId:"+idDependencia+"&fields=Id,TerceroId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {
			for u := 0; u < len(solicitud); u++ {
				//resultado solicitud descuento
				var descuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+
					fmt.Sprintf("%v", solicitud[u]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
				if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
					if descuento["Status"] != 404 {
						//resultado tipo descuento
						var tipo map[string]interface{}
						errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
						if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
							if tipo["Status"] != 404 {
								descuento["TipoDescuentoId"] = tipo
								solicitud[u]["DescuentosDependenciaId"] = descuento

								//resultado soporte descuento
								var soporte []map[string]interface{}
								errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=Activo:true,SolicitudDescuentoId:"+fmt.Sprintf("%v", solicitud[u]["Id"])+"&fields=DocumentoId", &soporte)
								if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
									if soporte[0]["Status"] != 404 {
										//fmt.Println("el resultado de los documentos es: ", resultado4)
										solicitud[u]["DocumentoId"] = soporte[0]["DocumentoId"]
									}
								} else {
									errorGetAll = true

									logs.Error(errSoporte)
									c.Data["message"] = "Error service GetDescuentoByPersonaPeriodoDependencia: " + errSoporte.Error()
									c.Abort("400")
								}
							}
						} else {
							errorGetAll = true

							logs.Error(errTipo)
							c.Data["message"] = "Error service GetDescuentoByPersonaPeriodoDependencia: " + errTipo.Error()
							c.Abort("400")
						}
					}
				} else {
					errorGetAll = true

					logs.Error(errDescuento)
					c.Data["message"] = "Error service GetDescuentoByPersonaPeriodoDependencia: " + errDescuento.Error()
					c.Abort("400")
				}
			}
			resultado = solicitud
			// c.Data["json"] = resultado
		} else {
			errorGetAll = true

			logs.Error("No data found")
			c.Data["message"] = "Error service GetDescuentoByPersonaPeriodoDependencia: " + "No data found"
			c.Abort("404")
		}
	} else {
		errorGetAll = true

		logs.Error(errSolicitud)
		c.Data["message"] = "Error service GetDescuentoByPersonaPeriodoDependencia: " + errSolicitud.Error()
		c.Abort("400")
	}
	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}
