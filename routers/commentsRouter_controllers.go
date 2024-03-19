package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "PostDocs",
            Router: "/comentarios",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:CuposController"],
        beego.ControllerComments{
            Method: "GetAllDocs",
            Router: "/comentarios",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "PostDescuentoAcademico",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoAcademico",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "PutDescuentoAcademico",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoAcademicoByPersona",
            Router: "/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoAcademicoByDependenciaID",
            Router: "/descuentoAcademicoByID/:dependencia_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoByPersonaPeriodoDependencia",
            Router: "/descuentopersonaperiododependencia/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "PostExperienciaLaboral",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "PutExperienciaLaboral",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "GetExperienciaLaboral",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "DeleteExperienciaLaboral",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "GetExperienciaLaboralByTercero",
            Router: "/by_tercero/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "GetInformacionEmpresa",
            Router: "/informacion_empresa/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "PostFormacionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "PutFormacionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetFormacionAcademicaByTercero",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "DeleteFormacionAcademica",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetFormacionAcademica",
            Router: "/info_complementaria/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetInfoUniversidad",
            Router: "/info_universidad/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetInfoUniversidadByNombre",
            Router: "/info_universidad_nombre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "PostTercero",
            Router: "/post_tercero",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GeneradorCodigoBarrasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GeneradorCodigoBarrasController"],
        beego.ControllerComments{
            Method: "GenerarCodigoBarras",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"],
        beego.ControllerComments{
            Method: "PostGenerarRecibo",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"],
        beego.ControllerComments{
            Method: "PostGenerarComprobanteInscripcion",
            Router: "/comprobante_inscripcion/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"],
        beego.ControllerComments{
            Method: "PostGenerarEstudianteRecibo",
            Router: "/estudiantes/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "ConsultarProyectosEventos",
            Router: "/consultar_proyectos_eventos/:evento_padre_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "GetEstadoInscripcion",
            Router: "/estado_recibos/:persona_id/:id_periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostGenerarInscripcion",
            Router: "/generar_inscripcion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoComplementariaTercero",
            Router: "/info_complementaria_tercero",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "GetInfoComplementariaTercero",
            Router: "/info_complementaria_tercero/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoComplementariaUniversidad",
            Router: "/info_complementaria_universidad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "ActualizarInfoContacto",
            Router: "/info_contacto",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoIcfesColegio",
            Router: "/post_info_icfes_colegio",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoIcfesColegioNuevo",
            Router: "/post_info_icfes_colegio_nuevo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInformacionFamiliar",
            Router: "/post_informacion_familiar",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostPreinscripcion",
            Router: "/post_preinscripcion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostReintegro",
            Router: "/post_reintegro",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostTransferencia",
            Router: "/post_transferencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PostProduccionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetAllProduccionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "DeleteProduccionAcademica",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PutProduccionAcademica",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetProduccionAcademica",
            Router: "/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PutEstadoAutorProduccionAcademica",
            Router: "/estado_autor_produccion/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetOneProduccionAcademica",
            Router: "/get_one/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetIdProduccionAcademica",
            Router: "/pr_academica/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"],
        beego.ControllerComments{
            Method: "PutResultadoSolicitud",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"],
        beego.ControllerComments{
            Method: "PostAlertSolicitudProduccion",
            Router: "/:tercero/:tipo_produccion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"],
        beego.ControllerComments{
            Method: "PostSolicitudEvaluacionCoincidencia",
            Router: "/coincidencia/:id_solicitud/:id_coincidencia/:id_tercero",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PostSolicitud",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PutInfoSolicitud",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PutInscripcion",
            Router: "/actualizar_estado/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetConsultarParametros",
            Router: "/consultar_parametros/:id_calendario/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetConsultarPeriodo",
            Router: "/consultar_periodo/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetEstadoInscripcion",
            Router: "/estado_recibos/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetEstados",
            Router: "/estados",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetInscripcion",
            Router: "/inscripcion/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PutSolicitud",
            Router: "/respuesta_solicitud/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetSolicitudesInscripcion",
            Router: "/solicitudes/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
