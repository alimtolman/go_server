package main

import (
	"Server/DataBase"
	"Server/DataBase/Tables"
	"Server/Helpers/Convert"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const http_port = "80"
const https_port = "443"
const cert_path = "_cert.pem"
const key_path = "_key.pem"
const root_msg = "It Works!"

/* AppData */

func AppDataAdd(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	q_id := request.URL.Query().Get("id")
	id := Convert.StrToUInt(q_id, 0)
	item := DataBase.NewAppData(id, string(body))

	Tables.AppData.AddOrSet(item)
	fmt.Fprint(writer, "OK")
}

func AppDataClear(writer http.ResponseWriter, request *http.Request) {
	Tables.AppData.RemoveAll()
	fmt.Fprint(writer, "OK")
}

func AppDataGet(writer http.ResponseWriter, request *http.Request) {
	q_id := request.URL.Query().Get("id")
	id := Convert.StrToUInt(q_id, 0)
	data := Tables.AppData.Get(func(item DataBase.AppData) bool { return item.Id == id })

	fmt.Fprint(writer, data.Data)
}

/* Root */

func RootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		fmt.Fprint(writer, root_msg)
	} else {
		http.NotFound(writer, request)
	}
}

func main() {
	var waitGroup sync.WaitGroup

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/api/data/add", AppDataAdd)
	http.HandleFunc("/api/data/clear", AppDataClear)
	http.HandleFunc("/api/data/get", AppDataGet)

	waitGroup.Add(1)
	waitGroup.Add(1)

	go func() {
		if err := http.ListenAndServe(":"+http_port, nil); err != nil {
			log.Fatal("Server listen error:", err)
		}

		waitGroup.Done()
	}()

	go func() {
		if err := http.ListenAndServeTLS(":"+https_port, cert_path, key_path, nil); err != nil {
			log.Fatal("Server tls listen error:", err)
		}

		waitGroup.Done()
	}()

	waitGroup.Wait()
}
