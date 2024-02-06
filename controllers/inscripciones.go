package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid_inscripcion/helpers"
	"github.com/udistrital/sga_mid_inscripcion/models"
	"github.com/udistrital/sga_mid_inscripcion/services"
	"github.com/udistrital/sga_mid_inscripcion/utils"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

type InscripcionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *InscripcionesController) URLMapping() {
	c.Mapping("PostInformacionFamiliar", c.PostInformacionFamiliar)
	c.Mapping("PostReintegro", c.PostReintegro)
	c.Mapping("PostTransferencia", c.PostTransferencia)
	c.Mapping("PostInfoIcfesColegio", c.PostInfoIcfesColegio)
	c.Mapping("PostPreinscripcion", c.PostPreinscripcion)
	c.Mapping("PostInfoComplementariaUniversidad", c.PostInfoComplementariaUniversidad)
	c.Mapping("PostInfoComplementariaTercero", c.PostInfoComplementariaTercero)
	c.Mapping("GetInfoComplementariaTercero", c.GetInfoComplementariaTercero)
	c.Mapping("PostInfoIcfesColegioNuevo", c.PostInfoIcfesColegioNuevo)
	c.Mapping("ConsultarProyectosEventos", c.ConsultarProyectosEventos)
	c.Mapping("ActualizarInfoContacto", c.ActualizarInfoContacto)
	c.Mapping("GetEstadoInscripcion", c.GetEstadoInscripcion)
	c.Mapping("PostGenerarInscripcion", c.PostGenerarInscripcion)
}

// GetEstadoInscripcion ...
// @Title GetEstadoInscripcion
// @Description consultar los estados de todos los recibos generados por el tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Param	id_periodo	path	int	true	"Id del ultimo periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /estado_recibos/:persona_id/:id_periodo [get]
func (c *InscripcionesController) GetEstadoInscripcion() {

	defer errorhandler.HandlePanic(&c.Controller)

	terceroId := c.GetString("Id")

	respuesta := services.GetFormacionAcademicaByIdTercero(terceroId)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInformacionFamiliar ...
// @Title PostInformacionFamiliar
// @Description Agregar Información Familiar
// @Param   body        body    {}  true        "body Agregar PostInformacionFamiliar content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_informacion_familiar [post]
func (c *InscripcionesController) PostInformacionFamiliar() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InformacionFamiliar(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostReintegro ...
// @Title PostReintegro
// @Description Agregar Reintegro
// @Param   body        body    {}  true        "body Agregar Reintegro content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_reintegro [post]
func (c *InscripcionesController) PostReintegro() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.Reintegro(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostTransferencia ...
// @Title PostTransferencia
// @Description Agregar Transferencia
// @Param   body        body    {}  true        "body Agregar Transferencia content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_transferencia [post]
func (c *InscripcionesController) PostTransferencia() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.TransferenciaPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInfoIcfesColegio ...
// @Title PostInfoIcfesColegio
// @Description Agregar InfoIcfesColegio
// @Param   body        body    {}  true        "body Agregar InfoIcfesColegio content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_info_icfes_colegio [post]
func (c *InscripcionesController) PostInfoIcfesColegio() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InfoIcfesColegio(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostPreinscripcion ...
// @Title PostPreinscripcion
// @Description Agregar Preinscripcion
// @Param   body        body    {}  true        "body Agregar Preinscripcion content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_preinscripcion [post]
func (c *InscripcionesController) PostPreinscripcion() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.PreinscripcionPost(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInfoIcfesColegioNuevo ...
// @Title PostInfoIcfesColegioNuevo
// @Description Agregar InfoIcfesColegio
// @Param   body        body    {}  true        "body Agregar InfoIcfesColegio content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_info_icfes_colegio_nuevo [post]
func (c *InscripcionesController) PostInfoIcfesColegioNuevo() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InfoNuevoColegioIcfes(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// PostInfoComplementariaUniversidad ...
// @Title PostInfoComplementariaUniversidad
// @Description Agregar InfoComplementariaUniversidad
// @Param   body        body    {}  true        "body Agregar InfoComplementariaUniversidad content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info_complementaria_universidad [post]
func (c *InscripcionesController) PostInfoComplementariaUniversidad() {

	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.PutInfoComplementaria(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// ConsultarProyectosEventos ...
// @Title ConsultarProyectosEventos
// @Description get ConsultarProyectosEventos by id
// @Param	evento_padre_id	path	int	true	"Id del Evento Padre"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_proyectos_eventos/:evento_padre_id [get]
func (c *InscripcionesController) ConsultarProyectosEventos() {
	
	defer errorhandler.HandlePanic(&c.Controller)
	
	//Id de la persona
	idStr := c.Ctx.Input.Param(":evento_padre_id")
	fmt.Println("El id es: " + idStr)

	respuesta := services.ConsultarEventos(idStr)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PostInfoComplementariaTercero ...
// @Title PostInfoComplementariaTercero
// @Description Agregar PostInfoComplementariaTercero
// @Param   body        body    {}  true        "body Agregar PostInfoComplementariaTercero content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info_complementaria_tercero [post]
func (c *InscripcionesController) PostInfoComplementariaTercero() {
	
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.InfoComplementariaTercero(data)

	c.Ctx.Output.SetStatus(respuesta.Status)

	c.Data["json"] = respuesta

	c.ServeJSON()
}

// GetInfoComplementariaTercero ...
// @Title GetInfoComplementariaTercero
// @Description consultar la información complementaria del tercero
// @Success 200 {}
// @Failure 404 not found resource
// @router /info_complementaria_tercero/:persona_id [get]
func (c *InscripcionesController) GetInfoComplementariaTercero() {
	//Id de la persona
	persona_id := c.Ctx.Input.Param(":persona_id")
	//resultado consulta
	resultado := map[string]interface{}{}
	// var resultado map[string]interface{}
	var errorGetAll bool
	var alerta models.Alert
	alertas := []interface{}{}

	// 41 = estrato
	IdEstrato, _ := helpers.IdInfoCompTercero("9", "ESTRATO")
	var resultadoEstrato []map[string]interface{}
	errEstratoResidencia := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdEstrato+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoEstrato)
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
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				alertas = append(alertas, errEstratoResidencia)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errEstratoResidencia)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// 55 = codigo postal
	IdCodPostal, _ := helpers.IdInfoCompTercero("10", "CODIGO_POSTAL")
	var resultadoCodigoPostal []map[string]interface{}
	errCodigoPostal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCodPostal+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoCodigoPostal)
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
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				alertas = append(alertas, errCodigoPostal)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errCodigoPostal)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// 51 = telefono
	IdTelefono, _ := helpers.IdInfoCompTercero("10", "TELEFONO")
	var resultadoTelefono []map[string]interface{}
	errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdTelefono+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoTelefono)
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
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				errorGetAll = true
				alertas = append(alertas, errTelefono)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errTelefono)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// 54 = direccion
	IdDireccion, _ := helpers.IdInfoCompTercero("10", "DIRECCIÓN")
	var resultadoDireccion []map[string]interface{}
	errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdDireccion+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoDireccion)
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
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				alertas = append(alertas, errDireccion)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errDireccion)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// Correo registro
	IdCorreo, _ := helpers.IdInfoCompTercero("10", "CORREO")
	var resultadoCorreo []map[string]interface{}
	errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCorreo+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoCorreo)
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
	errCorreoAlterno := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCorreoAlterno+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoCorreoAlterno)
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
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// ActualizarInfoContacto ...
// @Title ActualizarInfoContacto
// @Description Actualiza los datos de contacto del tercero
// @Param	body	body 	{}	true		"body for Actualizar la info de contacto del tercero content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info_contacto [put]
func (c *InscripcionesController) ActualizarInfoContacto() {
	var InfoContacto map[string]interface{}

	var alerta models.Alert
	alertas := []interface{}{}
	var algoFallo bool = false

	var revertPuts []map[string]interface{}
	var inactivePosts []map[string]interface{}

	var respuestas []interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoContacto); err == nil {
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
						alertas = append(alertas, errPutInfoComp.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
					}
				} else {
					algoFallo = true
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					alerta.Body = alertas
				}
			} else {
				var resp map[string]interface{}
				errPostInfoComp := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resp, InfoComplementaria)
				if errPostInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
					respuestas = append(respuestas, resp)
					inactivePosts = append(inactivePosts, resp)
				} else {
					algoFallo = true
					alertas = append(alertas, errPostInfoComp.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
				}
			}
			if algoFallo {
				break
			}
		}
	} else {
		algoFallo = true
		alertas = append(alertas, err.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
	}

	if !algoFallo {
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = respuestas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	} else {
		for _, revert := range revertPuts {
			var resp map[string]interface{}
			request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", revert["Id"].(float64)), "PUT", &resp, revert)
		}
		for _, disable := range inactivePosts {
			helpers.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + fmt.Sprintf("%.f", disable["Id"].(float64)))
		}
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// PostGenerarInscripcion ...
// @Title PostGenerarInscripcion
// @Description Registra una nueva inscripción con su respectivo recibo de pago
// @Param	body	body 	{}	true		"body for información de suministrada por el usuario par la inscripción"
// @Success 200 {}
// @Failure 403 body is empty
// @router /generar_inscripcion [post]
func (c *InscripcionesController) PostGenerarInscripcion() {
	var reciboVencido bool
	var respuesta models.Alert
	var SolicitudInscripcion map[string]interface{}
	var TipoParametro string
	var parametro map[string]interface{}
	var Valor map[string]interface{}
	var NuevoRecibo map[string]interface{}
	var inscripcionRealizada map[string]interface{}
	var contadorRecibos int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudInscripcion); err == nil {
		objTransaccion := map[string]interface{}{
			"codigo":              SolicitudInscripcion["Id"].(float64),
			"nombre":              SolicitudInscripcion["Nombre"].(string),
			"apellido":            SolicitudInscripcion["Apellido"].(string),
			"correo":              SolicitudInscripcion["Correo"].(string),
			"proyecto":            SolicitudInscripcion["ProgramaAcademicoId"].(float64),
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
		} else if SolicitudInscripcion["Nivel"].(float64) == 2 {
			TipoParametro = "12"
		}

		persona_id := fmt.Sprintf("%d", int(SolicitudInscripcion["PersonaId"].(float64)))
		id_periodo := fmt.Sprintf("%d", int(SolicitudInscripcion["PeriodoId"].(float64)))
		id_programa_academico := fmt.Sprintf("%d", int(SolicitudInscripcion["ProgramaAcademicoId"].(float64)))

		recibosResultado, err := helpers.VerificarRecibos(persona_id, id_periodo)

		if err == "" {
			if inscripciones, ok := recibosResultado["Inscripciones"]; ok {
				// Convertir la variable de tipo interface{} a un slice de mapas
				inscripcionesMap, ok := inscripciones.([]map[string]interface{})
				if len(inscripcionesMap) > 0 && ok {
					for i := 0; i < len(inscripcionesMap); i++ {
						id_programa_inscripciones := fmt.Sprintf("%d", int(inscripcionesMap[i]["ProgramaAcademicoId"].(float64)))
						estado_recibo_inscripciones := inscripcionesMap[i]["Estado"].(string)
						if id_programa_inscripciones == id_programa_academico {
							if estado_recibo_inscripciones == "Vencido" {
								reciboVencido = true
							}else {	
								reciboVencido = false
							}
						}else{
							contadorRecibos++
						}
					}
					if contadorRecibos == len(inscripcionesMap){
						reciboVencido = true
					}
				}

			}

			//Verificar si existe un recibo vencido o es la primera vez que inscribe el postgrado
			if reciboVencido || fmt.Sprintf("%v", recibosResultado) == "map[]" {
				errInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion", "POST", &inscripcionRealizada, inscripcion)
				if errInscripcion == nil && inscripcionRealizada["Status"] != "400" {
					errParam := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=Activo:true,ParametroId.TipoParametroId.Id:2,ParametroId.CodigoAbreviacion:"+TipoParametro+",PeriodoId.Year:"+fmt.Sprintf("%v", objTransaccion["aniopago"])+",PeriodoId.CodigoAbreviacion:VG", &parametro)
					if errParam == nil && fmt.Sprintf("%v", parametro["Data"].([]interface{})[0]) != "map[]" {
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
									respuesta.Type = "success"
									respuesta.Code = "200"
									respuesta.Body = inscripcionUpdate

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
									respuesta.Type = "error"
									respuesta.Code = "400"
									respuesta.Body = errInscripcionUpdate.Error()
								}
							} else {
								//var resDelete string
								//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
								helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
								logs.Error(errRecibo)
								respuesta.Type = "error"
								respuesta.Code = "400"
								respuesta.Body = errRecibo.Error()
							}
						} else {
							//var resDelete string
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
							helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
							logs.Error(errJson)
							respuesta.Type = "error"
							respuesta.Code = "403"
							respuesta.Body = errJson.Error()
						}
					} else {
						//var resDelete string
						//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
						helpers.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
						logs.Error(errParam)
						respuesta.Type = "error"
						respuesta.Code = "400"
						respuesta.Body = errParam.Error()
					}

				} else {
					logs.Error(errInscripcion)
					respuesta.Type = "success"
					respuesta.Code = "204"
					respuesta.Body = errInscripcion.Error()
				}
			} else {
				respuesta.Type = "success"
				respuesta.Code = "204"
				respuesta.Body = "Recipe already exist"
			}

		} else if err == "400" {
			respuesta.Code = "400"
			respuesta.Type = "error"
			respuesta.Body = "Bad request"
			c.Data["json"] = map[string]interface{}{"Response": "Bad request"}
		} else {
			respuesta.Code = "404"
			respuesta.Type = "error"
			respuesta.Body = "No data found"
			c.Data["json"] = map[string]interface{}{"Response": "No data found"}
		}

	} else {
		logs.Error(err)
		respuesta.Type = "error"
		respuesta.Code = "403"
		respuesta.Body = err.Error()
	}

	c.Data["json"] = respuesta
	c.ServeJSON()

}
