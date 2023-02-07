package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kachamaka/chaosgo/database"
	"github.com/kachamaka/chaosgo/models"
)

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	stringID, err := database.GetHeader(r)
	if err != nil {
		encoder.Encode(models.BasicResponse{Success: false, Message: "id from auth header not string"})
		log.Println("err by getting id from auth header")
		return
	}

	var event models.Event
	if err := decoder.Decode(&event); err != nil {
		encoder.Encode(models.BasicResponse{Success: false, Message: "can't read body"})
		log.Println("err by decoding body: ", err)
		return
	}

	err = database.Get().DeleteEvent(stringID, event)
	if err != nil {
		encoder.Encode(models.BasicResponse{Success: false, Message: err.Error()})
		log.Println("err deleting event: ", err)
		return
	}

	encoder.Encode(models.BasicResponse{Success: true, Message: "event deleted successfully"})
}
