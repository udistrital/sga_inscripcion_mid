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
	"github.com/udistrital/sga_inscripcion_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GenerarReciboV2(dataRecibo []byte) (APIResponseDTO requestresponse.APIResponse) {
	var data map[string]interface{}
	//First we fetch the data

	fmt.Println(data)

	if parseErr := json.Unmarshal(dataRecibo, &data); parseErr == nil {

		tipoRecibo := data["Tipo"].(string)

		switch tipoRecibo {
		case "Inscripcion":
			if parseErr := json.Unmarshal(dataRecibo, &data); parseErr == nil {

				var ReciboXML map[string]interface{}
				ReciboInscripcion := data["INSCRIPCION"].(map[string]interface{})["idRecibo"].(string)
				if ReciboInscripcion != "0/<nil>" {
					//errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ConsultarReciboJbpmService")+"consulta_recibo/"+ReciboInscripcion, &ReciboXML)
					errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ConsultarReciboJbpmService")+"consulta_recibo/8702/2021", &ReciboXML)
					fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
					fmt.Println(ReciboXML)
					fmt.Println("http://" + beego.AppConfig.String("ConsultarReciboJbpmService") + "consulta_recibo/8702/2021")
					if errRecibo == nil {
						if ReciboXML != nil && fmt.Sprintf("%v", ReciboXML) != "map[reciboCollection:map[]]" && fmt.Sprintf("%v", ReciboXML) != "map[]" {

							data["Valor"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["valor_extraordinario"].(string)
							//fmt.Println(ReciboXML)

							if fecha, exist := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_pagado"].(string); exist {
								if fecha != "" {
									data["fechaExiste"] = true
									data["fecha1"] = fecha
								} else {
									data["fechaExiste"] = false
								}
							} else {
								data["fechaExiste"] = false
							}

							if !data["fechaExiste"].(bool) {
								data["fecha1"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha"].(string)
							}

							data["Comprobante"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["secuencia"].(string)
							data["Estado"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["pago"].(string)

							pdf := generarComprobanteInscripcionV2(data)

							if pdf.Err() {
								logs.Error("Failed creating PDF voucher: %s\n", pdf.Error())
								APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, pdf.Error())
							}
							//fmt.Println(data)
							if pdf.Ok() {
								encodedFile := encodePDF(pdf)
								APIResponseDTO = requestresponse.APIResponseDTO(true, 200, encodedFile, nil)
								fecha_actual := time.Now()
								dataEmail := map[string]interface{}{
									"dia":     fecha_actual.Day(),
									"mes":     utils.GetNombreMes(fecha_actual.Month()),
									"anio":    fecha_actual.Year(),
									"nombre":  data["Nombre"].(string),
									"periodo": data["Periodo"].(string),
								}
								fmt.Println("data object", dataEmail)
								//utils.SendNotificationInscripcionSolicitud(dataEmail, objTransaccion["correo"].(string))
								attachments := []map[string]interface{}{}
								attachments = append(attachments, map[string]interface{}{
									"ContentType": "application/pdf",
									"FileName":    "Comprobante_inscripcion_" + data["Dependencia"].(map[string]interface{})["Nombre"].(string),
									"Base64File":  encodedFile,
								})
								utils.SendNotificationInscripcionComprobante(dataEmail, data["Correo"].(string), attachments)
							}

						} else {
							logs.Error("reciboCollection seems empty", ReciboXML)
							APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "ReciboCollection seems empty")
							return APIResponseDTO
						}
					} else {
						logs.Error(errRecibo)
						APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errRecibo.Error())
						return APIResponseDTO
					}
				} else {
					logs.Error("ReciboInscripcionId seems empty")
					APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, "ReciboInscripcionId seems empty")
					return APIResponseDTO
				}

			} else {
				logs.Error(parseErr)
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, parseErr.Error())
				return APIResponseDTO
			}
			fmt.Println("Inscripción")
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

func generarComprobanteInscripcionV2(data map[string]interface{}) *gofpdf.Fpdf {

	// características de página
	pdf := gofpdf.New("P", "mm", "Letter", "") //215.9 279.4

	// pps page properties and styling
	pps := styling{mL: 7, mT: 7, mR: 7, mB: 7, hF: 10}

	pps.wW, pps.hW = pdf.GetPageSize()
	pps.wW -= (pps.mL + pps.mR)
	pps.hW -= (pps.mT + pps.mB)

	pdf.SetMargins(pps.mL, pps.mT, pps.mR)
	pdf.SetAutoPageBreak(true, pps.mB+pps.hF) // margen inferior

	pdf.SetHeaderFunc(headerComprobanteV2(pdf, data, pps))

	pdf.SetFooterFunc(footerComprobanteV2(pdf, pps))
	pdf.AddPage()
	//pdf.Rect(pps.mL, pps.mT, pps.wW, pps.hW, "")

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.CellFormat(pps.wW*0.5, 5, tr(fmt.Sprintf("Inscripción No. %.f", data["INSCRIPCION"].(map[string]interface{})["id"].(float64))), "", 0, "C", false, 0, "")
	pdf.CellFormat(pps.wW*0.5, 5, tr(data["INSCRIPCION"].(map[string]interface{})["fechaInsripcion"].(string)), "", 0, "C", false, 0, "")
	pdf.Ln(9)

	informacionPersonalV2(pdf, data, pps)
	pdf.Ln(7)
	informacionPagoV2(pdf, data, pps)
	pdf.Ln(7)
	documentacionSuministradaV2(pdf, data, pps)

	return pdf
}

func informacionPersonalV2(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetDrawColor(100, 100, 100)
	pdf.RoundedRect(pps.mL, pdf.GetY()-1, pps.wW, 31, 2.5, "1234", "")
	pdf.Cell(pps.wW*0.01, 8, "")
	pdf.SetFillColor(0, 162, 255)
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+1, pps.wW*.98, 6, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pps.wW*0.98, 8, tr("INFORMACIÓN PERSONAL"), "", 0, "C", false, 0, "")
	pdf.Ln(8)

	pdf.SetFontStyle("")
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Nombre:", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["Nombre"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Tipo documento: ", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["TipoDocumento"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Número documento: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["Documento"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Teléfono contacto: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["Telefono"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Programa inscribe: ", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["Dependencia"].(map[string]interface{})["Nombre"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Correo contacto: ", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["Correo"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Énfasis: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["Dependencia"].(map[string]interface{})["Enfasis"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Periodo académico: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["Periodo"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	return pdf
}

func informacionPagoV2(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	textFecha := "Fecha de pago:"
	if !data["fechaExiste"].(bool) {
		textFecha = "Fecha de generación:"
		data["fecha1"] = strings.Split(data["fecha1"].(string), "T")[0]
		orderFecha := strings.Split(data["fecha1"].(string), "-")
		data["fecha1"] = fmt.Sprintf("%s/%s/%s", orderFecha[2], orderFecha[1], orderFecha[0])
	}

	estado := data["Estado"].(string)
	if estado == "S" {
		estado = "Pagado"
	} else if estado == "N" {
		estado = "Pendiente pago"
	} else if estado == "V" {
		estado = "Vencido"
	}
	data["Estado"] = estado

	pdf.SetDrawColor(100, 100, 100)
	pdf.RoundedRect(pps.mL, pdf.GetY()-1, pps.wW, 21, 2.5, "1234", "")
	pdf.Cell(pps.wW*0.01, 8, "")
	pdf.SetFillColor(0, 162, 255)
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+1, pps.wW*.98, 6, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pps.wW*0.98, 8, tr("INFORMACIÓN DE PAGO"), "", 0, "C", false, 0, "")
	pdf.Ln(8)

	pdf.SetFontStyle("")
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Valor inscripción:"), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")

	pdf.CellFormat(pps.wW*0.32, 5, tr(formatoDineroV2(0, "$", ",", data["Valor"].(string))), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr(textFecha), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["fecha1"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Código comprobante:"), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["Comprobante"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Estado del recibo: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["Estado"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)

	return pdf
}

func documentacionSuministradaV2(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	ystart := pdf.GetY()
	pdf.SetDrawColor(100, 100, 100)
	pdf.Cell(pps.wW*0.01, 8, "")
	pdf.SetFillColor(0, 162, 255)
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+1, pps.wW*.98, 6, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pps.wW*0.98, 8, tr("DOCUMENTACIÓN SUMINISTRADA"), "", 0, "C", false, 0, "")
	pdf.Ln(8)

	pdf.SetFontStyle("")
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY(), pps.wW*.38, 5, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.38, 5, "Componente", "", 0, "C", false, 0, "")
	pdf.Cell(pps.wW*0.04, 5, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY(), pps.wW*.54, 5, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.54, 5, "Documentos suministrados", "", 0, "C", false, 0, "")
	pdf.Ln(6)

	data = data["DOCUMENTACION"].(map[string]interface{})

	pdf = docsCarpetaV2(pdf, data, "Información Básica", true, pps)
	pdf = docsCarpetaV2(pdf, data, "Formación Académica", false, pps)
	pdf = docsCarpetaV2(pdf, data, "Experiencia Laboral", false, pps)
	pdf = docsCarpetaV2(pdf, data, "Producción Académica", true, pps)
	pdf = docsCarpetaV2(pdf, data, "Documentos Solicitados", false, pps)
	pdf = docsCarpetaV2(pdf, data, "Descuentos de Matrícula", false, pps)
	pdf = docsCarpetaV2(pdf, data, "Propuesta de Trabajo de Grado", false, pps)

	pdf.RoundedRect(pps.mL, ystart-1, pps.wW, 2+pdf.GetY()-ystart, 2.5, "1234", "")

	return pdf
}

func docsCarpetaV2(pdf *gofpdf.Fpdf, data map[string]interface{}, tagSuite string, subCarpeta bool, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	var ystart float64 = pdf.GetY()

	if thisTag, exist := data[tagSuite].(map[string]interface{}); exist {
		if subCarpeta {
			for sub, doc := range thisTag {
				total := len(doc.([]interface{}))
				pdf.SetX(pps.mL + pps.wW*.44)
				fontStyle(pdf, "B", 8, 100)
				pdf.MultiCell(pps.wW*.54, 5, tr(sub), "", "L", false)
				concatNames := fmt.Sprintf("(total docs: %d)  ", total)
				for _, docName := range doc.([]interface{}) {
					concatNames += (docName.(string) + ";  ")
				}
				concatNames = strings.Trim(concatNames, "; ")

				pdf.SetX(pps.mL + pps.wW*.44)
				fontStyle(pdf, "", 7, 0)
				pdf.MultiCell(pps.wW*.54, 4, tr(concatNames), "", "TL", false)
			}
			yend := pdf.GetY()
			pdf.RoundedRect(pps.mL+pps.wW*.44, ystart, pps.wW*.54, yend-ystart, 1, "1234", "")
			fontStyle(pdf, "B", 8, 0)
			pdf.SetXY(pps.mL+pps.wW*.02, ystart)
			pdf.CellFormat(pps.wW*0.38, yend-ystart, tr(tagSuite), "0", 0, "C", false, 0, "")
			pdf.RoundedRect(pps.mL+pps.wW*.02, ystart, pps.wW*.38, yend-ystart, 1, "1234", "")
			pdf.SetXY(pps.mL, yend+1)
		} else {
			for _, doc := range thisTag {
				total := len(doc.([]interface{}))
				fontStyle(pdf, "", 8, 0)
				concatNames := fmt.Sprintf("(total docs: %d)  ", total)
				for _, docName := range doc.([]interface{}) {
					concatNames += (docName.(string) + ";  ")
				}
				concatNames = strings.Trim(concatNames, "; ")
				pdf.SetX(pps.mL + pps.wW*.44)
				pdf.MultiCell(pps.wW*.54, 5, tr(concatNames), "", "L", false)
				yend := pdf.GetY()
				if (yend - ystart) <= 5 {
					yend += 2
				}
				pdf.RoundedRect(pps.mL+pps.wW*.44, ystart, pps.wW*.54, yend-ystart, 1, "1234", "")
				fontStyle(pdf, "B", 8, 0)
				pdf.SetXY(pps.mL+pps.wW*.02, ystart)
				pdf.CellFormat(pps.wW*0.38, yend-ystart, tr(tagSuite), "0", 0, "C", false, 0, "")
				pdf.RoundedRect(pps.mL+pps.wW*.02, ystart, pps.wW*.38, yend-ystart, 1, "1234", "")
				pdf.SetXY(pps.mL, yend+1)
				break
			}
		}
	}
	return pdf
}
func headerComprobanteV2(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) func() {
	return func() {
		pdf.SetHomeXY()
		tr := pdf.UnicodeTranslatorFromDescriptor("")

		path := beego.AppConfig.String("StaticPath")
		pdf = image(pdf, path+"/img/UDEscudo2.png", pps.mL, pps.mT, 0, 17.5)

		pdf.SetXY(pps.mL, pdf.GetY())
		fontStyle(pdf, "B", 10, 0)
		pdf.Cell(13, 10, "")
		pdf.Cell(140, 10, "UNIVERSIDAD DISTRITAL")
		pdf.Ln(4)

		pdf.Cell(13, 10, "")
		pdf.Cell(60, 10, tr("Francisco José de Caldas"))
		pdf.Cell(80, 10, tr("COMPROBANTE INSCRIPCIÓN"))
		pdf.Ln(4)

		fontStyle(pdf, "", 8, 0)
		pdf.Cell(13, 10, "")
		pdf.Cell(50, 10, "NIT 899.999.230-7")
		pdf.Ln(10)

		idPrograma := data["Dependencia"].(map[string]interface{})["Id"].(float64)
		idInscrip := data["INSCRIPCION"].(map[string]interface{})["id"].(float64)
		docAspirante := data["Documento"].(string)
		fechaInscrip := data["INSCRIPCION"].(map[string]interface{})["fechaInsripcion"].(string)
		fechaInscrip = strings.Split(fechaInscrip, ",")[0]

		codigo := fmt.Sprintf("%.f-%.f-%s-%s", idPrograma, idInscrip, docAspirante, fechaInscrip)
		bcode := barcode.RegisterCode128(pdf, codigo)
		barcode.Barcode(pdf, bcode, pps.mL+pps.wW-60, pps.mT+2.5, 58.5, 12, false)
	}
}

func footerComprobanteV2(pdf *gofpdf.Fpdf, pps styling) func() {
	return func() {
		pdf.SetXY(pps.mL, pps.mT+pps.hW-pps.hF)
		path := beego.AppConfig.String("StaticPath")
		pdf = image(pdf, path+"/img/sga_logo_name.png", pps.mL+pps.wW*0.5-17.66, pps.mT+pps.hW-pps.hF, 35.33, pps.hF)
	}
}
