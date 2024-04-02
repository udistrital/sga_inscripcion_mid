package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
            Method: "GetDescuentoAcademicoByDependenciaID",
            Router: "/:dependencia_id",
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
            Method: "GetDescuentoByPersonaPeriodoDependencia",
            Router: "/detalle",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoAcademicoByPersona",
            Router: "/persona/:persona_id",
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
            Method: "GetInformacionEmpresa",
            Router: "/informacion-empresa/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "GetExperienciaLaboralByTercero",
            Router: "/tercero/",
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
            Router: "/informacion-complementaria",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetInfoUniversidad",
            Router: "/informacion-universidad/nit/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetInfoUniversidadByNombre",
            Router: "/informacion-universidad/nombre/:nombre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "PostTercero",
            Router: "/tercero",
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
            Method: "PostGenerarComprobanteInscripcion",
            Router: "/comprobante-inscripcion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:GenerarReciboController"],
        beego.ControllerComments{
            Method: "PostGenerarRecibo",
            Router: "/estudiantes",
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
            Method: "GetEstadoInscripcion",
            Router: "/estado-recibos",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "ActualizarInfoContacto",
            Router: "/informacion-complementaria/tercero",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoComplementariaTercero",
            Router: "/informacion-complementaria/tercero",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "GetInfoComplementariaTercero",
            Router: "/informacion-complementaria/tercero/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoComplementariaUniversidad",
            Router: "/informacion-complementaria/universidad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInformacionFamiliar",
            Router: "/informacion-familiar",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostGenerarInscripcion",
            Router: "/nueva",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostPreinscripcion",
            Router: "/preinscripcion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "ConsultarProyectosEventos",
            Router: "/proyectos/eventos/:evento_padre_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoIcfesColegio",
            Router: "/pruebas-de-estado/informacion/saber-once",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoIcfesColegioNuevo",
            Router: "/pruebas-de-estado/informacion/saber-once-nuevo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostReintegro",
            Router: "/reintegro",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostTransferencia",
            Router: "/transferencia",
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
            Method: "PutProduccionAcademica",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetOneProduccionAcademica",
            Router: "/:id",
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
            Method: "GetIdProduccionAcademica",
            Router: "/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PutEstadoAutorProduccionAcademica",
            Router: "/autor/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetProduccionAcademica",
            Router: "/tercero/:tercero",
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
            Router: "/alerta-docente",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:SolicitudProduccionController"],
        beego.ControllerComments{
            Method: "PostSolicitudEvaluacionCoincidencia",
            Router: "/coincidencias",
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
            Router: "/actualizar-estado/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetConsultarParametros",
            Router: "/consultar-parametros",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetConsultarPeriodo",
            Router: "/consultar-periodo/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_inscripcion_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetEstadoInscripcion",
            Router: "/estado-recibos/:persona_id",
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
            Router: "/respuesta-solicitud/:id",
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
