package post

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nidal-Bakir/first_go/pkg/user"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request, user user.UserModel) {
	// do some work to create the post

	type doneStatus struct {
		Status string `json:"status"`
	}

	data, err := json.Marshal(doneStatus{Status: "done"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
