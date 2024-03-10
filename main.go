package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

var Users = make(map[string]User)

func main() {
	r := chi.NewRouter()

	r.Post("/create", CreateUser)
	r.Post("/make_friends", MakeFriends)
	r.Delete("/user", DeleteUser)
	r.Get("/friends/{user_id}", GetFriends)
	r.Put("/user/{user_id}", UpdateAge)

	http.ListenAndServe(":8080", r)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = strconv.Itoa(len(Users) + 1)

	Users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user.ID)
}

func MakeFriends(w http.ResponseWriter, r *http.Request) {
	var ids struct {
		SourceID string `json:"source_id"`
		TargetID string `json:"target_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sourceUser, sourceExists := Users[ids.SourceID]
	targetUser, targetExists := Users[ids.TargetID]

	if !sourceExists || !targetExists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	sourceUser.Friends = append(sourceUser.Friends, targetUser.ID)
	targetUser.Friends = append(targetUser.Friends, sourceUser.ID)

	Users[sourceUser.ID] = sourceUser
	Users[targetUser.ID] = targetUser

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sourceUser.Name + " and " + targetUser.Name + " are now friends")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var id struct {
		TargetID string `json:"target_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := Users[id.TargetID]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(Users, id.TargetID)

	for _, friend := range user.Friends {
		friendUser := Users[friend]
		for i, friendID := range friendUser.Friends {
			if friendID == user.ID {
				friendUser.Friends = append(friendUser.Friends[:i], friendUser.Friends[i+1:]...)
				break
			}
		}
		Users[friend] = friendUser
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.Name + " has been deleted")
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	user, exists := Users[userID]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user.Friends)
}

func UpdateAge(w http.ResponseWriter, r *http.Request) {
	var newAge struct {
		NewAge string `json:"new_age"`
	}

	err := json.NewDecoder(r.Body).Decode(&newAge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := chi.URLParam(r, "user_id")

	user, exists := Users[userID]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Age = newAge.NewAge
	Users[userID] = user

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User's age has been successfully updated")
}
