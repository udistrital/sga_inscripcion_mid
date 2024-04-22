package services

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/utils_oas/requestresponse"
)

/*
Data recibida
{
   "Nombre":"Prueba Aspirante uno ",
   "Documento":"1234456",
   "Periodo":"2024-1",
   "Proyecto":"Doctorado Estudio Artisticos",
   "Comprobante":"8702",
   "Fecha_pago":"15/03/2024",
   "Descripcion":"Matricula",
   "ValorMatricula":154700
}
*/
func GenerarReciboLiquidacionPost(dataRecibo []byte) (APIResponseDTO requestresponse.APIResponse) {
	var data map[string]interface{}

	if parseErr := json.Unmarshal(dataRecibo, &data); parseErr == nil {

		pdf := GenerarReciboLiquidacion(data)

		if pdf.Err() {
			logs.Error("Failed creating PDF report: %s\n", pdf.Error())
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, pdf.Error())
		}

		if pdf.Ok() {
			encodedFile := encodePDF(pdf)
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, encodedFile, nil)
		}

	} else {
		logs.Error(parseErr)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, parseErr.Error())
		return APIResponseDTO
	}
	return APIResponseDTO
}

func GenerarReciboLiquidacion(datos map[string]interface{}) *gofpdf.Fpdf {

	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	datos["Proyecto"] = strings.ToUpper(datos["Proyecto"].(string))
	datos["Descripcion"] = strings.ToUpper(datos["Descripcion"].(string))

	// características de página
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetMargins(7, 7, 7)
	pdf.SetAutoPageBreak(true, 7) // margen inferior
	pdf.SetHomeXY()

	pdf = header(pdf, numComprobante, true)
	pdf = agregarDatosEstudiante(pdf, datos)
	pdf = footer(pdf, "-COPIA ESTUDIANTE-")
	pdf = separador(pdf)

	pdf = header(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoAspirante(pdf, datos)
	pdf = footer(pdf, "-COPIA PROYECTO CURRICULAR-")
	pdf = separador(pdf)

	pdf = header(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoAspirante(pdf, datos)
	pdf = footer(pdf, "-COPIA BANCO-")
	//pdf = separador(pdf)

	return pdf
}

func agregarDatosEstudiante(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	valorDerecho := formatoDinero(int(datos["ValorMatricula"].(float64)), "$", ",") + "     "

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 45, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Estudiante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Documento de Identidad", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["Nombre"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(64, 5, tr(datos["Documento"].(string)), "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(20, 5, "Referencia", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(50, 5, tr("Descripción"), "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Valor", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(20, 5, "1     ", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "B", 7.5, 0)
	descripcion := dividirTexto(pdf, datos["Descripcion"].(string), 51)
	pdf.CellFormat(50, 5, tr(descripcion[0]), "", 0, "L", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(64, 5, tr(valorDerecho), "L", 0, "R", false, 0, "")
	pdf.Ln(5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(20, 5, "2     ", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "B", 7.5, 0)
	pdf.CellFormat(50, 5, "SEGURO", "", 0, "L", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(64, 5, tr(valorDerecho), "L", 0, "R", false, 0, "")
	pdf.Ln(5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(20, 5, "3     ", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "B", 7.5, 0)
	pdf.CellFormat(50, 5, "CARNET", "", 0, "L", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(64, 5, tr(valorDerecho), "L", 0, "R", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, "Proyecto Curricular", "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 7.5, 0)
	lineasProyecto := dividirTexto(pdf, datos["ProyectoAspirante"].(string), 67)
	var alturaRecuadro float64 = 20

	pdf.SetXY(142.9, pdf.GetY()+5)
	pdf.CellFormat(66, 5, tr(lineasProyecto[0]), "", 0, "L", false, 0, "")

	if len(lineasProyecto) > 1 {
		pdf.SetXY(142.9, pdf.GetY()+5)
		pdf.CellFormat(66, 5, tr(lineasProyecto[1]), "", 0, "L", false, 0, "")
		alturaRecuadro = 25
	}

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(36, 5, tr("Fecha de Expedición"), "TRB", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, "Periodo", "TB", 0, "C", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(36, 5, tr(fechaActual()), "R", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, tr(datos["Periodo"].(string)), "", 0, "C", false, 0, "")

	pdf.RoundedRect(142.9, ynow, 66, alturaRecuadro, 2.5, "1234", "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 7, 70)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+4)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, tr(datos["Descripcion"].(string)), "", 0, "TL", false, 0, "")

	pdf.SetXY(7, ynow+45)

	return pdf
}
