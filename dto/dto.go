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

	// Authentication
	MESSAGE_FAILED_LOGIN_USER    = "failed login user"
	MESSAGE_FAILED_REFRESH_TOKEN = "failed refresh token"

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

	// Member
	MESSAGE_FAILED_CREATE_MEMBER     = "failed create member"
	MESSAGE_FAILED_GET_LIST_MEMBER   = "failed get all member"
	MESSAGE_FAILED_GET_DETAIL_MEMBER = "failed get detail member"
	MESSAGE_FAILED_UPDATE_MEMBER     = "failed update member"
	MESSAGE_FAILED_DELETE_MEMBER     = "failed delete member"

	// Achievement Category
	MESSAGE_FAILED_CREATE_ACHIEVEMENT_CATEGORY     = "failed create achievement category"
	MESSAGE_FAILED_GET_LIST_ACHIEVEMENT_CATEGORY   = "failed get all achievement category"
	MESSAGE_FAILED_GET_DETAIL_ACHIEVEMENT_CATEGORY = "failed get detail achievement category"
	MESSAGE_FAILED_UPDATE_ACHIEVEMENT_CATEGORY     = "failed update achievement category"
	MESSAGE_FAILED_DELETE_ACHIEVEMENT_CATEGORY     = "failed delete achievement category"

	// Achievement
	MESSAGE_FAILED_CREATE_ACHIEVEMENT     = "failed create achievement"
	MESSAGE_FAILED_GET_LIST_ACHIEVEMENT   = "failed get all achievement"
	MESSAGE_FAILED_GET_DETAIL_ACHIEVEMENT = "failed get detail achievement"
	MESSAGE_FAILED_UPDATE_ACHIEVEMENT     = "failed update achievement"
	MESSAGE_FAILED_DELETE_ACHIEVEMENT     = "failed delete achievement"

	// Ship
	MESSAGE_FAILED_CREATE_SHIP     = "failed create ship"
	MESSAGE_FAILED_GET_LIST_SHIP   = "failed get all ship"
	MESSAGE_FAILED_GET_DETAIL_SHIP = "failed get detail ship"
	MESSAGE_FAILED_UPDATE_SHIP     = "failed update ship"
	MESSAGE_FAILED_DELETE_SHIP     = "failed delete ship"

	// Competition
	MESSAGE_FAILED_CREATE_COMPETITION     = "failed create competition"
	MESSAGE_FAILED_GET_LIST_COMPETITION   = "failed get all competition"
	MESSAGE_FAILED_GET_DETAIL_COMPETITION = "failed get detail competition"
	MESSAGE_FAILED_UPDATE_COMPETITION     = "failed update competition"
	MESSAGE_FAILED_DELETE_COMPETITION     = "failed delete competition"

	// News Category
	MESSAGE_FAILED_CREATE_NEWS_CATEGORY     = "failed create news category"
	MESSAGE_FAILED_GET_LIST_NEWS_CATEGORY   = "failed get all news category"
	MESSAGE_FAILED_GET_DETAIL_NEWS_CATEGORY = "failed get detail news category"
	MESSAGE_FAILED_UPDATE_NEWS_CATEGORY     = "failed update news category"
	MESSAGE_FAILED_DELETE_NEWS_CATEGORY     = "failed delete news category"

	// News
	MESSAGE_FAILED_CREATE_NEWS     = "failed create news"
	MESSAGE_FAILED_GET_LIST_NEWS   = "failed get all news"
	MESSAGE_FAILED_GET_DETAIL_NEWS = "failed get detail news"
	MESSAGE_FAILED_UPDATE_NEWS     = "failed update news"
	MESSAGE_FAILED_DELETE_NEWS     = "failed delete news"

	// Partner
	MESSAGE_FAILED_CREATE_PARTNER     = "failed create partner"
	MESSAGE_FAILED_GET_LIST_PARTNER   = "failed get all partner"
	MESSAGE_FAILED_GET_DETAIL_PARTNER = "failed get detail partner"
	MESSAGE_FAILED_UPDATE_PARTNER     = "failed update partner"
	MESSAGE_FAILED_DELETE_PARTNER     = "failed delete partner"

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

	// Member
	MESSAGE_SUCCESS_CREATE_MEMBER     = "success create member"
	MESSAGE_SUCCESS_GET_LIST_MEMBER   = "success get all member"
	MESSAGE_SUCCESS_GET_DETAIL_MEMBER = "success get detail member"
	MESSAGE_SUCCESS_UPDATE_MEMBER     = "success update member"
	MESSAGE_SUCCESS_DELETE_MEMBER     = "success delete member"

	// Achievement Category
	MESSAGE_SUCCESS_CREATE_ACHIEVEMENT_CATEGORY     = "success create achievement category"
	MESSAGE_SUCCESS_GET_LIST_ACHIEVEMENT_CATEGORY   = "success get all achievement category"
	MESSAGE_SUCCESS_GET_DETAIL_ACHIEVEMENT_CATEGORY = "success get detail achievement category"
	MESSAGE_SUCCESS_UPDATE_ACHIEVEMENT_CATEGORY     = "success update achievement category"
	MESSAGE_SUCCESS_DELETE_ACHIEVEMENT_CATEGORY     = "success delete achievement category"

	// Achievement
	MESSAGE_SUCCESS_CREATE_ACHIEVEMENT     = "success create achievement"
	MESSAGE_SUCCESS_GET_LIST_ACHIEVEMENT   = "success get all achievement"
	MESSAGE_SUCCESS_GET_DETAIL_ACHIEVEMENT = "success get detail achievement"
	MESSAGE_SUCCESS_UPDATE_ACHIEVEMENT     = "success update achievement"
	MESSAGE_SUCCESS_DELETE_ACHIEVEMENT     = "success delete achievement"

	// Ship
	MESSAGE_SUCCESS_CREATE_SHIP     = "success create ship"
	MESSAGE_SUCCESS_GET_LIST_SHIP   = "success get all ship"
	MESSAGE_SUCCESS_GET_DETAIL_SHIP = "success get detail ship"
	MESSAGE_SUCCESS_UPDATE_SHIP     = "success update ship"
	MESSAGE_SUCCESS_DELETE_SHIP     = "success delete ship"

	// Competition
	MESSAGE_SUCCESS_CREATE_COMPETITION     = "success create competition"
	MESSAGE_SUCCESS_GET_LIST_COMPETITION   = "success get all competition"
	MESSAGE_SUCCESS_GET_DETAIL_COMPETITION = "success get detail competition"
	MESSAGE_SUCCESS_UPDATE_COMPETITION     = "success update competition"
	MESSAGE_SUCCESS_DELETE_COMPETITION     = "success delete competition"

	// News Category
	MESSAGE_SUCCESS_CREATE_NEWS_CATEGORY     = "success create news category"
	MESSAGE_SUCCESS_GET_LIST_NEWS_CATEGORY   = "success get all news category"
	MESSAGE_SUCCESS_GET_DETAIL_NEWS_CATEGORY = "success get detail news category"
	MESSAGE_SUCCESS_UPDATE_NEWS_CATEGORY     = "success update news category"
	MESSAGE_SUCCESS_DELETE_NEWS_CATEGORY     = "success delete news category"

	// News
	MESSAGE_SUCCESS_CREATE_NEWS     = "success create news"
	MESSAGE_SUCCESS_GET_LIST_NEWS   = "success get all news"
	MESSAGE_SUCCESS_GET_DETAIL_NEWS = "success get detail news"
	MESSAGE_SUCCESS_UPDATE_NEWS     = "success update news"
	MESSAGE_SUCCESS_DELETE_NEWS     = "success delete news"

	// Partner
	MESSAGE_SUCCESS_CREATE_PARTNER     = "success create partner"
	MESSAGE_SUCCESS_GET_LIST_PARTNER   = "success get all partner"
	MESSAGE_SUCCESS_GET_DETAIL_PARTNER = "success get detail partner"
	MESSAGE_SUCCESS_UPDATE_PARTNER     = "success update partner"
	MESSAGE_SUCCESS_DELETE_PARTNER     = "success delete partner"
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

	// Parse
	ErrParseUUID                 = errors.New("failed parse to uuid format")
	ErrParseLimit                = errors.New("failed parse limit to int")
	ErrParseTimeFromStringToTime = errors.New("failed parse time format from string to time.Time")
	ErrParseTimeFromTimeToString = errors.New("failed parse time format from time.Time to string")

	// Middleware
	ErrDeniedAccess = errors.New("denied access")

	// Input Validation
	ErrEmptyEmail          = errors.New("email is required")
	ErrEmptyPassword       = errors.New("password is required")
	ErrEmptyName           = errors.New("name is required")
	ErrNameTooShort        = errors.New("name must be at least 3 characters")
	ErrEmptyDesc           = errors.New("description is required")
	ErrDescTooShort        = errors.New("description must be at least 5 characters")
	ErrEmptyImage          = errors.New("failed image is required")
	ErrFormatImage         = errors.New("format image must be has prefix assets/")
	ErrEmptyPhoneNumber    = errors.New("failed phone number is required")
	ErrEmptyMajor          = errors.New("failed major is required")
	ErrEmptyGeneration     = errors.New("failed generation is required")
	ErrTypeGeneration      = errors.New("failed generation is must be int")
	ErrEmptyYear           = errors.New("failed year is required")
	ErrEmptyDescription    = errors.New("failed description is required")
	ErrDescriptionTooShort = errors.New("description must be at least 5 characters")
	ErrEmptyImages         = errors.New("failed images is required")
	ErrEmptyDate           = errors.New("failed date is required")
	ErrEmptyLocation       = errors.New("failed location is required")
	ErrLocationTooShort    = errors.New("location must be at least 5 characters")
	ErrEmptyStatus         = errors.New("failed status is required")
	ErrEmptyNewsCategory   = errors.New("failed news category is required")
	ErrEmptyTeam           = errors.New("failed team is required")
	ErrEmptyTags           = errors.New("failed tags is required")

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
	ErrGetPositionByID              = errors.New("failed get position by id")
	ErrPositionNotFound             = errors.New("position not found")
	ErrCreatePosition               = errors.New("failed create position")
	ErrGetAllPosition               = errors.New("failed get all position")
	ErrGetAllPositionNoPagination   = errors.New("failed get all position no pagination")
	ErrGetAllPositionWithPagination = errors.New("failed get all position with pagination")
	ErrPositionAlreadyExists        = errors.New("failed position already exists")
	ErrUpdatePosition               = errors.New("failed update position")
	ErrDeletePositionByID           = errors.New("failed delete position by id")

	// Member
	ErrGetMemberByID              = errors.New("failed get member by id")
	ErrGetMemberByName            = errors.New("failed get member by name")
	ErrMemberNotFound             = errors.New("member not found")
	ErrCreateMember               = errors.New("failed create member")
	ErrGetAllMember               = errors.New("failed get all member")
	ErrGetAllMemberNoPagination   = errors.New("failed get all member no pagination")
	ErrGetAllMemberWithPagination = errors.New("failed get all member with pagination")
	ErrMemberAlreadyExists        = errors.New("failed member already exists")
	ErrUpdateMember               = errors.New("failed update member")
	ErrDeleteMemberByID           = errors.New("failed delete member by id")

	// Achievement category
	ErrGetAchievementCategoryByName     = errors.New("failed get achievement category by name")
	ErrGetAchievementCategoryByID       = errors.New("failed get achievement category by id")
	ErrAchievementCategoryNotFound      = errors.New("achievement category not found")
	ErrCreateAchievementCategory        = errors.New("failed create achievement category")
	ErrGetAllAchievementCategory        = errors.New("failed get all achievement category")
	ErrAchievementCategoryAlreadyExists = errors.New("failed achievement category already exists")
	ErrUpdateAchievementCategory        = errors.New("failed update achievement category")
	ErrDeleteAchievementCategoryByID    = errors.New("failed delete achievement category by id")

	// Achievement
	ErrGetAchievementByID                    = errors.New("failed get achievement by id")
	ErrGetAchievementImages                  = errors.New("failed get achievement images")
	ErrAchievementNotFound                   = errors.New("achievement not found")
	ErrGetAllFeaturedAchievement             = errors.New("failed get all featured achievement")
	ErrCreateAchievement                     = errors.New("failed create achievement")
	ErrCreateAchievementImage                = errors.New("failed create achievement image")
	ErrGetAllAchievement                     = errors.New("failed get all achievement")
	ErrGetAllAchievementNoPagination         = errors.New("failed get all achievement no pagination")
	ErrGetAllAchievementWithPagination       = errors.New("failed get all achievement with pagination")
	ErrAchievementAlreadyExists              = errors.New("failed achievement already exists")
	ErrUpdateAchievement                     = errors.New("failed update achievement")
	ErrDeleteAchievementByID                 = errors.New("failed delete achievement by id")
	ErrDeleteAchievementImageByAchievementID = errors.New("failed delete achievement image by achievement id")

	// Ship
	ErrGetShipByID              = errors.New("failed get ship by id")
	ErrGetShipImages            = errors.New("failed get ship images")
	ErrShipNotFound             = errors.New("ship not found")
	ErrCreateShip               = errors.New("failed create ship")
	ErrCreateShipImage          = errors.New("failed create ship image")
	ErrGetAllShip               = errors.New("failed get all ship")
	ErrGetAllShipNoPagination   = errors.New("failed get all ship no pagination")
	ErrGetAllShipWithPagination = errors.New("failed get all ship with pagination")
	ErrShipAlreadyExists        = errors.New("failed ship already exists")
	ErrUpdateShip               = errors.New("failed update ship")
	ErrDeleteShipByID           = errors.New("failed delete ship by id")
	ErrDeleteShipImageByShipID  = errors.New("failed delete ship image by ship id")

	// Competition
	ErrGetCompetitionByID                    = errors.New("failed get competition by id")
	ErrGetCompetitionImages                  = errors.New("failed get competition images")
	ErrCompetitionNotFound                   = errors.New("competition not found")
	ErrCreateCompetition                     = errors.New("failed create competition")
	ErrCreateCompetitionImage                = errors.New("failed create competition image")
	ErrGetAllCompetition                     = errors.New("failed get all competition")
	ErrGetAllCompetitionNoPagination         = errors.New("failed get all competition no pagination")
	ErrGetAllCompetitionWithPagination       = errors.New("failed get all competition with pagination")
	ErrCompetitionAlreadyExists              = errors.New("failed competition already exists")
	ErrUpdateCompetition                     = errors.New("failed update competition")
	ErrDeleteCompetitionByID                 = errors.New("failed delete competition by id")
	ErrDeleteCompetitionImageByCompetitionID = errors.New("failed delete competition image by ship id")

	// News category
	ErrGetNewsCategoryByName     = errors.New("failed get news category by name")
	ErrGetNewsCategoryByID       = errors.New("failed get news category by id")
	ErrNewsCategoryNotFound      = errors.New("news category not found")
	ErrCreateNewsCategory        = errors.New("failed create news category")
	ErrGetAllNewsCategory        = errors.New("failed get all news category")
	ErrNewsCategoryAlreadyExists = errors.New("failed news category already exists")
	ErrUpdateNewsCategory        = errors.New("failed update news category")
	ErrDeleteNewsCategoryByID    = errors.New("failed delete news category by id")

	// News
	ErrGetNewsByID              = errors.New("failed get news by id")
	ErrGetNewsImages            = errors.New("failed get news images")
	ErrNewsNotFound             = errors.New("news not found")
	ErrIncrementViews           = errors.New("failed increment views")
	ErrCreateNews               = errors.New("failed create news")
	ErrCreateNewsImage          = errors.New("failed create news image")
	ErrGetAllNews               = errors.New("failed get all news")
	ErrGetAllFeaturedNews       = errors.New("failed get all featured news")
	ErrInvalidStatus            = errors.New("failed invalid status")
	ErrGetAllNewsNoPagination   = errors.New("failed get all news no pagination")
	ErrGetAllNewsWithPagination = errors.New("failed get all news with pagination")
	ErrNewsAlreadyExists        = errors.New("failed news already exists")
	ErrUpdateNews               = errors.New("failed update news")
	ErrDeleteNewsByID           = errors.New("failed delete news by id")
	ErrDeleteNewsImageByNewsID  = errors.New("failed delete news image by news id")

	// Partner
	ErrGetPartnerByID              = errors.New("failed get partner by id")
	ErrGetPartnerImage             = errors.New("failed get partner image")
	ErrPartnerNotFound             = errors.New("partner not found")
	ErrCreatePartner               = errors.New("failed create partner")
	ErrGetAllPartner               = errors.New("failed get all partner")
	ErrGetAllPartnerNoPagination   = errors.New("failed get all partner no pagination")
	ErrGetAllPartnerWithPagination = errors.New("failed get all partner with pagination")
	ErrPartnerAlreadyExists        = errors.New("failed partner already exists")
	ErrUpdatePartner               = errors.New("failed update partner")
	ErrDeletePartnerByID           = errors.New("failed delete partner by id")
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

// Member
type (
	MemberResponse struct {
		ID         string           `json:"id"`
		Name       string           `json:"name"`
		Image      string           `json:"image"`
		Major      string           `json:"major"`
		Generation *int             `json:"generation"`
		Position   PositionResponse `json:"position"`
	}
	CreateMemberRequest struct {
		Name       string `json:"name"`
		Image      string `json:"image"`
		Major      string `json:"major"`
		Generation *int   `json:"generation"`
		PositionID string `json:"position_id"`
	}
	UpdateMemberRequest struct {
		ID         string `json:"-"`
		Name       string `json:"name,omitempty"`
		Image      string `json:"image,omitempty"`
		Major      string `json:"major,omitempty"`
		Generation *int   `json:"generation,omitempty"`
		PositionID string `json:"position_id,omitempty"`
	}
	MemberPaginationResponse struct {
		response.PaginationResponse
		Data []MemberResponse `json:"data"`
	}
	MemberPaginationRepositoryResponse struct {
		response.PaginationResponse
		Members []entity.Member
	}
)

// AchievementCategory
type (
	AchievementCategoryResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	CreateAchievementCategoryRequest struct {
		Name string `json:"name"`
	}
	UpdateAchievementCategoryRequest struct {
		ID   string `json:"-"`
		Name string `json:"name,omitempty"`
	}
)

// Achievement
type (
	AchievementImageResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	AchievementResponse struct {
		ID          string                      `json:"id"`
		Name        string                      `json:"name"`
		Year        int                         `json:"year"`
		Description string                      `json:"description"`
		Location    string                      `json:"location"`
		Rank        string                      `json:"rank"`
		Competition string                      `json:"competition"`
		Team        []string                    `json:"team"`
		Impact      string                      `json:"impact"`
		VideoURL    string                      `json:"video_url"`
		Featured    bool                        `json:"featured"`
		Tags        []string                    `json:"tags"`
		Images      []AchievementImageResponse  `json:"images"`
		Category    AchievementCategoryResponse `json:"category"`
	}
	CreateAchievementRequest struct {
		Name        string   `json:"name" binding:"required"`
		Year        int      `json:"year" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Location    string   `json:"location" binding:"required"`
		Rank        string   `json:"rank" binding:"required"`
		Competition string   `json:"competition" binding:"required"`
		Team        []string `json:"team" binding:"required"`
		Impact      string   `json:"impact" binding:"required"`
		VideoURL    string   `json:"video_url"`
		Featured    bool     `json:"featured"`
		Tags        []string `json:"tags" binding:"required"`
		Images      []string `json:"images" binding:"required"`
		CategoryID  string   `json:"category_id" binding:"required"`
	}
	UpdateAchievementRequest struct {
		ID          string   `json:"-"`
		Name        string   `json:"name,omitempty"`
		Year        *int     `json:"year,omitempty"`
		Description string   `json:"description,omitempty"`
		Location    string   `json:"location,omitempty"`
		Rank        string   `json:"rank,omitempty"`
		Competition string   `json:"competition,omitempty"`
		Team        []string `json:"team,omitempty"`
		Impact      string   `json:"impact,omitempty"`
		VideoURL    string   `json:"video_url,omitempty"`
		Featured    bool     `json:"featured,omitempty"`
		Tags        []string `json:"tags,omitempty"`
		Images      []string `json:"images,omitempty"`
		CategoryID  string   `json:"category_id,omitempty"`
	}
	AchievementPaginationResponse struct {
		response.PaginationResponse
		Data []AchievementResponse `json:"data"`
	}
	AchievementPaginationRepositoryResponse struct {
		response.PaginationResponse
		Achievements []entity.Achievement
	}
)

// Ship
type (
	ShipImageResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	ShipResponse struct {
		ID          string              `json:"id"`
		Name        string              `json:"name"`
		Description string              `json:"description"`
		Images      []ShipImageResponse `json:"images"`
	}
	CreateShipRequest struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
	}
	UpdateShipRequest struct {
		ID          string   `json:"-"`
		Name        string   `json:"name,omitempty"`
		Description string   `json:"description,omitempty"`
		Images      []string `json:"images,omitempty"`
	}
	ShipPaginationResponse struct {
		response.PaginationResponse
		Data []ShipResponse `json:"data"`
	}
	ShipPaginationRepositoryResponse struct {
		response.PaginationResponse
		Ships []entity.Ship
	}
)

// Competition
type (
	CompetitionImageResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	CompetitionResponse struct {
		ID          string                     `json:"id"`
		Name        string                     `json:"name"`
		Date        string                     `json:"date"`
		Description string                     `json:"description"`
		Images      []CompetitionImageResponse `json:"images"`
	}
	CreateCompetitionRequest struct {
		Name        string   `json:"name"`
		Date        string   `json:"date"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
	}
	UpdateCompetitionRequest struct {
		ID          string   `json:"-"`
		Name        string   `json:"name,omitempty"`
		Date        string   `json:"date,omitempty"`
		Description string   `json:"description,omitempty"`
		Images      []string `json:"images,omitempty"`
	}
	CompetitionPaginationResponse struct {
		response.PaginationResponse
		Data []CompetitionResponse `json:"data"`
	}
	CompetitionPaginationRepositoryResponse struct {
		response.PaginationResponse
		Competitions []entity.Competition
	}
)

// NewsCategory
type (
	NewsCategoryResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	CreateNewsCategoryRequest struct {
		Name string `json:"name"`
	}
	UpdateNewsCategoryRequest struct {
		ID   string `json:"-"`
		Name string `json:"name,omitempty"`
	}
)

// News
type (
	NewsImageResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	NewsResponse struct {
		ID          string               `json:"id"`
		Name        string               `json:"name"`
		Description string               `json:"description"`
		PublishedAt string               `json:"published_at"`
		Location    string               `json:"location"`
		Status      string               `json:"status"` // Completed, Ongoing, Upcoming
		Views       int                  `json:"views"`
		Featured    bool                 `json:"featured"`
		Category    NewsCategoryResponse `json:"category"`
		Images      []NewsImageResponse  `json:"images"`
	}
	CreateNewsRequest struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Location    string   `json:"location"`
		Status      string   `json:"status"` // Completed, Ongoing, Upcoming
		Featured    bool     `json:"featured"`
		CategoryID  string   `json:"category_id"`
		Images      []string `json:"images"`
	}
	UpdateNewsRequest struct {
		ID          string   `json:"-"`
		Name        string   `json:"name,omitempty"`
		Description string   `json:"description,omitempty"`
		Location    string   `json:"location,omitempty"`
		Status      string   `json:"status,omitempty"` // Completed, Ongoing, Upcoming
		Featured    bool     `json:"featured,omitempty"`
		CategoryID  string   `json:"category_id,omitempty"`
		Images      []string `json:"images,omitempty"`
	}
	NewsPaginationResponse struct {
		response.PaginationResponse
		Data []NewsResponse `json:"data"`
	}
	NewsPaginationRepositoryResponse struct {
		response.PaginationResponse
		Newss []entity.News
	}
)

// Partner
type (
	PartnerResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	CreatePartnerRequest struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	UpdatePartnerRequest struct {
		ID    string `json:"-"`
		Name  string `json:"name,omitempty"`
		Image string `json:"image,omitempty"`
	}
	PartnerPaginationResponse struct {
		response.PaginationResponse
		Data []PartnerResponse `json:"data"`
	}
	PartnerPaginationRepositoryResponse struct {
		response.PaginationResponse
		Partners []entity.Partner
	}
)
