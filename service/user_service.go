package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/namle133/Login1.git/LOGIN/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	Db *gorm.DB
}

var jwtKey = []byte("my-secrect-key")

var user string

func AddUser(s string) {
	user = s
}

func GetUser() string {
	return user
}

func hash(s string) []byte {
	bsp, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return bsp
}

func (us *UserService) SignIn(c context.Context, ui *domain.UserInit) (*domain.Claims, error) {
	var u *domain.User
	e := us.Db.First(&u, "username=? and email=?", ui.Username, ui.Email).Scan(&u).Error
	if e != nil {
		return nil, e
	}
	er := bcrypt.CompareHashAndPassword(u.Password, []byte(ui.Password))
	if er != nil {
		return nil, er
	} else {
		fmt.Println("password are equal")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &domain.Claims{
		Username: ui.Username,
		Password: string(hash(ui.Password)),
		Email:    ui.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, e
	}
	tk := &domain.Token{Username: ui.Username, TokenString: tokenString}
	failed := us.Db.Create(tk).Error
	if failed != nil {
		return nil, failed
	}
	AddUser(claims.Username)
	return claims, nil
}

func (us *UserService) CreateUser(c context.Context, u *domain.UserInit) error {
	if u == nil {
		return fmt.Errorf("user does not empty")
	}
	if u.Username == "" {
		return fmt.Errorf("Error with username")
	}
	if u.Email == "" {
		return fmt.Errorf("Error with email")
	}
	if u.Password == "" {
		return fmt.Errorf("Error with password")
	}
	uh := &domain.User{Username: u.Username, Password: hash(u.Password), Email: u.Email}
	err := us.Db.Create(uh).Error
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UserAdmin() error {
	u := domain.UserInit{Username: "admin", Password: "admin1234", Email: "admin@gmail.com"}
	uh := &domain.User{Username: u.Username, Password: hash(u.Password), Email: u.Email}
	err := us.Db.Create(uh).Error
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) CheckRowToken(c context.Context) error {
	var t *domain.Token
	rToken := us.Db.Find(&t).RowsAffected
	if rToken < 1 {
		return fmt.Errorf("Cannot create user")
	}
	return nil
}

func (us *UserService) LogOut(c context.Context) error {
	s := GetUser()
	var t *domain.Token
	err := us.Db.Where("username=?", s).Delete(&t).Error
	if err != nil {
		return err
	}
	return nil
}
