package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger *logrus.Logger
	client *http.Client
	uri    string
}

type User struct {
	Name    string `json:"username"`
	ID      int64  `json:"id"`
	IsAdmin bool   `json:"is_admin"`
}

type Users []User

type Request struct {
	IP       string `json:"ip"`
	Response string `json:"response"`
}

type Requests []Request

func NewHandler(logger *logrus.Logger, uri string) *Handler {
	h := &Handler{
		logger: logger,
		uri:    uri,
	}
	h.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return h
}

func (h *Handler) GetUser(c echo.Context) error {

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	req := &http.Request{
		Method: http.MethodGet,
	}

	req.URL, _ = url.Parse(h.uri + "/users/" + strconv.Itoa(int(id)))

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	u := &User{}

	if err := json.Unmarshal(data, u); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, u)
}

func (h *Handler) GetUserHistory(c echo.Context) error {

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	req := &http.Request{
		Method: http.MethodGet,
	}

	req.URL, _ = url.Parse(h.uri + "/requests/" + strconv.Itoa(int(id)))

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	requestsToUnmarshal := &Requests{}

	if err := json.Unmarshal(data, requestsToUnmarshal); err != nil {
		return err
	}

	requests := []Request{}

	for _, u := range *requestsToUnmarshal {
		requests = append(requests, u)
	}

	return c.JSON(http.StatusOK, requests)
}

func (h *Handler) DeleteOneRequest(c echo.Context) error {

	id, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	ip := c.QueryParam("ip")

	req := &http.Request{
		Method: http.MethodDelete,
	}

	req.URL, _ = url.Parse(h.uri + "/requests?user_id=" + strconv.Itoa(int(id)) + "&ip=" + ip)

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	h.logger.Infoln(string(data))

	return c.JSON(http.StatusOK, "deleted")
}

func (h *Handler) GetUsers(c echo.Context) error {

	req := &http.Request{
		Method: http.MethodGet,
	}

	req.URL, _ = url.Parse(h.uri + "/users")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	usersToUnmarshal := &Users{}

	if err := json.Unmarshal(data, usersToUnmarshal); err != nil {
		return err
	}

	users := []User{}

	for _, u := range *usersToUnmarshal {
		users = append(users, u)
	}

	return c.JSON(http.StatusOK, users)
}
