package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{listenAddr: listenAddr, store: store}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.Handle("/account/{id}", makeHandler(s.handleAccount))
	router.Handle("/account", makeHandler(s.handleAccount))

	log.Println("Server started at ", s.listenAddr, " ...")
	log.Fatal(http.ListenAndServe(s.listenAddr, router))

}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handleCreateAccount(w, r)
	case "PUT":
		return s.handleEditAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	case "GET":
		return s.handleGetAccount(w, r)

	}

	return nil
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	reqJson := CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqJson)
	if err != nil {
		return err
	}

	newAcc := Account{
		ID:        0,
		FirstName: reqJson.FirstName,
		LastName:  reqJson.LastName,
		Number:    int64(rand.Intn(100)),
		Balance:   0,
		CreatedAt: time.Now().UTC(),
	}

	// make database query
	res, err := s.store.CreateAccount(&newAcc)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, fmt.Sprintf("New account created with ID: %d", res))
	return nil
}

func (s *ApiServer) handleEditAccount(w http.ResponseWriter, r *http.Request) error {
	strId := mux.Vars(r)["id"]
	if strId == "" {
		return fmt.Errorf("an ID is required in the PUT request")
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		return err
	}
	editReq := CreateAccountRequest{}
	decErr := json.NewDecoder(r.Body).Decode(&editReq)
	if decErr != nil {
		return err
	}

	updatedAcc := Account{
		ID:        id,
		FirstName: editReq.FirstName,
		LastName:  editReq.LastName,
	}
	updateErr := s.store.UpdateAccount(&updatedAcc)

	if updateErr != nil {
		WriteJson(w, http.StatusOK, fmt.Sprintf("Unable to update account with ID: %d", id))
	} else {

		WriteJson(w, http.StatusOK, updatedAcc)
	}
	return nil

}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	strId := vars["id"]

	if strId == "" {
		return fmt.Errorf("an ID needs to be provided to the DELETE method")
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		return err
	}

	deletedId := s.store.DeleteAccount(id)
	if deletedId == -1 {
		WriteJson(w, http.StatusOK, fmt.Sprintf("No account delete with the ID: %d", id))
	} else {

		WriteJson(w, http.StatusOK, fmt.Sprintf("Account with ID: %d successfully deleted", deletedId))
	}

	return nil
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	strId := mux.Vars(r)["id"]

	// handle case of GET request with ID
	if strId != "" {
		id, err := strconv.Atoi(strId)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			WriteJson(w, http.StatusOK, fmt.Sprintf("Unable to find account with ID: %d", id))
		} else {
			WriteJson(w, http.StatusOK, account)
		}

	} else {
		// Handle case of GET request with no ID, ie. List
		accounts, err := s.store.GetAccounts()
		if err != nil {
			return err
		}
		WriteJson(w, http.StatusOK, accounts)
	}

	return nil
}

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func WriteJson(w http.ResponseWriter, status int, value any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func makeHandler(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Println(err)
			WriteJson(w, http.StatusBadRequest, err)
		}
	}
}
