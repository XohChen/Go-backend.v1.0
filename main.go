package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Informations struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Freindes []int  `json:"freindes"`
}

var infos []Informations

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/users", addUser).Methods("POST")
	router.HandleFunc("/api/users/update/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/users/removeUser/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/api/users/addFr/{id}", addFreinds).Methods("POST")
	router.HandleFunc("/api/users/deleteFr/{id}", deleteFreinds).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(infos)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, info := range infos {
		if info.ID == id {
			json.NewEncoder(w).Encode(info)
			//re-view-1-bug getUser: always returns not found
			return
		}
	}
	http.NotFound(w, r)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var info Informations
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//re-view --> add user validation: name, phone, friends
	//1- name validation
	if info.Name == "" {
		http.Error(w, "#ERROR - Validation in `Name` feild, should not empty", http.StatusBadRequest)
		return
	}

	//2- Phone validation

	if info.Phone == "" {
		http.Error(w, "#ERROR - Validation in `Phone` feild, should not empty", http.StatusBadRequest)
		return

	}

	//3- Freindes validation
	if info.Freindes == nil {
		http.Error(w, "#ERROR - Validation in `Freindes` feild, should not empty", http.StatusBadRequest)
		return

	}
	//test
	// add user validation: friends must be existed
	for i := 0; i < len(info.Freindes); i++ {
		var found bool
		for _, j := range infos {
			ok := info.Freindes[i] == j.ID
			if ok {
				found = true
				break
			}
		}
		if !found {
			http.Error(w, "#ERROR - Frindes not exists", http.StatusBadRequest)
			return
		}
	}

	info.ID = len(infos) + 1
	infos = append(infos, info)
	json.NewEncoder(w).Encode(info)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedUser Informations
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, info := range infos {
		if info.ID == id {
			// re-view-2 -> a user should not be able to update it's id
			updatedUser.ID = id
			infos[i] = updatedUser
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}

	http.NotFound(w, r)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, info := range infos {
		if info.ID == id {
			infos = append(infos[:i], infos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.NotFound(w, r)
}

// re-view2- add freindes
func addFreinds(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var addFreindesOfUsers Informations
	id, errId := strconv.Atoi(params["id"])

	if errId != nil {
		http.Error(w, errId.Error(), http.StatusBadRequest)
		return
	}

	errId = json.NewDecoder(r.Body).Decode(&addFreindesOfUsers)
	if errId != nil {
		http.Error(w, errId.Error(), http.StatusBadRequest)
	}

	for i, info := range infos {
		if info.ID == id {
			// note: should be check freinde id not exists in user freind list

			//
			addFreindesOfUsers.ID = id
			infos[i].Freindes = append(infos[i].Freindes, addFreindesOfUsers.Freindes...)
			json.NewEncoder(w).Encode(addFreindesOfUsers)
			return
		}
	}

	http.NotFound(w, r)
}

func deleteFreinds(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var deleteFreindesOfUsers Informations
	err = json.NewDecoder(r.Body).Decode(&deleteFreindesOfUsers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for j, info := range infos {
		if info.ID == id {
			for i := 0; i < len(deleteFreindesOfUsers.Freindes); i++ {
				index := getIndex(info.Freindes, deleteFreindesOfUsers.Freindes[i])
				if index == -1 {
					http.Error(w, "Index of array not found", http.StatusBadRequest)
					continue
				}
				log.Printf("%d", index)
				deleteElement(&infos[j], index)
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.NotFound(w, r)
}

func deleteElement(me *Informations, indexOFfreindToRemove int) {
	a := me.Freindes[:indexOFfreindToRemove]
	b := me.Freindes[indexOFfreindToRemove+1:]
	x := append(a, b...)
	me.Freindes = x
}

func getIndex(arrayOfInfo []int, freindID int) int {
	for i, arrInfo := range arrayOfInfo {
		if arrInfo == freindID {
			return i
		}
	}
	return -1
}
