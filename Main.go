package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/wordnet-world/Conductor/database"

	"github.com/wordnet-world/Conductor/service"
)

// PORT to expose server
const PORT = 8675

func main() {
	log.SetFlags(log.LUTC | log.Llongfile | log.Ldate | log.Ltime)
	log.Printf("Starting on port %d\n", PORT)

	// flush/setup DB
	rdb := database.GetCacheDatabase()
	rdb.SetupDB()

	res, err := database.HelloWorld("bolt://127.0.0.1:7687", "neo4j", "neo4j1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	// start router to allow connections
	router := service.NewRouter()

	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(PORT), router))
}
