package service

import (
	"fmt"
	"net/http"
	//"github.com/google/go-cmp/cmp"
)

// End point to verify connection
func HeartBeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func AdminPasswordCheck(w http.ResponseWriter, r *http.Request) {

}

func JoinGame(w http.ResponseWriter, r *http.Request) {

}

func CreateGame(w http.ResponseWriter, r *http.Request) {

}

func DeleteGame(w http.ResponseWriter, r *http.Request) {

}

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
