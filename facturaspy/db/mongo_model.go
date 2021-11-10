package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

//Keep the sequence of the struct fileds as it may affect when serializing to files

//LIE type: libro ingresos/egresos
type LIE struct {
	Informante     Informante         `json:"informante" bson:"informante" xml:"informante"`
	Cantidades     []Cantidad         `json:"cantidades" bson:"cantidades" xml:"cantidades"`
	Identificacion Identificacion     `json:"identificacion" bson:"identificacion" xml:"identificacion"`
	Totales        Totales            `json:"totales" bson:"totales" xml:"totales"`
	Mongo_ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

//Informante type
// the person who is loading its thata
type Informante struct {
	Ruc               string       `json:"ruc" bson:"ruc" xml:"ruc"`
	Dv                string       `json:"dv" bson:"dv" xml:"dv"`
	Nombre            string       `json:"nombre" bson:"nombre" xml:"nombre"`
	TipoContribuyente string       `json:"tipoContribuyente" bson:"tipoContribuyente" xml:"tipoContribuyente"`
	TipoSociedad      string       `json:"tipoSociedad" bson:"tipoSociedad" xml:"tipoSociedad"`
	NombreFantasia    string       `json:"nombreFantasia" bson:"nombreFantasia" xml:"nombreFantasia"`
	Clasificacion     string       `json:"clasificacion" bson:"clasificacion" xml:"clasificacion"`
	Obligaciones      []Obligacion `json:"obligaciones" bson:"obligaciones" xml:"obligaciones"`
	LieID             uint
	LIEDetallesId     uint
}

//Obligacion type
type Obligacion struct {
	Impuesto     int    `json:"impuesto" bson:"impuesto" xml:"impuesto"`
	Nombre       string `json:"nombre" bson:"nombre" xml:"nombre"`
	FechaDesde   string `json:"fechaDesde" bson:"fechaDesde" xml:"fechaDesde"`
	InformanteId uint
}

//Identificacion type
type Identificacion struct {
	Periodo          string `json:"periodo" bson:"periodo" xml:"periodo"`
	TipoMovimiento   string `json:"tipoMovimiento" bson:"tipoMovimiento" xml:"tipoMovimiento"`
	TipoPresentacion string `json:"tipoPresentacion" bson:"tipoPresentacion" xml:"tipoPresentacion"`
	Version          string `json:"version" bson:"version" xml:"version"`
	LieID            uint
	LIEDetallesId    uint
}

//Cantidad Type
type Cantidad struct {
	Ingresos int `json:"ingresos" bson:"ingresos" xml:"ingresos"`
	Egresos  int `json:"egresos" bson:"egresos" xml:"egresos"`
	LieID    uint
}

//Totales type
type Totales struct {
	gorm.Model

	Ingresos      []Ingreso     `json:"ingresos" bson:"ingresos" xml:"ingresos"`
	Egresos       []Egreso      `json:"egresos" bson:"egresos" xml:"egresos"`
	ArbolIngresos ArbolIngresos `json:"arbolIngresos" bson:"arbolIngresos" xml:"arbolIngresos"`
	ArbolEgresos  ArbolEgresos  `json:"arbolEgresos" bson:"arbolEgresos" xml:"arbolEgresos"`
	LieID         uint
}

//Ingreso type
type Ingreso struct {
	Ruc            string `json:"ruc" bson:"ruc" xml:"ruc"`
	Periodo        int    `json:"periodo" bson:"periodo" xml:"periodo"`
	TipoIngreso    string `json:"tipoIngreso" bson:"tipoIngreso" xml:"tipoIngreso"`
	ValorGravado   int64  `json:"valorGravado" bson:"valorGravado" xml:"valorGravado"`
	ValorNoGravado int64  `json:"valorNoGravado" bson:"valorNoGravado" xml:"valorNoGravado"`
	CantidadId     uint
	TotalesId      uint
}

//Egreso type
type Egreso struct {
	Ruc           string `json:"ruc" bson:"ruc" xml:"ruc"`
	Periodo       int    `json:"periodo" bson:"periodo" xml:"periodo"`
	TipoEgreso    string `json:"tipoEgreso" bson:"tipoEgreso" xml:"tipoEgreso"`
	Clasificacion string `json:"clasificacion" bson:"clasificacion" xml:"clasificacion"`
	Valor         int    `json:"valor" bjson:"valor" xml:"valor"`
	CantidadId    uint
	TotalesId     uint
}

//ArbolIngresos type
type ArbolIngresos struct {
	SubtotalGravado   int   `json:"subtotalGravado" bson:"subtotalGravado" xml:"subtotalGravado"`
	SubtotalNoGravado int   `json:"subtotalNoGravado" bson:"subtotalNoGravado" xml:"subtotalNoGravado"`
	HPRSP             HPRSP `json:"HPRSP" bson:"HPRSP" xml:"HPRSP"`
	TotalesId         uint
}

//HPRSP type
type HPRSP struct {
	Gravado         int `json:"gravado" bson:"gravado" xml:"gravado"`
	NoGravado       int `json:"noGravadoint" bson:"noGravadoint" xml:"noGravadoint"`
	ArbolIngresosId uint
}

//ArbolEgresos type
type ArbolEgresos struct {
	Gasto     Gasto `json:"gasto" bson:"gasto" xml:"gasto"`
	TotalesId uint
}

//Gasto type
type Gasto struct {
	Total          int `json:"total" bson:"total" xml:"total"`
	GACT           int `json:"GACT" bson:"GACT" xml:"GACT"`
	GPERS          int `json:"GPERS" bson:"GPERS" xml:"GPERS"`
	ArbolEgresosId uint
}

// detalles

//LIEDetalles type
type LIEDetalles struct {
	Informante     Informante       `json:"informante" bson:"informante" xml:"informante"`
	Identificacion Identificacion   `json:"identificacion" bson:"identificacion" xml:"identificacion"`
	Ingresos       []IngresoDetalle `json:"ingresos" bson:"ingresos" bson:"ingresos"`
	Egresos        []EgresoDetalle  `json:"egresos" bson:"egresos" xml:"egresos"`
	//Familiares []Familiar `json:"familiares"`
	_ID      string             `bson:"_id,omitempty"`
	LieID    string             `json:"-" bson:"lie_id,omitempty"`
	Mongo_ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

//IngresoDetalle type
//the order of the filds must not change as it is used
//to generate the xls file headers
type IngresoDetalle struct {
	Tipo                            string `json:"tipo,omitempty" bson:"tipo" xml:"tipo"`
	TipoTexto                       string `json:"tipoTexto,omitempty" bson:"tipoTexto" xml:"tipoTexto"`
	TipoIngreso                     string `json:"tipoIngreso,omitempty" bson:"tipoIngreso" xml:"tipoIngreso"`
	TipoIngresoTexto                string `json:"tipoIngresoTexto,omitempty" bson:"tipoIngresoTexto" xml:"tipoIngresoTexto"`
	Fecha                           string `json:"fecha,omitempty" bson:"fecha" xml:"fecha"`
	Mes                             string `json:"mes,omitempty" bson:"mes" xml:"mes"` //?
	RelacionadoTipoIdentificacion   string `json:"relacionadoTipoIdentificacion,omitempty" bson:"relacionadoTipoIdentificacion" xml:"relacionadoTipoIdentificacion"`
	RelacionadoNumeroIdentificacion string `json:"relacionadoNumeroIdentificacion,omitempty" bson:"relacionadoNumeroIdentificacion" xml:"relacionadoNumeroIdentificacion"`
	RelacionadoNombres              string `json:"relacionadoNombres,omitempty" bson:"relacionadoNombres" xml:"relacionadoNombres"`
	TimbradoNumero                  string `json:"timbradoNumero,omitempty" bson:"timbradoNumero" xml:"timbradoNumero"`
	TimbradoDocumento               string `json:"timbradoDocumento,omitempty" bson:"timbradoDocumento" xml:"timbradoDocumento"`
	TimbradoCondicion               string `json:"timbradoCondicion,omitempty" bson:"timbradoCondicion" xml:"timbradoCondicion"`
	IngresoMontoGravado             int64  `json:"ingresoMontoGravado" bson:"ingresoMontoGravado" xml:"ingresoMontoGravado"`
	IngresoMontoNoGravado           int64  `json:"ingresoMontoNoGravado" bson:"ingresoMontoNoGravado" xml:"ingresoMontoNoGravado"`
	IngresoMontoTotal               int64  `json:"ingresoMontoTotal" bson:"ingresoMontoTotal" xml:"ingresoMontoTotal"`
	// TODO
	// TMPNumerodeCuenta                           string `json:"NumerodeCuenta"`
	// TMPRazonSocialdelBancoFinancieraCooperativa string `json:"RazonSocialdelBancoFinancieraCooperativa"`
	// TMPOtroTipoDeDocumento                      string `json:"OtroTipoDeDocumento"`
	Periodo       string `json:"periodo" bson:"periodo" xml:"periodo"`
	ID            int    `json:"id" bson:"id" xml:"id"`
	Ruc           string `json:"ruc" bson:"ruc" xml:"ruc"`
	_ID           string `bson:"_id,omitempty"`
	LieId         uint
	LIEDetallesId uint
}

//EgresoDetalle type
//TODO: USE tags and positions to generate xls
type EgresoDetalle struct {
	Tipo                            string `json:"tipo,omitempty" bson:"tipo" xml:"tipo" xls:"Tipo de Documento"`
	TipoTexto                       string `json:"tipoTexto,omitempty" bson:"tipoTexto" xml:"tipoTexto" xls:"Tipo de Documento (Texto)"`
	Fecha                           string `json:"fecha,omitempty" bson:"fecha" xml:"fecha" xls:"Fecha"`
	Mes                             string `json:"mes,omitempty" bson:"mes" xml:"mes" xls:"Mes"`
	Ruc                             string `json:"ruc,omitempty" bson:"ruc" xml:"ruc" xls:"Tipo de Identificación"`
	RelacionadoTipoIdentificacion   string `json:"relacionadoTipoIdentificacion,omitempty" bson:"relacionadoTipoIdentificacion" xml:"relacionadoTipoIdentificacion" xls:"Número de Identificación"`
	RelacionadoNumeroIdentificacion string `json:"relacionadoNumeroIdentificacion,omitempty" bson:"relacionadoNumeroIdentificacion" xml:"relacionadoNumeroIdentificacion" xls:"Número de Timbrado"`
	RelacionadoNombres              string `json:"relacionadoNombres,omitempty" bson:"relacionadoNombres" xml:"relacionadoNombres" xls:"Nombres y Apellidos o Razón Social"`
	TimbradoNumero                  string `json:"timbradoNumero,omitempty" bson:"timbradoNumero" xml:"timbradoNumero" xls:"Número de Timbrado"`
	TimbradoDocumento               string `json:"timbradoDocumento,omitempty" bson:"timbradoDocumento" xml:"timbradoDocumento" xls:"Número de Documento"`
	TimbradoCondicion               string `json:"timbradoCondicion,omitempty" bson:"timbradoCondicion" xml:"timbradoCondicion" xls:"Condición de la Venta"`
	EgresoMontoTotal                int64  `json:"egresoMontoTotal,omitempty" xml:"egresoMontoTotal" xls:"Monto Total"`
	NroCuenta                       string `json:"nroCuenta,omitempty" bson:"nroCuenta" xml:"nroCuenta" xls:"Número de Cuenta"` //?
	RazonSocial                     string `xls:"Razón Social del Banco / Financiera / Cooperativa"`                            //?
	OtroTipoDeDocumento             string `xls:"Otro Tipo de Documento"`
	NumeroDocumento                 string `json:"numeroDocumento,omitempty" bson:"numeroDocumento" xml:"numeroDocumento" xls:"Número de Documento"` //?
	NroDespacho                     string `xls:"Número de Despacho"`
	PeriodoCuenta                   string `xls:"Período de la Cuenta"`
	IDEmpleador                     string `xls:"Identificador del Empleador"`
	TipoEgreso                      string `json:"tipoEgreso,omitempty" bson:"tipoEgreso" xml:"tipoEgreso" xls:"Tipo de Egreso"`
	TipoEgresoTexto                 string `json:"tipoEgresoTexto,omitempty" bson:"tipoEgresoTexto" xml:"tipoEgresoTexto" xls:"Tipo de Egreso (Texto)"`
	SubtipoEgreso                   string `json:"subtipoEgreso,omitempty" bson:"subtipoEgreso" xml:"subtipoEgreso" xls:"Clasificación de Egreso"`
	SubtipoEgresoTexto              string `json:"subtipoEgresoTexto,omitempty" bson:"subtipoEgresoTexto" xml:"subtipoEgresoTexto" xls:"Clasificación de Egreso (Texto)"`
	NroIDEmpleador                  string `xls:"Número de Identificación Del Empleador"`
	ID                              int    `json:"id,omitempty" bson:"id"`
	Periodo                         string `json:"periodo,omitempty" bson:"periodo" xml:"periodo"`
	_ID                             string `bson:"_id,omitempty"`
	LieId                           uint
	LIEDetallesId                   uint
}
