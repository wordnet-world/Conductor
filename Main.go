package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/wordnet-world/Conductor/service"
)

// PORT to expose server
const PORT = 8675

func main() {
	log.SetFlags(log.LUTC | log.Llongfile | log.Ldate | log.Ltime)
	log.Printf("Starting on port %d\n", PORT)
	router := service.NewRouter()

	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(PORT), router))
}
