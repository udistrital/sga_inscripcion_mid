package controllers

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_inscripcion_mid/utils"
	requestmanager "github.com/udistrital/sga_inscripcion_mid/utils"
	"github.com/udistrital/utils_oas/request"
)

// LegalizacionController operations for Legalizacion
type LegalizacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *LegalizacionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// PostBaseLegalizacionMatricula ...
// @Title PostBaseLegalizacionMatricula
// @Description create Legalizacion
// @Param   body        body    {}  true		"body for Legalizacion content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /base [post]
func (c *LegalizacionController) PostBaseLegalizacionMatricula() {
	var legalizacionMatriculaRequest map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &legalizacionMatriculaRequest); err == nil {
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
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{
				"Success": true, "Status": "201",
				"Message": "Created",
				"Data":    resLegalizacion,
			}
		} else {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{
				"Success": false, "Status": "400",
				"Message": "Error al crear legalización de matricula",
				"Data":    errLegalizacion,
			}
		}
	} else {
		errResponse, statusCode := requestmanager.MidResponseFormat("CreacionLegalizacionMatriculaBase", "POST", false, err.Error())
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = errResponse
	}
	c.ServeJSON()
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
		fmt.Println(newInfoComp)
		return newInfoComp, nil
	} else {
		return newInfoComp, fmt.Errorf("TercerosService Error creando información complementaria de un tercero")
	}
}

// GetOne ...
// @Title GetOne
// @Description get Legalizacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Legalizacion
// @Failure 403 :id is empty
// @router /:id [get]
func (c *LegalizacionController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Legalizacion
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Legalizacion
// @Failure 403
// @router / [get]
func (c *LegalizacionController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Legalizacion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Legalizacion	true		"body for Legalizacion content"
// @Success 200 {object} models.Legalizacion
// @Failure 403 :id is not int
// @router /:id [put]
func (c *LegalizacionController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Legalizacion
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *LegalizacionController) Delete() {

}
