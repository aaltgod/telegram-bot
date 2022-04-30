package bot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Api struct {
	logger *logrus.Logger
	client *http.Client
	uri    string
}

type (
	User struct {
		Name    string `json:"username"`
		ID      int64  `json:"id"`
		IsAdmin bool   `json:"is_admin"`
	}

	Users []User

	CreateUser struct {
		Name    string `json:"username"`
		ID      int64  `json:"id"`
		IsAdmin bool   `json:"is_admin"`
	}

	UpdateUser struct {
		IsAdmin bool `json:"is_admin"`
	}

	Request struct {
		IP       string `json:"ip"`
		Response string `json:"response"`
	}

	Requests []Request
)

func NewApi(logger *logrus.Logger, uri string) *Api {
	api := &Api{
		logger: logger,
		uri:    uri,
	}
	api.client = &http.Client{
		Timeout: time.Second * 10,
	}

	return api
}

func (api *Api) GetUser(id int64) (*User, error) {

	u := &User{}

	req := &http.Request{
		Method: http.MethodGet,
	}

	req.URL, _ = url.Parse(api.uri + "/users/" + strconv.Itoa(int(id)))

	resp, err := api.client.Do(req)
	if err != nil {
		return u, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return u, err
	}

	if err := json.Unmarshal(data, u); err != nil {
		return u, err
	}

	return u, nil
}

func (api *Api) InsertUser(u *CreateUser) error {

	result, err := json.Marshal(u)
	if err != nil {
		return err
	}

	url := api.uri + "/users/"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(result))
	if err != nil {
		return err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	api.logger.Infoln(string(data))

	return nil
}

func (api *Api) UpdateUser(id int64, u *UpdateUser) error {

	result, err := json.Marshal(u)
	if err != nil {
		return err
	}

	url := api.uri + "/users/" + strconv.Itoa(int(id))

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(result))
	if err != nil {
		return err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	api.logger.Infoln(string(data))

	return nil
}

func (api *Api) GetUsers() ([]User, error) {

	users := []User{}

	req := &http.Request{
		Method: http.MethodGet,
	}

	req.URL, _ = url.Parse(api.uri + "/users")

	resp, err := api.client.Do(req)
	if err != nil {
		return users, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return users, err
	}

	usersToUnmarshal := &Users{}

	if err := json.Unmarshal(data, usersToUnmarshal); err != nil {
		return users, err
	}

	for _, u := range *usersToUnmarshal {
		users = append(users, u)
	}

	return users, nil
}

func (api *Api) AppendRequest(id int64, r *Request) error {

	result, err := json.Marshal(r)
	if err != nil {
		return err
	}

	url := api.uri + "/requests/" + strconv.Itoa(int(id))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(result))
	if err != nil {
		return err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	api.logger.Infoln(string(data))

	return nil
}

func (api *Api) GetAllRequestsByID(id int64) ([]Request, error) {

	requests := []Request{}

	req := &http.Request{
		Method: http.MethodGet,
	}

	req.URL, _ = url.Parse(api.uri + "/requests/" + strconv.Itoa(int(id)))

	resp, err := api.client.Do(req)
	if err != nil {
		return requests, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return requests, err
	}

	requestsToUnmarshal := &Requests{}

	if err := json.Unmarshal(data, requestsToUnmarshal); err != nil {
		return requests, err
	}

	for _, u := range *requestsToUnmarshal {
		requests = append(requests, u)
	}

	return requests, nil
}
