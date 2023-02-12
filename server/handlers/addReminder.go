package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kachamaka/chaosgo/database"
	"github.com/kachamaka/chaosgo/models"
	"github.com/kachamaka/chaosgo/status"
)

func AddReminderHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	stringID, objectID, err := database.GetHeaders(r)
	if err != nil {
		encoder.Encode(models.BasicResponse{Success: false, Message: err.Error(), Status: status.AUTHORIZATION_ERROR})
		log.Println("error by get headers: ", err)
		return
	}

	var req models.ReminderRequest
	if err := decoder.Decode(&req); err != nil {
		encoder.Encode(models.BasicResponse{Success: false, Message: "can't read body", Status: status.BODY_ERROR})
		log.Println("error by decode body: ", err)
		return
	}

	user, err := database.GetUser(objectID)
	if err != nil {
		encoder.Encode(models.BasicResponse{Success: false, Message: err.Error(), Status: status.GET_USER_ERROR})
		log.Println("err by get user: ", err)
		return
	}

	reminder := models.Reminder{
		UserID:     stringID,
		Email:      user.Email,
		Subject:    req.Subject,
		Time:       req.EventStart - req.TimeAhead,
		EventStart: req.EventStart,
	}

	err = database.Get().AddReminder(reminder)
	if err != nil {
		encoder.Encode(models.BasicResponse{Success: false, Message: err.Error(), Status: status.ADD_REMINDER_ERROR})
		log.Println("err by add reminder: ", err)
		return
	}

	// reminder.Time = time.Now().Add(time.Second * 30).Unix()

	//send reminder
	// fmt.Println(reminder)
	go database.Send(reminder)

	encoder.Encode(models.BasicResponse{Success: true, Message: "reminder added successfully", Status: status.OK})
}
