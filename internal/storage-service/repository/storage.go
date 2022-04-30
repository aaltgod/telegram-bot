package repository

type Repository interface {
	Get(id int64) (*User, error)
	Insert(u *CreateUser) error
	Update(id int64, u *UpdateUser) error
	GetAll() ([]*User, error)

	AppendRequest(id int64, r *Request) error
	DeleteRequest(id int64, r *DeleteRequest) error
	GetAllRequestsByID(id int64) ([]Request, error)
}
