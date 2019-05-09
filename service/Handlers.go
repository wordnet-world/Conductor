package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wordnet-world/Conductor/database"
	"github.com/wordnet-world/Conductor/models"
	//"github.com/google/go-cmp/cmp"
)

// HeartBeat end point to verify connection
func HeartBeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I'm Alive!")
}

// AdminPasswordCheck determines if the client can access the admin pages
func AdminPasswordCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Remember to defer some recovery code here
}

// JoinGame this will be fun, will need to return a websocket
func JoinGame(w http.ResponseWriter, r *http.Request) {

}

// CreateGame will create a game with the specified configuration
func CreateGame(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	// Check admin password
	verifyPassword(r)

	db := database.GetCacheDatabase()

	game := models.CreateGame{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln("Could not read the body of the message")
	}
	log.Printf("Received body: %s\n", string(body))
	err = json.Unmarshal(body, &game)
	if err != nil {
		log.Panicln("Could not unmarshall body into game object")
	}
	log.Printf("This is the json %v\n", game)

	// TODO use graph database to get a random root node
	// from the graph for the start node
	game.StartNode = "startNode123"
	gameID := db.CreateGame(game)
	fmt.Fprintf(w, models.CreateHTTPResponse(nil, map[string]interface{}{"gameID": gameID}, true).ToJSON())
}

// DeleteGame will delete the game with the matching id
func DeleteGame(w http.ResponseWriter, r *http.Request) {

}

// ListGames will return an array of games
// For convenience I will return everything if the no fields are specified
func ListGames(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()
	// TODO Will need to have special handling if the string Teams is specified in fields

	// Check admin password
	verifyPassword(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln("Could not read the body of the message")
	}
	log.Printf("Received body: %s\n", string(body))

	fields := models.ListGameFields{}
	err = json.Unmarshal(body, &fields)
	if err != nil {
		log.Panicln("Could not unmarshall body into fields object")
	}
	log.Printf("This is the json %v\n", fields)

	db := database.GetCacheDatabase()

	if len(fields.Fields) == 0 {
		games := db.GetGames()
		log.Printf("Here are the games%v\n", games)
		fmt.Fprintln(w, models.CreateHTTPResponse("blah", nil, true).ToJSON())
	} else {
		// do later, handle sending only a subset
		fmt.Fprintln(w, models.CreateHTTPResponse("blah", nil, true).ToJSON())
	}

}

func verifyPassword(r *http.Request) {
	adminPassword := r.Header["Adminpassword"] // TODO: Need to capitalize P and hopefully it'll still work, Postman sends this
	if len(adminPassword) != 1 {
		log.Panicln("Malformed header 'AdminPassword'")
	} else if adminPassword[0] != models.Config.Wordnet.AdminPassword {
		log.Panicln("Incorrect Admin Password")
	}
}

// Store handles POST requests to /store
// This converts the object to a checkoff model and
// sends it to long term storage
/*func Store(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, false).ToJSON())
		}
	}()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Received in body: %v\n", string(body))
	activeUserModel := models.ActiveUserModel{}
	json.Unmarshal(body, &activeUserModel)
	log.Printf("Unmarshalled response %v\n", activeUserModel)

	if cmp.Equal(activeUserModel, models.ActiveUserModel{}) {
		log.Panic("Incorrect format of Body")
	}

	db := database.GetDriver()
	db.Store(models.CreateCheckoff(activeUserModel))
	log.Println("Stored in DB")
	fmt.Fprintln(w, models.CreateHTTPResponse(nil, true).ToJSON())
}

// CSV offers a .csv file for download
func CSV(w http.ResponseWriter, r *http.Request) {
	modtime := time.Now()

	w.Header().Add("Content-Disposition", "Attachment")

	db := database.GetDriver()
	csvString := db.GenerateCSV()

	http.ServeContent(w, r, "random.csv", modtime, strings.NewReader(csvString))
} */
