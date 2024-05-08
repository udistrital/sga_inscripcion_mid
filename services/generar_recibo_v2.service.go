package services

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/barcode"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GenerarReciboV2(dataRecibo []byte) (APIResponseDTO requestresponse.APIResponse) {
	var data map[string]interface{}
	//First we fetch the data

	fmt.Println(data)

	if parseErr := json.Unmarshal(dataRecibo, &data); parseErr == nil {

		tipoRecibo := data["Tipo"].(string)

		switch tipoRecibo {
		case "Aspirante":
			if parseErr := json.Unmarshal(dataRecibo, &data); parseErr == nil {

				pdf := GenerarReciboAspiranteV2(data)

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
			fmt.Println("Admitido")
		case "Estudiante":
			fmt.Println("Estudiante")
			if parseErr := json.Unmarshal(dataRecibo, &data); parseErr == nil {
				//Then we create a new PDF document and write the title and the current date.

				pdf := GenerarEstudianteReciboV2(data)

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

		default:
			fmt.Println("La referencia de recibo es erronea")
		}

	} else {
		logs.Error(parseErr)
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, parseErr.Error())
		return APIResponseDTO
	}

	return APIResponseDTO
}

// GenerarRecibo Version Aspirante
func GenerarReciboAspiranteV2(datos map[string]interface{}) *gofpdf.Fpdf {

	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	datos["Dependencia"].(map[string]interface{})["Nombre"] = strings.ToUpper(datos["Dependencia"].(map[string]interface{})["Nombre"].(string))

	// características de página
	pdf := gofpdf.New("P", "mm", "Legal", "")
	pdf.AddPage()
	pdf.SetMargins(7, 7, 7)
	pdf.SetAutoPageBreak(true, 7) // margen inferior
	pdf.SetHomeXY()

	pdf = headerV2(pdf, numComprobante, true)
	pdf = agregarDatosAspiranteV2(pdf, datos)
	pdf = footerV2(pdf, "-COPIA ASPIRANTE-")
	pdf = separador(pdf)

	pdf = headerV2(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoAspiranteV2(pdf, datos)
	pdf = footerV2(pdf, "-COPIA PROYECTO CURRICULAR-")
	pdf = separador(pdf)

	pdf = headerV2(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoAspiranteV2(pdf, datos)
	pdf = footerV2(pdf, "-COPIA BANCO-")
	pdf = separador(pdf)

	return pdf
}

// GenerarRecibo Version Estudiante
func GenerarEstudianteReciboV2(datos map[string]interface{}) *gofpdf.Fpdf {

	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	datos["Dependencia"].(map[string]interface{})["Nombre"] = strings.ToUpper(datos["Dependencia"].(map[string]interface{})["Nombre"].(string))

	// características de página

	pdf := gofpdf.New("P", "mm", "Legal", "")
	pdf.AddPage()
	pdf.SetMargins(7, 7, 7)
	pdf.SetAutoPageBreak(true, 7) // margen inferior
	pdf.SetHomeXY()

	pdf = headerV2(pdf, numComprobante, true)
	pdf = agregarDatosEstudianteReciboV2(pdf, datos)
	pdf = footerV2(pdf, "-COPIA ESTUDIANTE-")
	pdf = separador(pdf)

	pdf = headerV2(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoEstudianteReciboV2(pdf, datos)
	pdf = footerV2(pdf, "-COPIA PROYECTO CURRICULAR-")
	pdf = separador(pdf)

	pdf = headerV2(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoEstudianteReciboV2(pdf, datos)
	pdf = footerV2(pdf, "-COPIA BANCO-")
	pdf = separador(pdf)

	return pdf
}

// Description: genera el encabezado reutilizable del recibo de pago
func headerV2(pdf *gofpdf.Fpdf, comprobante string, banco bool) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	path := beego.AppConfig.String("StaticPath")
	pdf = image(pdf, path+"/img/UDEscudo2.png", 7, pdf.GetY(), 0, 17.5)

	if banco {
		pdf = image(pdf, path+"/img/banco.PNG", 198, pdf.GetY(), 0, 12.5)
	}

	pdf.SetXY(7, pdf.GetY())
	fontStyle(pdf, "B", 10, 0)
	pdf.Cell(13, 10, "")
	pdf.Cell(140, 10, "UNIVERSIDAD DISTRITAL")
	if banco {
		fontStyle(pdf, "B", 8, 0)
		pdf.Cell(50, 10, "PAGUE UNICAMENTE EN")
		fontStyle(pdf, "B", 10, 0)
	}
	pdf.Ln(4)
	pdf.Cell(13, 10, "")
	pdf.Cell(60, 10, tr("Francisco José de Caldas"))
	pdf.Cell(80, 10, "COMPROBANTE DE PAGO No "+comprobante)

	if banco {
		fontStyle(pdf, "B", 8, 0)
		pdf.Cell(50, 10, "BANCO DE OCCIDENTE")
	} /* else {
		fontStyle(pdf, "", 8, 70)
		pdf.Cell(50, 10, "espacio para serial")
	} */

	pdf.Ln(4)
	fontStyle(pdf, "", 8, 0)
	pdf.Cell(13, 10, "")
	pdf.Cell(50, 10, "NIT 899.999.230-7")
	pdf.Ln(10)
	return pdf
}

// Description: genera el pie de paǵina reutilizable del recibo de pago
func footerV2(pdf *gofpdf.Fpdf, copiaPara string) *gofpdf.Fpdf {
	fontStyle(pdf, "", 8, 70)
	pdf.CellFormat(134, 5, copiaPara, "", 0, "C", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY())
	pdf.CellFormat(66, 5, "-Espacio para timbre o sello Banco-", "", 0, "C", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(5)

	return pdf
}

// Description: genera linea de corte reutilizable del recibo de pago
func separadorV2(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	fontStyle(pdf, "", 8, 70)
	pdf.CellFormat(201.9, 5, "...........................................................................................................................Doblar...........................................................................................................................", "", 0, "TC", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(5)
	return pdf
}

// Description: genera el código de barras reutilizable del recibo de pago
func generarCodigoBarrasV2(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	//Se genera el codigo de barras y se agrega al archivo
	documento := datos["Documento"].(string)
	for len(documento) < 12 {
		documento = "0" + documento
	}

	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	con := datos["Conceptos"].([]interface{})

	var totalValor float64

	for _, c := range con {
		conMap := c.(map[string]interface{})
		valor := conMap["Valor"].(float64)
		totalValor += valor
	}

	costo := fmt.Sprintf("%.f", totalValor)
	for len(costo) < 10 {
		costo = "0" + costo
	}
	FNC1 := '\u00f1'
	fecha := strings.Split(datos["Fecha1"].(string), "/")
	codigo := string(FNC1) + "41577099980004218020" + documento + numComprobante + string(FNC1) + "3900" + costo + string(FNC1) + "96" + fecha[2] + fecha[1] + fecha[0]
	codigoTexto := "(415)7709998000421(8020)" + documento + numComprobante + "(3900)" + costo + "(96)" + fecha[2] + fecha[1] + fecha[0]
	bcode := barcode.RegisterCode128(pdf, codigo)
	barcode.Barcode(pdf, bcode, 8, pdf.GetY()+2, 132, 12, false)
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(13.5)
	pdf.CellFormat(134, 5, codigoTexto, "", 0, "C", false, 0, "")
	pdf.Ln(5)

	return pdf
}

// Copia de recibo para aspirante (sin codigo)
func agregarDatosAspiranteV2(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 60, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Aspirante", "RB", 0, "C", false, 0, "")
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

	var cant float64 = 25
	conceptos := datos["Conceptos"].([]interface{})
	cantidadConceptos := len(conceptos)

	switch cantidadConceptos {
	case 1:
		cant = 25
	case 2:
		cant = 20
	case 3:
		cant = 15
	case 4:
		cant = 10
	case 5:
		cant = 5
	case 6:
		cant = 0
	}

	var totalValor float64

	for _, concepto := range conceptos {
		conceptoMap := concepto.(map[string]interface{})
		descripcion := dividirTexto(pdf, conceptoMap["Descripcion"].(string), 51)
		valor := conceptoMap["Valor"].(float64)

		totalValor += valor

		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(20, 5, conceptoMap["Ref"].(string), "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(50, 5, tr(descripcion[0]), "", 0, "L", false, 0, "")
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(64, 5, tr(formatoDinero(int(valor), "$", ",")+"     "), "L", 0, "R", false, 0, "")
		pdf.Ln(0)
		pdf.CellFormat(20, 10, "", "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		if len(descripcion) > 1 {
			pdf.CellFormat(50, 10, tr(descripcion[1]), "", 0, "TL", false, 0, "")
		} else {
			pdf.CellFormat(50, 10, "", "", 0, "TL", false, 0, "")
		}
		fontStyle(pdf, "", 8, 0)
		pdf.Ln(5)
	}

	valorTotal := formatoDinero(int(totalValor), "$", ",") + "     "
	valorRecargo := formatoDinero(int(totalValor*datos["Recargo"].(float64)), "$", ",") + "     "

	pdf.CellFormat(20, cant, "", "R", 0, "R", false, 0, "")
	pdf.CellFormat(50, cant, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(64, cant, "", "L", 0, "R", false, 0, "")
	pdf.Ln(cant)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha1"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorTotal), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha2"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorRecargo), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, datos["Dependencia"].(map[string]interface{})["Tipo"].(string), "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 8, 0)
	lineasProyecto := dividirTexto(pdf, datos["Dependencia"].(map[string]interface{})["Nombre"].(string), 67)
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
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	observaciones := datos["Observaciones"].([]interface{})

	for _, observacion := range observaciones {
		observacionMap := observacion.(map[string]interface{})
		ref := observacionMap["Ref"].(string)
		descripcion := observacionMap["Descripcion"].(string)

		pdf.SetXY(142.9, pdf.GetY()+4)
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(66, 4, tr(ref+" "+descripcion), "", 0, "TL", false, 0, "")
	}

	pdf.SetXY(7, ynow+65)

	return pdf
}

// Copia de recibo para estudiante (con codigo)
func agregarDatosEstudianteReciboV2(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 60, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Estudiante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr("Código"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, "Doc. Identidad", "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["Nombre"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["CodigoEstudiante"].(string)), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["Documento"].(string)), "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(20, 5, "Referencia", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(50, 5, tr("Descripción"), "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Valor", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	var cant float64 = 25
	conceptos := datos["Conceptos"].([]interface{})
	cantidadConceptos := len(conceptos)

	switch cantidadConceptos {
	case 1:
		cant = 25
	case 2:
		cant = 20
	case 3:
		cant = 15
	case 4:
		cant = 10
	case 5:
		cant = 5
	case 6:
		cant = 0
	}

	var totalValor float64

	for _, concepto := range conceptos {
		conceptoMap := concepto.(map[string]interface{})
		descripcion := dividirTexto(pdf, conceptoMap["Descripcion"].(string), 51)
		valor := conceptoMap["Valor"].(float64)

		totalValor += valor

		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(20, 5, conceptoMap["Ref"].(string), "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(50, 5, tr(descripcion[0]), "", 0, "L", false, 0, "")
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(64, 5, tr(formatoDinero(int(valor), "$", ",")+"     "), "L", 0, "R", false, 0, "")
		pdf.Ln(0)
		pdf.CellFormat(20, 10, "", "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		if len(descripcion) > 1 {
			pdf.CellFormat(50, 10, tr(descripcion[1]), "", 0, "TL", false, 0, "")
		} else {
			pdf.CellFormat(50, 10, "", "", 0, "TL", false, 0, "")
		}
		fontStyle(pdf, "", 8, 0)
		pdf.Ln(5)
	}

	valorTotal := formatoDinero(int(totalValor), "$", ",") + "     "
	valorRecargo := formatoDinero(int(totalValor*datos["Recargo"].(float64)), "$", ",") + "     "

	pdf.CellFormat(20, cant, "", "R", 0, "R", false, 0, "")
	pdf.CellFormat(50, cant, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(64, cant, "", "L", 0, "R", false, 0, "")
	pdf.Ln(cant)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha1"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorTotal), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha2"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorRecargo), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, datos["Dependencia"].(map[string]interface{})["Tipo"].(string), "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 8, 0)
	lineasProyecto := dividirTexto(pdf, datos["Dependencia"].(map[string]interface{})["Nombre"].(string), 67)
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
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")

	observaciones := datos["Observaciones"].([]interface{})

	for _, observacion := range observaciones {
		observacionMap := observacion.(map[string]interface{})
		ref := observacionMap["Ref"].(string)
		descripcion := observacionMap["Descripcion"].(string)

		pdf.SetXY(142.9, pdf.GetY()+4)
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(66, 4, tr(ref+" "+descripcion), "", 0, "TL", false, 0, "")
	}

	pdf.SetXY(7, ynow+65)

	return pdf
}

// Copia de recibo version aspirante (sin codigo)
func agregarDatosCopiaBancoProyectoAspiranteV2(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 20, 2.5, "1234", "")

	con := datos["Conceptos"].([]interface{})

	var totalValor float64

	for _, c := range con {
		conMap := c.(map[string]interface{})
		valor := conMap["Valor"].(float64)
		totalValor += valor
	}

	valorTotal := formatoDinero(int(totalValor), "$", ",") + "     "
	valorRecargo := formatoDinero(int(totalValor*datos["Recargo"].(float64)), "$", ",") + "     "

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Aspirante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Documento de Identidad", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["Nombre"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(64, 5, tr(datos["Documento"].(string)), "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha1"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorTotal)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(10)

	datos["Documento"] = datos["Documento"]
	pdf = generarCodigoBarrasV2(pdf, datos)
	pdf.Ln(2)

	pdf.RoundedRect(7, pdf.GetY(), 134, 10, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "R", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "R", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha2"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorRecargo)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(10)

	pdf = generarCodigoBarrasV2(pdf, datos)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, datos["Dependencia"].(map[string]interface{})["Tipo"].(string), "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 8, 0)
	lineasProyecto := dividirTextoV2(pdf, datos["Dependencia"].(map[string]interface{})["Nombre"].(string), 67)
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
	pdf.CellFormat(36, 5, fechaActual(), "R", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, datos["Periodo"].(string), "", 0, "C", false, 0, "")

	pdf.RoundedRect(142.9, ynow, 66, alturaRecuadro, 2.5, "1234", "")

	pdf.RoundedRect(142.9, pdf.GetY()+8, 66, 35, 2.5, "1234", "")

	pdf.SetXY(142.9, pdf.GetY()+8)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(9.5, 5, "Ref.", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(37, 5, tr("Descripción"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(19.5, 5, "Valor", "LB", 0, "C", false, 0, "")

	var cant float64 = 25
	conceptos := datos["Conceptos"].([]interface{})
	cantidadConceptos := len(conceptos)

	switch cantidadConceptos {
	case 1:
		cant = 25
	case 2:
		cant = 20
	case 3:
		cant = 15
	case 4:
		cant = 10
	case 5:
		cant = 5
	case 6:
		cant = 0
	}

	pdf.SetXY(142.9, pdf.GetY()+5)
	for _, concepto := range conceptos {
		pdf.SetXY(142.9, pdf.GetY())
		conceptoMap := concepto.(map[string]interface{})
		descripcion := dividirTexto(pdf, conceptoMap["Descripcion"].(string), 51)
		valor := conceptoMap["Valor"].(float64)

		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(9.5, 5, conceptoMap["Ref"].(string), "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(37, 5, tr(descripcion[0]), "R", 0, "R", false, 0, "")
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(19.5, 5, tr(formatoDinero(int(valor), "$", ",")+"     "), "R", 0, "R", false, 0, "")
		pdf.Ln(0)
		//pdf.CellFormat(9.5, 7, "", "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		if len(descripcion) > 1 {
			pdf.CellFormat(37, 7, tr(descripcion[1]), "", 0, "TL", false, 0, "")
		} else {
			pdf.CellFormat(37, 7, "", "", 0, "TL", false, 0, "")
		}
		fontStyle(pdf, "", 8, 0)
		pdf.Ln(5)
		pdf.SetXY(142.9, pdf.GetY())
	}
	pdf.CellFormat(9.5, cant, "", "R", 0, "R", false, 0, "")
	pdf.CellFormat(37, cant, "", "R", 0, "R", false, 0, "")
	pdf.CellFormat(19.5, cant, "", "R", 0, "R", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+cant)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	observaciones := datos["Observaciones"].([]interface{})

	for _, observacion := range observaciones {
		observacionMap := observacion.(map[string]interface{})
		ref := observacionMap["Ref"].(string)
		descripcion := observacionMap["Descripcion"].(string)

		pdf.SetXY(142.9, pdf.GetY()+4)
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(66, 4, tr(ref+" "+descripcion), "", 0, "TL", false, 0, "")
	}

	pdf.SetXY(7, ynow+88)

	return pdf
}

// Copia de recibo version estudiante (con codigo)
func agregarDatosCopiaBancoProyectoEstudianteReciboV2(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	con := datos["Conceptos"].([]interface{})

	var totalValor float64

	for _, c := range con {
		conMap := c.(map[string]interface{})
		valor := conMap["Valor"].(float64)
		totalValor += valor
	}

	valorTotal := formatoDinero(int(totalValor), "$", ",") + "     "
	valorRecargo := formatoDinero(int(totalValor*datos["Recargo"].(float64)), "$", ",") + "     "

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 20, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Estudiante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr("Código"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, "Doc. Identidad", "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["Nombre"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["CodigoEstudiante"].(string)), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["Documento"].(string)), "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha1"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorTotal)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	datos["Documento"] = datos["Documento"]
	pdf = generarCodigoBarrasV2(pdf, datos)
	pdf.Ln(2)

	pdf.RoundedRect(7, pdf.GetY(), 134, 10, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "R", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "R", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha2"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorRecargo)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf = generarCodigoBarrasV2(pdf, datos)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, datos["Dependencia"].(map[string]interface{})["Tipo"].(string), "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 8, 0)
	lineasProyecto := dividirTexto(pdf, datos["Dependencia"].(map[string]interface{})["Nombre"].(string), 67)
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
	pdf.CellFormat(36, 5, fechaActual(), "R", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, datos["Periodo"].(string), "", 0, "C", false, 0, "")

	pdf.RoundedRect(142.9, ynow, 66, alturaRecuadro, 2.5, "1234", "")

	pdf.RoundedRect(142.9, pdf.GetY()+8, 66, 35, 2.5, "1234", "")

	pdf.SetXY(142.9, pdf.GetY()+8)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(9.5, 5, "Ref.", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(37, 5, tr("Descripción"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(19.5, 5, "Valor", "LB", 0, "C", false, 0, "")

	var cant float64 = 25
	conceptos := datos["Conceptos"].([]interface{})
	cantidadConceptos := len(conceptos)

	switch cantidadConceptos {
	case 1:
		cant = 25
	case 2:
		cant = 20
	case 3:
		cant = 15
	case 4:
		cant = 10
	case 5:
		cant = 5
	case 6:
		cant = 0
	}

	pdf.SetXY(142.9, pdf.GetY()+5)
	for _, concepto := range conceptos {
		pdf.SetXY(142.9, pdf.GetY())
		conceptoMap := concepto.(map[string]interface{})
		descripcion := dividirTexto(pdf, conceptoMap["Descripcion"].(string), 51)
		valor := conceptoMap["Valor"].(float64)

		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(9.5, 5, conceptoMap["Ref"].(string), "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		pdf.CellFormat(37, 5, tr(descripcion[0]), "R", 0, "R", false, 0, "")
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(19.5, 5, tr(formatoDinero(int(valor), "$", ",")+"     "), "R", 0, "R", false, 0, "")
		pdf.Ln(0)
		//pdf.CellFormat(9.5, 7, "", "R", 0, "R", false, 0, "")
		fontStyle(pdf, "B", 8, 0)
		if len(descripcion) > 1 {
			pdf.CellFormat(37, 7, tr(descripcion[1]), "", 0, "TL", false, 0, "")
		} else {
			pdf.CellFormat(37, 7, "", "", 0, "TL", false, 0, "")
		}
		fontStyle(pdf, "", 8, 0)
		pdf.Ln(5)
		pdf.SetXY(142.9, pdf.GetY())
	}
	pdf.CellFormat(9.5, cant, "", "R", 0, "R", false, 0, "")
	pdf.CellFormat(37, cant, "", "R", 0, "R", false, 0, "")
	pdf.CellFormat(19.5, cant, "", "R", 0, "R", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+cant)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	observaciones := datos["Observaciones"].([]interface{})

	for _, observacion := range observaciones {
		observacionMap := observacion.(map[string]interface{})
		ref := observacionMap["Ref"].(string)
		descripcion := observacionMap["Descripcion"].(string)

		pdf.SetXY(142.9, pdf.GetY()+4)
		fontStyle(pdf, "", 8, 0)
		pdf.CellFormat(66, 4, tr(ref+" "+descripcion), "", 0, "TL", false, 0, "")
	}

	pdf.SetXY(7, ynow+88)

	return pdf
}

// agrega imagen de archivo a pdf, w o h en cero autoajusta segun ratio imagen
func imageV2(pdf *gofpdf.Fpdf, image string, x, y, w, h float64) *gofpdf.Fpdf {
	//The ImageOptions method takes a file path, x, y, width, and height parameters, and an ImageOptions struct to specify a couple of options.
	pdf.ImageOptions(image, x, y, w, h, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

// convierte pdf a base64
func encodePDFV2(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("../docs/recibo.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

// Fecha de expedición del recibo
func fechaActualV2() string {
	hoy := time.Now()
	return fmt.Sprintf("%02d/%02d/%d", hoy.Day(), hoy.Month(), hoy.Year())
}

// Estilo de fuente usando Helvetica
func fontStyleV2(pdf *gofpdf.Fpdf, style string, size float64, bw int) {
	pdf.SetTextColor(bw, bw, bw)
	pdf.SetFont("Helvetica", style, size)
}

// Divide texto largo en lineas
// func dividirTexto(pdf *gofpdf.Fpdf, text string, w float64) []string {
// 	fmt.Println("Texto: ", text)
// 	lineasraw := pdf.SplitLines([]byte(text), w)
// 	var lineas []string
// 	for _, lineraw := range lineasraw {
// 		lineas = append(lineas, string(lineraw))
// 	}
// 	return lineas
// }

func dividirTextoV2(pdf *gofpdf.Fpdf, text string, w float64) []string {
	palabras := strings.Fields(text)
	var lineas []string
	var lineaActual string

	for _, palabra := range palabras {
		pruebaLinea := lineaActual
		if len(pruebaLinea) > 0 {
			pruebaLinea += " "
		}
		pruebaLinea += palabra

		// Calcula el ancho de la línea con la palabra añadida
		anchoLinea := pdf.GetStringWidth(pruebaLinea)

		if anchoLinea > w && len(lineaActual) > 0 {
			// Si la línea excede el ancho permitido, guarda la línea actual y comienza una nueva
			lineas = append(lineas, lineaActual)
			lineaActual = palabra
		} else {
			// Si la línea no excede el ancho, añade la palabra a la línea actual
			lineaActual = pruebaLinea
		}
	}

	// Añade la última línea si queda alguna palabra
	if len(lineaActual) > 0 {
		lineas = append(lineas, lineaActual)
	}

	return lineas
}

func formatoDineroV2(valor int, simbolo string, separador string, valorStr ...string) string {
	if simbolo != "" {
		simbolo = simbolo + " "
	}
	var caracteres []string
	if valor > 0 {
		caracteres = strings.Split(fmt.Sprintf("%d", valor), "")
	} else {
		caracteres = strings.Split(valorStr[0], "")
	}

	valorTexto := ""

	for i := len(caracteres) - 1; i >= 0; i-- {
		sep := ((i % 3) == 0) && (i > 0)
		valorTexto += caracteres[len(caracteres)-1-i]
		if sep {
			valorTexto += separador
		}
	}

	return simbolo + valorTexto
}
