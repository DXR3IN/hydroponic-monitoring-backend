package repository

import (
	"errors"
	"time"

	models "github.com/DXR3IN/auth-service-v2/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook generates UUID before inserting a new user
// This ensures UUID is always set, even if database default fails
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// BeforeUpdate hook updates the UpdatedAt timestamp
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ToDomain() *models.User {
	if u == nil {
		return nil
	}

	return &models.User{
		Name:      u.Name,
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UserRepository interface {
	Create(u *User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	EditPasswordByID(ID, newPassword string) error
	EditNameByID(ID, newName string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *userRepo) EditPasswordByID(ID, newPassword string) error {
	return r.db.Model(&User{}).Where("id = ?", ID).Update("password", newPassword).Error
}

func (r *userRepo) EditNameByID(ID, newName string) error {
	return r.db.Model(&User{}).Where("id = ?", ID).Update("name", newName).Error
}

func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return u.ToDomain(), nil
}

func (r *userRepo) FindByID(id string) (*models.User, error) {
	var u User
	if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return u.ToDomain(), nil
}
