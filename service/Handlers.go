package service

import (
	"fmt"
	"net/http"

	"github.com/wordnet-world/Conductor/models"

	//"github.com/google/go-cmp/cmp"
	"github.com/wordnet-world/Conductor/database"
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
	game := models.Game{}
	db := database.GetDatabase()
	db.CreateGame(game)
	fmt.Fprintln(w, "Did the thing")
}

// DeleteGame will delete the game with the matching id
func DeleteGame(w http.ResponseWriter, r *http.Request) {

}

// ListGames will return an array of games
func ListGames(w http.ResponseWriter, r *http.Request) {

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
