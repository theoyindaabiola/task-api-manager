package dao

import (
	"taskapi/models"
	"gorm.io/gorm"
	// "fmt"
)

// create an instance of the data access object
type UserDAO struct {
	DB *gorm.DB
}

// create the instance of the dao passing in the pointer
func NewUserDAO(db *gorm.DB) *UserDAO {
	// return the address of the instance of the db we are passing in
	return &UserDAO{DB: db}
}

/** 
	to create user: pointer function: these functions are created 
	under the dao package of UserDAOclass, and the parameter is called dao
	These are methods under UserDAOclass e.g: UserDAO.CreateUserDB
	a pointer to the models is needed to connect db with the model structure, 
	because they are now on different packages
**/

func (dao *UserDAO) CreateUserDB(user *models.User) error {
	// gorm needs the instance of User{} not the user struct
	if err := dao.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// to read user, get the array of the user
func (dao *UserDAO) GetUsersDB() ([]models.User, error) {
	var users []models.User
	if err := dao.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDAO) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    if err := dao.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (dao *UserDAO) GetUserByIdDB(userID string) (*models.User, error) {
    var user models.User
    if err := dao.DB.Where("id = ?", userID).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (dao *UserDAO) GetUserDB(username string) (*models.User, error) {
	var user models.User
	// store db user data in user when found
	if err := dao.DB.Where("username = ?", username).First(&user).Error; err != nil {
		// struct in go must return a value
		return nil, err
	}
	// fmt.Println("User:....", user, username)
	return &user, nil
}

func (dao *UserDAO) GetUserVerification(VerificationToken string) (*models.User, error) {
	var user models.User
	if err := dao.DB.Where("verification_token = ?", VerificationToken).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserResetToken(ResetToken string) (*models.User, error) {
	var user models.User
	if err := dao.DB.Where("reset_token = ?", ResetToken).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) Update(user *models.User) error {
    return dao.DB.Save(user).Error
}


