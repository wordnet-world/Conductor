package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/wordnet-world/Conductor/database"
	"github.com/wordnet-world/Conductor/models"
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

	graph := database.GetGraphDatabase()
	err := graph.PopulateDummy(models.Config.Neo4j.URI, models.Config.Neo4j.Username, models.Config.Neo4j.Password)
	if err != nil {
		log.Println(err)
	}
	graph.Close()

	// start router to allow connections
	router := service.NewRouter()

	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(PORT), router))
}
