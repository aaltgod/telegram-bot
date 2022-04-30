package postgresql

import (
	"github.com/aaltgod/telegram-bot/internal/storage-service/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Storage struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewStorage(logger *logrus.Logger, db *gorm.DB) *Storage {
	return &Storage{
		logger: logger,
		db:     db,
	}
}

func (st *Storage) Migrate() error {
	st.logger.Infoln("migrate")
	return st.db.AutoMigrate(&repository.User{}, &repository.Request{})
}

func (st *Storage) InitAdmin(name string, id int64) error {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "initadmin",
		"args":   []interface{}{name, id},
	}).Infoln()

	createAdmin := &repository.CreateUser{
		Name:    name,
		ID:      id,
		IsAdmin: true,
	}

	return st.Insert(createAdmin)
}

func (st *Storage) Get(id int64) (*repository.User, error) {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "get",
		"args":   id,
	}).Infoln()

	u := &repository.User{}

	st.db.Where("id = ?", id).Find(u)

	return u, nil
}

func (st *Storage) Insert(u *repository.CreateUser) error {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "insert",
		"args":   *u,
	}).Infoln()

	user, err := st.Get(u.ID)
	if err != nil {
		return err
	}

	if user.Name != "" {
		st.logger.Infof("user [%s] exists\n", user.Name)
		return nil
	}

	createUser := &repository.User{
		Name:    u.Name,
		ID:      u.ID,
		IsAdmin: u.IsAdmin,
	}

	tx := st.db.Create(createUser)
	if tx.Error != nil {
		st.logger.Warnln(tx.Error)
		return tx.Error
	}

	return nil
}

func (st *Storage) Update(id int64, u *repository.UpdateUser) error {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "update",
		"args":   []interface{}{id, *u},
	}).Infoln()

	user, err := st.Get(id)
	if err != nil {
		return err
	}

	if user.Name == "" {
		st.logger.Infof("user [%s] doesn't exist\n", user.Name)
		return nil
	}

	tx := st.db.Model(&repository.User{}).Where("id = ?", id).Update("is_admin", u.IsAdmin)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (st *Storage) GetAll() ([]*repository.User, error) {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "getall",
		"args":   "",
	}).Infoln()

	var users []*repository.User

	if err := st.db.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (st *Storage) AppendRequest(id int64, r *repository.Request) error {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "appendrequest",
		"args":   []interface{}{id, *r},
	}).Infoln()

	user, err := st.Get(id)
	if err != nil {
		return err
	}

	requests := []*repository.Request{}
	if err := st.db.Model(user).Association("Requests").Find(&requests); err != nil {
		return err
	}

	for _, req := range requests {
		if req.IP == r.IP {
			return nil
		}
	}

	if err := st.db.Model(user).Association("Requests").Append(r); err != nil {
		return err
	}

	return nil
}

func (st *Storage) DeleteRequest(id int64, r *repository.DeleteRequest) error {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "deleterequest",
		"args":   []interface{}{id, *r},
	}).Infoln()

	user, err := st.Get(id)
	if err != nil {
		return err
	}

	req := &repository.Request{}

	err = st.db.Model(user).Where("ip = ?", r.IP).Association("Requests").Find(req)
	if err != nil {
		return err
	}
	if req.Response == "" {
		return nil
	}

	return st.db.Unscoped().Delete(req).Error
}

func (st *Storage) GetAllRequestsByID(id int64) ([]repository.Request, error) {
	st.logger.WithFields(logrus.Fields{
		"struct": "storage",
		"method": "getallrequestsbyid",
		"args":   id,
	}).Infoln()

	var requests []repository.Request

	user, err := st.Get(id)
	if err != nil {
		return requests, err
	}

	if err := st.db.Model(user).Association("Requests").Find(&requests); err != nil {
		return requests, err
	}

	return requests, nil
}
