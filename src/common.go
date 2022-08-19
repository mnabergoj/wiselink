package src

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"mime/quotedprintable"
	"net/smtp"
	"strconv"
	"time"
)

//Acceso a DataBase
const MaxConnections int = 10 //ex 20
const DBip string = "127.0.0.1"
const DBuser string = "root"
const DBpass string = "ManVp2020DB"

//const DBpass string = "root"

//Datos acceso Auth0 SandBox(mnabergoj@gmail.com)
/*const Auth_Domain string = "https://dev-mk8y7bal.auth0.com"
const Auth_ClientId string = "lbwS1TkdM68Lii1z99vGHcFhOWW7aP9e"
const Auth_Secret string = "OS5lNyOJp5o8thaEuDAVi-7bhXvTGJ5AJ157F-sRiWHwt2CEth2J4VhHKjUUwzn-"
const Auth_Connection string = "Username-Password-Authentication" */

//Datos acceso Auth0 Produccion(smartalarmplat@gmail.com)
const Auth_Domain string = "https://dev-dwteghkl.us.auth0.com"
const Auth_ClientId string = "QUUi313KvCbcjGHnfdkeOAIZzM9fZLE6"
const Auth_Secret string = "dZYNgk7EyrU_QYUXeoWNi_URt3irQ-124nH_nSy2k5qqrCY8s2hnSuiHGBA2CRHU"
const Auth_Connection string = "Username-Password-Authentication"

//Datos de Conexion a Bridge API. PORTS SandBox:8082  Produccion:8086
const BridgeDomain string = "http://127.0.0.1:8086" //In Server solo soporta HTTP
//const BridgeDomain string = "http://127.0.0.1:8082" //In Server solo soporta HTTP
//const BridgeDomain string = "http://dev.local:8082"

//Conexiones a DBs
/*var Connection string = fmt.Sprintf("%s:%s@tcp(%s:3306)/smartalarm_crmsb", DBuser, DBpass, DBip)
var ConnectFinance string = fmt.Sprintf("%s:%s@tcp(%s:3306)/smartalarm_financesb", DBuser, DBpass, DBip) */

var Connection string = fmt.Sprintf("%s:%s@tcp(%s:3306)/smartalarm_crm", DBuser, DBpass, DBip)
var ConnectFinance string = fmt.Sprintf("%s:%s@tcp(%s:3306)/smartalarm_finances", DBuser, DBpass, DBip)

//IDP_DISTRI MAESTRO
const IDP_DISTRI_MASTER string = "6YJowNPMSoBM9zoIwWYtKgaRJETJJpTJ"

//Datos Cuenta de Email (SMTP)
const from_email string = "smartalarmplat@gmail.com"

//const password string = "ManVp2020Ad"
const password string = "znijnjkboyqsucmf" //"ManVp2020Ad"
const host string = "smtp.gmail.com:587"

//Matriz General Modulos -----------------------------------------------------------------------------------
type Modulos struct {
	Distribuidores Permisos `json:"distribuidores"`
	Usuarios       Permisos `json:"usuarios"`
	Clientes       Permisos `json:"clientes"`
	Servicios      Permisos `json:"servicios"`
	Equipos        Permisos `json:"equipos"`
}
type Permisos struct {
	View   bool `json:"view"`
	New    bool `json:"new"`
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
}

//Request Error Response -----------------------------------------------------------------------------------
type Request_Err struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Err_code int    `json:"err_code"`
}

//Conexion a DB --------------------------------------------------------------------------------------------
var db *sql.DB

func InitDB() {
	var dberr error
	db, dberr = sql.Open("mysql", Connection)
	if dberr != nil {
		log.Panic(dberr)
	}
	db.SetMaxOpenConns(MaxConnections)
	db.SetMaxIdleConns(MaxConnections)
	db.SetConnMaxLifetime(5 * time.Minute)

	if dberr = db.Ping(); dberr != nil {
		log.Panic(dberr)
	}
	fmt.Println("DB Iniciada con exito!")
} //--------------------------------------------------------------------------------------------------------
//Verifica la Conexion a DB---------------------------------------------------------------------------------
func PingDB() bool {
	err := db.Ping()
	if err != nil {
		db.Close()
		InitDB()
		return false
	}
	return true
} //------------------------------------------------------------------------------------------------------

//Conexion a DB2 -------------------------------------------------------------------------------------------
var db2 *sql.DB

func InitDB2() {
	var dberr error
	db2, dberr = sql.Open("mysql", ConnectFinance)
	if dberr != nil {
		log.Panic(dberr)
	}
	db2.SetMaxOpenConns(MaxConnections)
	db2.SetMaxIdleConns(MaxConnections)
	db2.SetConnMaxLifetime(5 * time.Minute)

	if dberr = db2.Ping(); dberr != nil {
		log.Panic(dberr)
	}
	fmt.Println("DB2 Iniciada con exito!")
} //--------------------------------------------------------------------------------------------------------
//Verifica la Conexion a DB---------------------------------------------------------------------------------
func PingDB2() bool {
	err := db2.Ping()
	if err != nil {
		db2.Close()
		InitDB2()
		return false
	}
	return true
} //--------------------------------------------------------------------------------------------------------
//Retorna Timestamp Actual----------------------------------------------------------------------------------
func Timestamp_now() string {
	loc, _ := time.LoadLocation("America/Buenos_Aires") //loc, _ := time.LoadLocation("Local")
	t := time.Now().In(loc)
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return timestamp
} //---------------------------------------------------------------------------------------------------------
//Retorna Time Actual(time.Time)-----------------------------------------------------------------------------
func Time_now() time.Time {
	loc, _ := time.LoadLocation("America/Buenos_Aires")
	t := time.Now().In(loc)

	return t
} //---------------------------------------------------------------------------------------------------------

//Creacion de ID HEX 32Bits ---------------------------------------------------------------------------------
func GetToken(length int) string {
	token := ""
	codeAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	codeAlphabet += "abcdefghijklmnopqrstuvwxyz"
	codeAlphabet += "0123456789"

	for i := 0; i < length; i++ {
		token += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
	}
	return token
}
func cryptoRandSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
} //-----------------------------------------------------------------------------------------------------

//Ortener un nuevo Password con Matysculas, Minusculas y Especials --------------------------------------
func GetNewPASSWORD(length int) string {
	token := ""
	codeAlphabet1 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	codeAlphabet2 := "abcdefghijklmnopqrstuvwxyz"
	codeAlphabet3 := "0123456789"
	codeAlphabet4 := "!@#$%^&*"
	codeAlphabet := codeAlphabet1 + codeAlphabet2 + codeAlphabet3 + codeAlphabet4
	//Genera el Password
	for i := 0; i < length; i++ {
		token += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
	}
	//Verifica Contenido Seguro
	grupo1 := false
	grupo2 := false
	grupo3 := false
	grupo4 := false
	count := 0
	//Intentando Generar un Password Valido
	for {
		//Verifica que existan algun caracter especial
		for i := 0; i < length; i++ {
			for r := 0; r < len(codeAlphabet1); r++ {
				if token[i:(i+1)] == codeAlphabet1[r:(r+1)] {
					grupo1 = true
					break
				}
			}
			for r := 0; r < len(codeAlphabet2); r++ {
				if token[i:(i+1)] == codeAlphabet2[r:(r+1)] {
					grupo2 = true
					break
				}
			}
			for r := 0; r < len(codeAlphabet3); r++ {
				if token[i:(i+1)] == codeAlphabet3[r:(r+1)] {
					grupo3 = true
					break
				}
			}
			for r := 0; r < len(codeAlphabet4); r++ {
				if token[i:(i+1)] == codeAlphabet4[r:(r+1)] {
					grupo4 = true
					break
				}
			}
			if grupo1 && grupo2 && grupo3 && grupo4 {
				return token
			}
		}
		//Reintenta generar el Password
		token = ""
		for i := 0; i < length; i++ {
			token += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
		}
		//Controla Limite DO-WHILE
		count++
		if count > 65530 {
			break
		}
	}

	return token
} //-----------------------------------------------------------------------------------------------------

//Validamos Distribuidor --------------------------------------------------------------------------------
func ValidateDistri(idp_distri string) bool {
	//Verificamos idp_distri MAESTRO (El siguiente ID es elegido por convencion)
	if idp_distri != IDP_DISTRI_MASTER {

		//Verificamos idp_distri General
		results, err := db.Query("SELECT id FROM distribuidores WHERE idp_distri='" + idp_distri + "'")
		defer results.Close()
		if err != nil {
			//panic(err.Error())
			return false
		}

		if !results.Next() {
			return false
		}
	}
	return true
} //-----------------------------------------------------------------------------------------------------

//Validamos Plan ----------------------------------------------------------------------------------------
func ValidatePlan(plan string, tipo int) bool {

	//Verificamos Plan
	results, err := db2.Query("SELECT id FROM planes WHERE idp_plan='" + plan + "' and tipo=" + strconv.Itoa(tipo))
	defer results.Close()
	if err != nil {
		return false
	}

	if !results.Next() {
		return false
	}
	return true
} //-----------------------------------------------------------------------------------------------------

//Validamos Numero de Serie -----------------------------------------------------------------------------
func ValidateSerial(serial string) bool {

	//Verificamos Plan
	results, err := db.Query("SELECT id FROM equipos WHERE idp_serie='" + serial + "' AND syncro != 'S'")
	defer results.Close()
	if err != nil {
		return false
	}

	if !results.Next() {
		return false
	}
	return true
} //-----------------------------------------------------------------------------------------------------

//Validamos Distribuidor-Cliente ------------------------------------------------------------------------
func ValidateDistriClient(idp_distri string, idp_cliente string) bool {
	//Verificamos idp_distri MAESTRO (El siguiente ID es elegido por convencion)

	//Conectamos a DB CRM
	/*db, err := sql.Open("mysql", Connection)
	if err != nil {
		//panic(err.Error())
		return false
	}
	defer db.Close()
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 10. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(5) */

	//Verificamos idp_distri/idp_cliente
	Sql := ""
	if idp_distri != IDP_DISTRI_MASTER {
		Sql = "SELECT id FROM clientes WHERE idp_distri='" + idp_distri + "' AND idp_cliente='" + idp_cliente + "'"
	} else {
		Sql = "SELECT id FROM clientes WHERE idp_cliente='" + idp_cliente + "'"
	}
	results, err := db.Query(Sql)
	defer results.Close()

	if err != nil {
		//panic(err.Error())
		return false
	}

	if !results.Next() {
		return false
	}

	return true
} //-----------------------------------------------------------------------------------------------------

/*Validate Edit Status ----------------------------------------------------------------------------------
0 - Null
1 - Eliminado
2 - Pre inscripcion
3 - Inscripto
4 - Pauseado
5 - Mora
6 - Bloqueado
10 - Activo
11 - Activo Usuario Free
*/
func ValidateEditStatus(statusIn int, status int) bool {
	//Evalua las posibilidades segun status actual
	if statusIn >= 0 && statusIn <= 1 {
		return false
	}
	//Evalua Posibilidades de estado inicial
	if status == 3 && statusIn == 2 {
		return false
	}
	if (status == 10 || status == 11) && statusIn < 4 {
		return false
	}

	return true
} //-----------------------------------------------------------------------------------------------------

//Change Client Status ----------------------------------------------------------------------------------
func ChangeClientStatus(status int, idp string) bool {

	//Editamos el Status
	statusx := strconv.Itoa(status)
	update, err := db.Query("UPDATE clientes SET status=" + statusx + " WHERE idp_cliente='" + idp + "'")
	defer update.Close()
	if err != nil {
		return false
	}
	update.Close()

	return true
} //---------------------------------------------------------------------------------------------------

//Validamos My Idp User ---------------------------------------------------------------------------------
func ValidateMyIdp(myIdp string) (bool, Modulos, string) {
	var out Modulos

	//Verificamos idp_user
	results, err := db.Query("SELECT id,idp_distri,str_modulos FROM usuarios WHERE idp_user='" + myIdp + "'")
	defer results.Close()
	if err != nil {
		//panic(err.Error())
		return false, out, ""
	}

	if !results.Next() {
		return false, out, ""
	}
	//Carga los Modulos
	var item Usuario
	err = results.Scan(&item.Id, &item.Idp_distri, &item.StrModulos)
	if err != nil {
		return false, out, ""
	}

	json.Unmarshal([]byte(item.StrModulos), &out)
	return true, out, item.Idp_distri
} //-----------------------------------------------------------------------------------------------------

//Verificamos que el Modulo a conceder no sea Superior al propio ----------------------------------------
func CompareModulos(my Modulos, then Modulos) bool {
	//Compara uno a uno los Permisos
	if then.Distribuidores.Delete && !my.Distribuidores.Delete {
		return false
	}
	if then.Distribuidores.Edit && !my.Distribuidores.Edit {
		return false
	}
	if then.Distribuidores.View && !my.Distribuidores.View {
		return false
	}
	if then.Distribuidores.New && !my.Distribuidores.New {
		return false
	}

	if then.Usuarios.Delete && !my.Usuarios.Delete {
		return false
	}
	if then.Usuarios.Edit && !my.Usuarios.Edit {
		return false
	}
	if then.Usuarios.View && !my.Usuarios.View {
		return false
	}
	if then.Usuarios.New && !my.Usuarios.New {
		return false
	}
	if then.Clientes.Delete && !my.Clientes.Delete {
		return false
	}
	if then.Clientes.Edit && !my.Clientes.Edit {
		return false
	}
	if then.Clientes.View && !my.Clientes.View {
		return false
	}
	if then.Clientes.New && !my.Clientes.New {
		return false
	}
	if then.Servicios.Delete && !my.Servicios.Delete {
		return false
	}
	if then.Servicios.Edit && !my.Servicios.Edit {
		return false
	}
	if then.Servicios.View && !my.Servicios.View {
		return false
	}
	if then.Servicios.New && !my.Servicios.New {
		return false
	}
	if then.Equipos.Delete && !my.Equipos.Delete {
		return false
	}
	if then.Equipos.Edit && !my.Equipos.Edit {
		return false
	}
	if then.Equipos.View && !my.Equipos.View {
		return false
	}
	if then.Equipos.New && !my.Equipos.New {
		return false
	}
	return true
} //-----------------------------------------------------------------------------------------------------

//Verificamos que el Modulo sea NULL --------------------------------------------------------------------
func DetecModulosNull(my Modulos) bool {
	//Compara uno a uno los Permisos
	if my.Distribuidores.Delete {
		return false
	}
	if my.Distribuidores.Edit {
		return false
	}
	if my.Distribuidores.View {
		return false
	}
	if my.Distribuidores.New {
		return false
	}

	if my.Usuarios.Delete {
		return false
	}
	if my.Usuarios.Edit {
		return false
	}
	if my.Usuarios.View {
		return false
	}
	if my.Usuarios.New {
		return false
	}

	if my.Clientes.Delete {
		return false
	}
	if my.Clientes.Edit {
		return false
	}
	if my.Clientes.View {
		return false
	}
	if my.Clientes.New {
		return false
	}

	if my.Servicios.Delete {
		return false
	}
	if my.Servicios.Edit {
		return false
	}
	if my.Servicios.View {
		return false
	}
	if my.Servicios.New {
		return false
	}

	if my.Equipos.Delete {
		return false
	}
	if my.Equipos.Edit {
		return false
	}
	if my.Equipos.View {
		return false
	}
	if my.Equipos.New {
		return false
	}

	return true
} //-----------------------------------------------------------------------------------------------------

//Obtenemos un Modulo Resultante de los permisos del Distribuidos y el Usuario --------------------------
func MergeModulos(my Modulos, then Modulos) Modulos {
	var merge Modulos

	//Compara uno a uno los Permisos
	if then.Distribuidores.Delete && my.Distribuidores.Delete {
		merge.Distribuidores.Delete = true
	}
	if then.Distribuidores.Edit && my.Distribuidores.Edit {
		merge.Distribuidores.Edit = true
	}
	if then.Distribuidores.View && my.Distribuidores.View {
		merge.Distribuidores.View = true
	}
	if then.Distribuidores.New && my.Distribuidores.New {
		merge.Distribuidores.New = true
	}

	if then.Usuarios.Delete && my.Usuarios.Delete {
		merge.Usuarios.Delete = true
	}
	if then.Usuarios.Edit && my.Usuarios.Edit {
		merge.Usuarios.Edit = true
	}
	if then.Usuarios.View && my.Usuarios.View {
		merge.Usuarios.View = true
	}
	if then.Usuarios.New && my.Usuarios.New {
		merge.Usuarios.New = true
	}

	if then.Clientes.Delete && my.Clientes.Delete {
		merge.Clientes.Delete = true
	}
	if then.Clientes.Edit && my.Clientes.Edit {
		merge.Clientes.Edit = true
	}
	if then.Clientes.View && my.Clientes.View {
		merge.Clientes.View = true
	}
	if then.Clientes.New && my.Clientes.New {
		merge.Clientes.New = true
	}

	if then.Servicios.Delete && my.Servicios.Delete {
		merge.Servicios.Delete = true
	}
	if then.Servicios.Edit && my.Servicios.Edit {
		merge.Servicios.Edit = true
	}
	if then.Servicios.View && my.Servicios.View {
		merge.Servicios.View = true
	}
	if then.Servicios.New && my.Servicios.New {
		merge.Servicios.New = true
	}

	if then.Equipos.Delete && my.Equipos.Delete {
		merge.Equipos.Delete = true
	}
	if then.Equipos.Edit && my.Equipos.Edit {
		merge.Equipos.Edit = true
	}
	if then.Equipos.View && my.Equipos.View {
		merge.Equipos.View = true
	}
	if then.Equipos.New && my.Equipos.New {
		merge.Equipos.New = true
	}

	return merge
} //-----------------------------------------------------------------------------------------------------

//Obtenemos el ID del Usuario desde IDP -----------------------------------------------------------------
func GetUserId(Idp string) int32 {

	//Verificamos idp_user
	results, err := db.Query("SELECT id FROM usuarios WHERE idp_user='" + Idp + "'")
	defer results.Close()
	if err != nil {
		//panic(err.Error())
		return 0
	}

	if !results.Next() {
		return 0
	}
	//Carga los Modulos
	var item Usuario
	err = results.Scan(&item.Id)
	if err != nil {
		return 0
	}

	return item.Id
} //-----------------------------------------------------------------------------------------------------

//Obtenemos el nombre del Distribuidor desde IDP --------------------------------------------------------
func GetDistriName(Idp string) string {

	//Detectamos SuperUsuario
	if Idp == IDP_DISTRI_MASTER {
		return "SuperUsuario"
	}

	//Verificamos idp_user
	results, err := db.Query("SELECT nombre FROM distribuidores WHERE idp_distri='" + Idp + "'")
	defer results.Close()
	if err != nil {
		//panic(err.Error())
		return "error(2)"
	}

	if !results.Next() {
		return "error(3)"
	}
	//Carga los Modulos
	var item Distribuidor
	err = results.Scan(&item.Nombre)
	if err != nil {
		return "error(4)"
	}

	return item.Nombre
} //-----------------------------------------------------------------------------------------------------

//Obtenemos el Email del Distribuidor desde IDP ---------------------------------------------------------
func GetDistriEmail(Idp string) (email string, err error) {

	//Verificamos idp_user
	results, err := db.Query("SELECT mail FROM distribuidores WHERE idp_distri='" + Idp + "'")
	defer results.Close()
	if err != nil {
		return "", err
	}

	if !results.Next() {
		return "", err
	}
	//Carga los Modulos
	var item Distribuidor
	err = results.Scan(&item.Mail)
	if err != nil {
		return "", err
	}

	return item.Mail, err
} //-----------------------------------------------------------------------------------------------------

//Obtenemos datos varios del Distribuidor desde IDP para uso en LogIn -----------------------------------
func GetDistriVar(Idp string) (int, int, Modulos, string) {
	var out Modulos

	//Detectamos SuperUsuario
	if Idp == IDP_DISTRI_MASTER {
		return 1, 10, out, "SuperUsuario"
	}

	//Verificamos idp_user
	results, err := db.Query("SELECT status,tipo,str_modulos FROM distribuidores WHERE idp_distri='" + Idp + "'")
	defer results.Close()
	if err != nil {
		//panic(err.Error())
		return 0, 0, out, "error(2)"
	}

	if !results.Next() {
		return 0, 0, out, "error(3)"
	}
	//Carga los Modulos
	var item Distribuidor
	err = results.Scan(&item.Status, &item.Tipo, &item.StrModulos)
	if err != nil {
		return 0, 0, out, "error(4)"
	}
	json.Unmarshal([]byte(item.StrModulos), &out)

	return item.Status, item.Tipo, out, "ok"
} //-----------------------------------------------------------------------------------------------------

//Validamos Mail y DNI de Cliente como unicos -----------------------------------------------------------
func ValidateClient(email string, dni string) bool {

	//Verificamos idp_distri General
	results, err := db.Query("SELECT id FROM clientes WHERE email='" + email + "' OR dni='" + dni + "'")
	defer results.Close()
	if err != nil {
		//panic(err.Error())
		return false
	}

	if results.Next() {
		return false
	}
	return true
} //-----------------------------------------------------------------------------------------------------

//Validamos Mail de Cliente como unico ------------------------------------------------------------------
func ValidateClientX(email string) bool {

	//Verificamos idp_distri General
	results, err := db.Query("SELECT id FROM clientes WHERE email='" + email + "'")
	defer results.Close()
	if err != nil {
		//panic(err.Error())
		return false
	}

	if results.Next() {
		return false
	}
	return true
} //-----------------------------------------------------------------------------------------------------

//Obtenemos IDP del Cliente por Mail --------------------------------------------------------------------
func GetClientIDP(email string) (st bool, idp string, err error) {

	//Verificamos idp_distri General
	results, err := db.Query("SELECT idp_cliente FROM clientes WHERE email='" + email + "'")
	defer results.Close()
	if err != nil {
		return
	}
	if !results.Next() {
		return
	}
	//Carga los Modulos
	err = results.Scan(&idp)
	if err != nil {
		return
	}

	return true, idp, err
} //-----------------------------------------------------------------------------------------------------

//Obtenemos la Cantidad de servicios del Cliente --------------------------------------------------------
func GetServicesNumber(Idp string) (int, string) {

	var count int
	//Obtiene el numero de Registros
	err := db.QueryRow("SELECT COUNT(*) FROM servicios WHERE idp_cliente='" + Idp + "'").Scan(&count)
	switch {
	case err != nil:
		return 0, "error2"
	default:
		return count, "ok"
	}
} //-----------------------------------------------------------------------------------------------------

//Send Email---------------------------------------------------------------------------------------------
func SendEmail(to_email string, subject string, body string) (bool, string) {

	fmt.Println("Enviando Email!")
	auth := smtp.PlainAuth("", from_email, password, "smtp.gmail.com")

	header := make(map[string]string)
	header["From"] = "Administración SmartAlarm" //Queda fijo por el momento
	header["To"] = to_email
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", "text/html")
	header["Content-Disposition"] = "inline"
	header["Content-Transfer-Encoding"] = "quoted-printable"
	//Empaqueta header
	header_message := ""
	for key, value := range header {
		header_message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	//body := "<h1>This is your HTML Body</h1>"
	var body_message bytes.Buffer
	temp := quotedprintable.NewWriter(&body_message)
	temp.Write([]byte(body))
	temp.Close()

	final_message := header_message + "\r\n" + body_message.String()
	status := smtp.SendMail(host, auth, from_email, []string{to_email}, []byte(final_message))
	if status != nil {
		fmt.Printf("Error from SMTP Server: %s \n", status)
		return false, "Error from SMTP Server: " + status.Error()
	} else {
		fmt.Println("Email Sent Successfully!!")
		return true, ""
	}

} //-----------------------------------------------------------------------------------------------------
