package main

import (
	"fmt"
	"net/http"

	"./src"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

//Soporte TLS

func handleRequests() {

	defer func() {
		recuperado := recover()
		if recuperado != nil {
			res := "handleRequests Error:" + fmt.Sprint(recuperado)
			fmt.Println(res)
		}
	}()

	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter := mux.NewRouter()
	myRouter.HandleFunc("/{div}", hello).Queries("idp", "{idp}", "key", "{key}").Methods("GET")
	myRouter.HandleFunc("/api/v1/test", src.TestAuth0).Methods("GET")
	myRouter.HandleFunc("/api/v1/post_msg", postmsg).Methods("POST")
	//Distribuidores
	myRouter.HandleFunc("/api/v1/get_distributor", src.GetDistributor).Queries("myidp", "{myidp}", "idp", "{idp}").Methods("GET")
	myRouter.HandleFunc("/api/v1/get_distributors", src.GetDistributors).Queries("myidp", "{myidp}").Methods("GET")
	myRouter.HandleFunc("/api/v1/create_distributor", src.CreateDistributor).Methods("POST")
	myRouter.HandleFunc("/api/v1/edit_distributor", src.EditDistributor).Queries("myidp", "{myidp}", "idp", "{idp}").Methods("POST")
	myRouter.HandleFunc("/api/v1/get_distristatusaccount", src.GetDistriStatusAccount).Queries("myidp", "{myidp}").Methods("GET")
	//Usuarios
	myRouter.HandleFunc("/api/v1/create_user", src.CreateUser).Methods("POST")
	myRouter.HandleFunc("/api/v1/get_user", src.GetUser).Queries("myidp", "{myidp}", "id", "{id}").Methods("GET")
	myRouter.HandleFunc("/api/v1/get_log_user", src.GetLogUser).Queries("authid", "{authid}").Methods("GET")
	myRouter.HandleFunc("/api/v1/get_users", src.GetUsers).Queries("myidp", "{myidp}").Methods("GET")
	myRouter.HandleFunc("/api/v1/edit_user", src.EditUser).Queries("myidp", "{myidp}", "id", "{id}").Methods("POST")
	//Clientes
	myRouter.HandleFunc("/api/v1/get_client", src.GetClient).Queries("myidp", "{myidp}", "idp", "{idp}").Methods("GET")
	myRouter.HandleFunc("/api/v1/get_clients", src.GetClients).Queries("myidp", "{myidp}").Methods("GET")
	myRouter.HandleFunc("/api/v1/create_client", src.CreateClient).Methods("POST")
	myRouter.HandleFunc("/api/v1/edit_client", src.EditClient).Queries("myidp", "{myidp}", "idp", "{idp}").Methods("POST")
	//Servicios
	myRouter.HandleFunc("/api/v1/get_service", src.GetService).Queries("myidp", "{myidp}", "idp", "{idp}").Methods("GET")
	myRouter.HandleFunc("/api/v1/get_services", src.GetServices).Queries("myidp", "{myidp}").Methods("GET")
	myRouter.HandleFunc("/api/v1/create_service", src.CreateService).Methods("POST")
	myRouter.HandleFunc("/api/v2/create_service", src.CreateServiceX).Methods("POST")
	myRouter.HandleFunc("/api/v1/edit_service", src.EditService).Queries("myidp", "{myidp}", "idp", "{idp}").Methods("POST")
	myRouter.HandleFunc("/api/v1/get_planes", src.GetPlanes).Queries("myidp", "{myidp}", "tipo", "{tipo}").Methods("GET")
	myRouter.HandleFunc("/api/v1/get_serialS", src.GetSerialS).Queries("serial", "{serial}").Methods("GET")
	//Registro por Serial Number
	myRouter.HandleFunc("/api/v1/get_serial", src.GetSerial).Queries("serial", "{serial}", "email", "{email}").Methods("GET")
	myRouter.HandleFunc("/api/v1/create_cli_service", src.CreateCliService).Methods("POST")
	myRouter.HandleFunc("/api/v1/create_equipo", src.CreateEquipo).Methods("POST")
	myRouter.HandleFunc("/api/v1/val_serial", src.ValSerial).Queries("serial", "{serial}").Methods("GET")
	myRouter.HandleFunc("/api/v1/create_serial", src.CreateSerial).Methods("POST")
	//PORTS  SandBox:8081  Produccion:8085
	//log.Fatal(http.ListenAndServe(":8085", myRouter))
	err := http.ListenAndServe(":8085", myRouter) //Produccion 8090 SandBox:8092
	if err != nil {
		panic(err)
	}

}

func main() {

	fmt.Println("WiseLink GO-API V1!")

	src.InitDB()
	fmt.Println("Hora Local: " + src.Timestamp_now())

	handleRequests()
}
