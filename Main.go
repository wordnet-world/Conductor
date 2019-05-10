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

	graph := database.GetGraphDatabase()
	err := graph.Connect("bolt://neo4j:7687", "neo4j", "neo4j")
	if err != nil {
		log.Println(err)
	} else {
		root, _ := graph.GetRoot()
		nodes, _ := graph.GetNeighbors(root)
		fmt.Println(root)
		fmt.Println(nodes)
	}

	// start router to allow connections
	router := service.NewRouter()

	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(PORT), router))
}
