package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/aaltgod/telegram-bot/internal/storage-service/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger  *logrus.Logger
	storage repository.Repository
}

type User struct {
	Name    string `json:"username"`
	ID      int64  `json:"id"`
	IsAdmin bool   `json:"is_admin"`
}

type UpdateUser struct {
	IsAdmin bool `json:"is_admin"`
}

type Request struct {
	IP       string `json:"ip"`
	Response string `json:"response"`
}

type Requests []Request

func NewHandler(logger *logrus.Logger, storage repository.Repository) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
	}
}

func (h *Handler) GetUsers(c echo.Context) error {
	h.logger.WithFields(logrus.Fields{
		"struct": "handler",
		"method": "getusers",
		"path":   c.Request().URL,
	}).Infoln()

	users, err := h.storage.GetAll()
	if err != nil {
		return err
	}

	usersToResponse := make([]User, 0, len(users))
	for _, u := range users {
		usersToResponse = append(usersToResponse, User{
			Name:    u.Name,
			ID:      u.ID,
			IsAdmin: u.IsAdmin,
		})
	}

	return c.JSON(http.StatusOK, usersToResponse)
}

func (h *Handler) GetUser(c echo.Context) error {
	h.logger.WithFields(logrus.Fields{
		"struct": "handler",
		"method": "getuser",
		"path":   c.Request().URL,
	}).Infoln()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	user, err := h.storage.Get(int64(id))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	userToResponse := &User{
		Name:    user.Name,
		ID:      user.ID,
		IsAdmin: user.IsAdmin,
	}

	return c.JSON(http.StatusOK, userToResponse)
}

func (h *Handler) InsertUser(c echo.Context) error {
	h.logger.WithFields(logrus.Fields{
		"struct": "handler",
		"method": "insertuser",
		"path":   c.Request().URL,
	}).Infoln()

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	u := &User{}

	if err := json.Unmarshal(data, u); err != nil {
		return err
	}

	createUser := &repository.CreateUser{
		Name:    u.Name,
		ID:      u.ID,
		IsAdmin: u.IsAdmin,
	}

	if err := h.storage.Insert(createUser); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "ok")
}

func (h *Handler) UpdateUser(c echo.Context) error {
	h.logger.WithFields(logrus.Fields{
		"struct": "handler",
		"method": "updateuser",
		"path":   c.Request().URL,
	}).Infoln()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	userToUnmarshal := &UpdateUser{}

	if err := json.Unmarshal(data, userToUnmarshal); err != nil {
		return err
	}

	u := &repository.UpdateUser{
		IsAdmin: userToUnmarshal.IsAdmin,
	}

	if err := h.storage.Update(int64(id), u); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "updated")
}

func (h *Handler) AppendRequest(c echo.Context) error {
	h.logger.WithFields(logrus.Fields{
		"struct": "handler",
		"method": "appendrequest",
		"path":   c.Request().URL,
	}).Infoln()

	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	requestToUnmarshal := &Request{}

	if err := json.Unmarshal(data, requestToUnmarshal); err != nil {
		return err
	}

	u := &repository.Request{
		IP:       requestToUnmarshal.IP,
		Response: requestToUnmarshal.Response,
	}

	if err := h.storage.AppendRequest(int64(id), u); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "appended")
}

func (h *Handler) DeleteRequest(c echo.Context) error {
	h.logger.WithFields(logrus.Fields{
		"struct": "handler",
		"method": "deleterequest",
		"path":   c.Request().URL,
	}).Infoln()

	id, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	ip := c.QueryParam("ip")

	u := &repository.DeleteRequest{
		IP: ip,
	}

	if err := h.storage.DeleteRequest(int64(id), u); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "deleted")
}

func (h *Handler) GetAllRequestsByID(c echo.Context) error {
	h.logger.WithFields(logrus.Fields{
		"struct": "handler",
		"method": "getallrequestsbyid",
		"path":   c.Request().URL,
	}).Infoln()

	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.logger.Warnln(err)
		return err
	}

	requests, err := h.storage.GetAllRequestsByID(int64(id))
	if err != nil {
		return err
	}

	requestsToResponse := Requests{}
	for _, u := range requests {
		requestsToResponse = append(requestsToResponse, Request{
			IP:       u.IP,
			Response: u.Response,
		})
	}

	return c.JSON(http.StatusOK, requestsToResponse)
}
