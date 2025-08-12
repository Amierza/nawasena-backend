package dto

import (
	"errors"

	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/response"
)

const (
	// ====================================== Failed ======================================
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"

	// Token
	MESSAGE_FAILED_PROSES_REQUEST    = "failed proses request"
	MESSAGE_FAILED_ACCESS_DENIED     = "failed access denied"
	MESSAGE_FAILED_TOKEN_NOT_FOUND   = "failed token not found"
	MESSAGE_FAILED_TOKEN_NOT_VALID   = "failed token not valid"
	MESSAGE_FAILED_GET_CUSTOM_CLAIMS = "failed get custom claims"
	MESSAGE_FAILED_GET_ROLE_USER     = "failed get role user"

	// Middleware
	MESSAGE_FAILED_TOKEN_DENIED_ACCESS = "failed token denied access"

	// File
	MESSAGE_FAILED_PARSE_MULTIPART_FORM = "failed to parse multipart form"
	MESSAGE_FAILED_NO_FILES_UPLOADED    = "failed no files uploaded"
	MESSAGE_FAILED_FILES_IS_EMPTY       = "failed files is empty"
	MESSAGE_FAILED_UPLOAD_FILES         = "failed upload files"

	// Admin
	MESSAGE_FAILED_CREATE_ADMIN     = "failed create admin"
	MESSAGE_FAILED_GET_LIST_ADMIN   = "failed get all admin"
	MESSAGE_FAILED_GET_DETAIL_ADMIN = "failed get detail admin"
	MESSAGE_FAILED_UPDATE_ADMIN     = "failed update admin"
	MESSAGE_FAILED_DELETE_ADMIN     = "failed delete admin"

	// Position
	MESSAGE_FAILED_CREATE_POSITION     = "failed create position"
	MESSAGE_FAILED_GET_LIST_POSITION   = "failed get all position"
	MESSAGE_FAILED_GET_DETAIL_POSITION = "failed get detail position"
	MESSAGE_FAILED_UPDATE_POSITION     = "failed update position"
	MESSAGE_FAILED_DELETE_POSITION     = "failed delete position"

	// Authentication
	MESSAGE_FAILED_LOGIN_USER    = "failed login user"
	MESSAGE_FAILED_REFRESH_TOKEN = "failed refresh token"

	// ====================================== Success ======================================
	// File
	MESSAGE_SUCCESS_UPLOAD_FILES = "success upload files"

	// Authentication
	MESSAGE_SUCCESS_LOGIN_USER    = "success login user"
	MESSAGE_SUCCESS_REFRESH_TOKEN = "success refresh token"

	// Admin
	MESSAGE_SUCCESS_CREATE_ADMIN     = "success create admin"
	MESSAGE_SUCCESS_GET_LIST_ADMIN   = "success get all admin"
	MESSAGE_SUCCESS_GET_DETAIL_ADMIN = "success get detail admin"
	MESSAGE_SUCCESS_UPDATE_ADMIN     = "success update admin"
	MESSAGE_SUCCESS_DELETE_ADMIN     = "success delete admin"

	// Position
	MESSAGE_SUCCESS_CREATE_POSITION     = "success create position"
	MESSAGE_SUCCESS_GET_LIST_POSITION   = "success get all position"
	MESSAGE_SUCCESS_GET_DETAIL_POSITION = "success get detail position"
	MESSAGE_SUCCESS_UPDATE_POSITION     = "success update position"
	MESSAGE_SUCCESS_DELETE_POSITION     = "success delete position"
)

var (
	// Token
	ErrGenerateAccessToken       = errors.New("failed to generate access token")
	ErrGenerateRefreshToken      = errors.New("failed to generate refresh token")
	ErrUnexpectedSigningMethod   = errors.New("unexpected signing method")
	ErrDecryptToken              = errors.New("failed to decrypt token")
	ErrTokenInvalid              = errors.New("token invalid")
	ErrValidateToken             = errors.New("failed to validate token")
	ErrGetAdminIDFromToken       = errors.New("failed get admin id from token")
	ErrGetAdminRoleNameFromToken = errors.New("failed get admin role name from token")

	// Middleware
	ErrDeniedAccess = errors.New("denied access")

	// Input Validation
	ErrEmptyEmail       = errors.New("email is required")
	ErrEmptyPassword    = errors.New("password is required")
	ErrEmptyName        = errors.New("name is required")
	ErrNameTooShort     = errors.New("name must be at least 3 characters")
	ErrEmptyDesc        = errors.New("description is required")
	ErrDescTooShort     = errors.New("description must be at least 5 characters")
	ErrEmptyImages      = errors.New("failed image is required")
	ErrEmptyPhoneNumber = errors.New("failed phone number is required")

	// Phone Number
	ErrFormatPhoneNumber = errors.New("failed format phone number")

	// File
	ErrNoFilesUploaded    = errors.New("failed no files uploaded")
	ErrInvalidFileType    = errors.New("only jpg/jpeg/png allowed")
	ErrSaveFile           = errors.New("failed save file")
	ErrCreateFolderAssets = errors.New("failed create folder assets")
	ErrDeleteOldImage     = errors.New("failed to delete old image")

	// Auth
	ErrInvalidEmail      = errors.New("email is required and must be in a valid format (ex: admin@example.com)")
	ErrInvalidPassword   = errors.New("password is required and must be at least 8 characters long")
	ErrIncorrectPassword = errors.New("incorrect password")

	// Admin
	ErrGetAdminByEmail           = errors.New("failed get admin by email")
	ErrAdminNotFound             = errors.New("admin not found")
	ErrEmailAlreadyExists        = errors.New("failed email already exists")
	ErrHashPassword              = errors.New("failed hash password")
	ErrCreateAdmin               = errors.New("failed create admin")
	ErrGetAllAdmin               = errors.New("failed get all admin")
	ErrGetAllAdminNoPagination   = errors.New("failed get all admin no pagination")
	ErrGetAllAdminWithPagination = errors.New("failed get all admin with pagination")
	ErrAdminAlreadyExists        = errors.New("failed admin already exists")
	ErrUpdateAdmin               = errors.New("failed update admin")
	ErrDeleteAdminByID           = errors.New("failed delete admin by id")

	// Position
	ErrGetPositionByName            = errors.New("failed get position by name")
	ErrPositionNotFound             = errors.New("position not found")
	ErrCreatePosition               = errors.New("failed create position")
	ErrGetAllPosition               = errors.New("failed get all position")
	ErrGetAllPositionNoPagination   = errors.New("failed get all position no pagination")
	ErrGetAllPositionWithPagination = errors.New("failed get all position with pagination")
	ErrPositionAlreadyExists        = errors.New("failed position already exists")
	ErrUpdatePosition               = errors.New("failed update position")
	ErrDeletePositionByID           = errors.New("failed delete position by id")
)

// Authentiation for Admin
type (
	LoginRequest struct {
		Email    string `json:"email" example:"user@example.com"`
		Password string `json:"password" example:"secret123"`
	}
	LoginResponse struct {
		AccessToken  string `json:"access_token" example:"<access_token_here>"`
		RefreshToken string `json:"refresh_token" example:"<refresh_token_here>"`
	}
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" example:"<refresh_token_here>"`
	}
	RefreshTokenResponse struct {
		AccessToken string `json:"access_token" example:"<new_access_token_here>"`
	}
)

// Admin
type (
	AdminResponse struct {
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Email       string      `json:"email"`
		Password    string      `json:"password"`
		Role        entity.Role `json:"role"`
		PhoneNumber string      `json:"phone_number"`
	}
	CreateAdminRequest struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
	}
	UpdateAdminRequest struct {
		ID          string `json:"-"`
		Name        string `json:"name,omitempty"`
		Email       string `json:"email,omitempty"`
		Password    string `json:"password,omitempty"`
		PhoneNumber string `json:"phone_number,omitempty"`
	}
	AdminPaginationResponse struct {
		response.PaginationResponse
		Data []AdminResponse `json:"data"`
	}
	AdminPaginationRepositoryResponse struct {
		response.PaginationResponse
		Admins []entity.Admin
	}
)

// Position
type (
	PositionResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	CreatePositionRequest struct {
		Name string `json:"name"`
	}
	UpdatePositionRequest struct {
		ID   string `json:"-"`
		Name string `json:"name,omitempty"`
	}
	PositionPaginationResponse struct {
		response.PaginationResponse
		Data []PositionResponse `json:"data"`
	}
	PositionPaginationRepositoryResponse struct {
		response.PaginationResponse
		Positions []entity.Position
	}
)
