package main

import (
	"../common"
	"encoding/json"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
	var videos []common.Video

	db := common.Init()
	db.Find(&videos)
	json.NewEncoder(w).Encode(videos)
}

func CreateNew(w http.ResponseWriter, r *http.Request) {
	var video common.Video
	db := common.Init()

	body := json.NewDecoder(r.Body)
	err := body.Decode(&video)

	if err != nil {
		panic(err)
	}

	uuid, _ := uuid.NewRandom()
	video.JobId = uuid
	t := time.Now()
	video.Created = &t
	video.Status = "CREATED"
	db.Create(&video)
	server, error := common.MachineryInstance()
	if error != nil {
		panic("Could not create server")
	}
	transcode_task := &tasks.Signature{
		Name:       "Transcode",
		RetryCount: 3,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: uuid,
			},
		},
	}

	_, err1 := server.SendTask(transcode_task)

	if err1 != nil {
		panic("Could not send transcode task")
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"JobID": uuid.String()})
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	var video common.Video
	params := mux.Vars(r)
	db := common.Init()
	db.First(&video, common.Video{JobId: uuid.MustParse(params["id"])})
	json.NewEncoder(w).Encode(video)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/jobs", GetJobs).Methods("GET")
	router.HandleFunc("/api/v1/jobs", CreateNew).Methods("POST")
	router.HandleFunc("/api/v1/jobs/{id}", GetJob).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}
