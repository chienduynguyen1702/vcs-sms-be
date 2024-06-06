package services

import (
	"fmt"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type IUserService interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type UserService struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewUserService(userRepo *repositories.UserRepository, roleRepo *repositories.RoleRepository) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

func (us *UserService) CreateUser(user *dtos.CreateUserRequest, adminID string) error {
	// get admin user
	adminUser, err := us.userRepo.GetUserByID(adminID)
	if err != nil {
		return err
	}
	if adminUser == nil {
		return fmt.Errorf("admin user not found")
	}
	// get organizationID from admin user
	orgID := adminUser.OrganizationID
	if orgID == 0 {
		return fmt.Errorf("admin user does not belong to any organization")
	}
	// get role by name
	var findingRoleID uint
	var errGetRole error
	var findingRole *models.Role
	if user.IsOrganizationAdmin {
		findingRole, errGetRole = us.roleRepo.GetRoleByName("Admin")
	} else {
		findingRole, errGetRole = us.roleRepo.GetRoleByName("User")
	}
	if errGetRole != nil {
		return errGetRole
	}
	if findingRole == nil {
		return fmt.Errorf("role not found")
	}

	findingRoleID = findingRole.ID

	// set organizationID for new user
	newUser := &models.User{
		Email:               user.Email,
		Password:            user.Password,
		Username:            user.Username,
		OrganizationID:      orgID,
		Phone:               user.Phone,
		IsOrganizationAdmin: user.IsOrganizationAdmin,
		RoleID:              findingRoleID,
	}
	return us.userRepo.CreateUser(newUser)
}

func (us *UserService) GetUserByEmail(email string) (dtos.UserResponse, error) {
	user := us.userRepo.GetUserByEmail(email)
	if user == nil {
		return dtos.UserResponse{}, fmt.Errorf("User not found")
	}

	return dtos.MakeUserResponse(*user), nil
}

func (us *UserService) GetUserByID(id string) (dtos.UserResponse, error) {
	user, err := us.userRepo.GetUserByID(id)
	if user == nil {
		return dtos.UserResponse{}, fmt.Errorf("User not found")
	}
	if err != nil {
		return dtos.UserResponse{}, err
	}
	return dtos.MakeUserResponse(*user), nil
}

func (us *UserService) UpdateUser(userBodyRequest dtos.UpdateUserRequest, userID string) error {
	user, err := us.userRepo.GetUserByID(userID)
	if user == nil {
		return fmt.Errorf("User not found")
	}
	if err != nil {
		return err
	}
	// get role by name
	var findingRoleID uint
	var errGetRole error
	var findingRole *models.Role
	if user.IsOrganizationAdmin {
		findingRole, errGetRole = us.roleRepo.GetRoleByName("Admin")
	} else {
		findingRole, errGetRole = us.roleRepo.GetRoleByName("User")
	}
	if errGetRole != nil {
		return errGetRole
	}
	if findingRole == nil {
		return fmt.Errorf("role not found")
	}

	findingRoleID = findingRole.ID
	user.Email = userBodyRequest.Email
	user.Username = userBodyRequest.Username
	user.Name = userBodyRequest.Name
	user.Phone = userBodyRequest.Phone
	// user.IsOrganizationAdmin = userBodyRequest.IsOrganizationAdmin
	user.RoleID = findingRoleID
	user.Password = userBodyRequest.ConfirmPassword

	return us.userRepo.UpdateUser(user)
}

func (us *UserService) ArchiveUser(id, adminID string) error {
	user, err := us.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	user.IsArchived = true
	user.ArchivedBy = adminID
	return us.userRepo.UpdateUser(user)
}
func (us *UserService) UnarchiveUser(id, adminID string) error {
	user, err := us.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	user.IsArchived = false
	user.ArchivedBy = ""
	return us.userRepo.UpdateUser(user)
}

func (us *UserService) GetUsers(orgID string, email, username string, page, limit int) (dtos.PaginateListUserResponse, error) {
	// if email or username is empty, get all users
	if email != "" || username != "" {
		// if email or username is not empty, get users by email or username
		users, err := repositories.UserRepo.GetUsersByOrganizationIDAndSearchByEmailAndUsername(orgID, email, username)
		if err != nil {
			return dtos.PaginateListUserResponse{}, err
		}
		return dtos.MakePaginateListUserResponse(users, page, limit), nil
	}
	users, err := repositories.UserRepo.GetUsersByOrganizationID(orgID)
	if err != nil {
		return dtos.PaginateListUserResponse{}, err
	}
	return dtos.MakePaginateListUserResponse(users, page, limit), nil

}

func (us *UserService) GetUsersArchived(orgID string) (dtos.ListUserResponse, error) {
	users, err := repositories.UserRepo.GetUsersArchivedByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}
	return dtos.MakeListUserResponse(users), nil
}
