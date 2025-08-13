package service

import (
	"context"
	"strings"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/google/uuid"
)

type (
	IMemberService interface {
		Create(ctx context.Context, req dto.CreateMemberRequest) (dto.MemberResponse, error)
		GetAll(ctx context.Context) ([]dto.MemberResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.MemberPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.MemberResponse, error)
		Update(ctx context.Context, req dto.UpdateMemberRequest) (dto.MemberResponse, error)
		Delete(ctx context.Context, id string) (dto.MemberResponse, error)
	}

	memberService struct {
		memberRepo repository.IMemberRepository
		jwt        jwt.IJWT
	}
)

func NewMemberService(memberRepo repository.IMemberRepository, jwt jwt.IJWT) *memberService {
	return &memberService{
		memberRepo: memberRepo,
		jwt:        jwt,
	}
}

func (ms *memberService) Create(ctx context.Context, req dto.CreateMemberRequest) (dto.MemberResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.MemberResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.MemberResponse{}, dto.ErrNameTooShort
	}

	// handle image request
	if req.Image == "" {
		return dto.MemberResponse{}, dto.ErrEmptyImage
	}
	if !strings.HasPrefix(req.Image, "assets/") {
		return dto.MemberResponse{}, dto.ErrFormatImage
	}

	// handle major request
	if req.Major == "" {
		return dto.MemberResponse{}, dto.ErrEmptyMajor
	}

	// handle generation request
	if req.Generation == nil {
		return dto.MemberResponse{}, dto.ErrEmptyGeneration
	}

	// handle position id request
	position, found, err := ms.memberRepo.GetPositionByPositionID(ctx, nil, req.PositionID)
	if err != nil {
		return dto.MemberResponse{}, dto.ErrGetPositionByID
	}
	if !found {
		return dto.MemberResponse{}, dto.ErrPositionNotFound
	}

	// handle double data
	_, found, _ = ms.memberRepo.GetByNameMajorGenerationAndPositionID(ctx, nil, req.Name, req.Major, req.PositionID, *req.Generation)
	if found {
		return dto.MemberResponse{}, dto.ErrMemberAlreadyExists
	}

	id := uuid.New()
	member := &entity.Member{
		ID:         id,
		Name:       req.Name,
		Image:      req.Image,
		Major:      req.Major,
		Generation: req.Generation,
		PositionID: &position.ID,
	}

	err = ms.memberRepo.Create(ctx, nil, member)
	if err != nil {
		return dto.MemberResponse{}, dto.ErrCreateMember
	}

	return dto.MemberResponse{
		ID:         member.ID.String(),
		Name:       member.Name,
		Image:      member.Image,
		Major:      member.Major,
		Generation: member.Generation,
		Position: dto.PositionResponse{
			ID:   position.ID.String(),
			Name: position.Name,
		},
	}, nil
}

func (ms *memberService) GetAll(ctx context.Context) ([]dto.MemberResponse, error) {
	members, err := ms.memberRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllMemberNoPagination
	}

	var datas []dto.MemberResponse
	for _, member := range members {
		data := dto.MemberResponse{
			ID:         member.ID.String(),
			Name:       member.Name,
			Image:      member.Image,
			Major:      member.Major,
			Generation: member.Generation,
			Position: dto.PositionResponse{
				ID:   member.Position.ID.String(),
				Name: member.Position.Name,
			},
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (ms *memberService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.MemberPaginationResponse, error) {
	dataWithPaginate, err := ms.memberRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.MemberPaginationResponse{}, dto.ErrGetAllMemberWithPagination
	}

	var datas []dto.MemberResponse
	for _, member := range dataWithPaginate.Members {
		data := dto.MemberResponse{
			ID:         member.ID.String(),
			Name:       member.Name,
			Image:      member.Image,
			Major:      member.Major,
			Generation: member.Generation,
			Position: dto.PositionResponse{
				ID:   member.Position.ID.String(),
				Name: member.Position.Name,
			},
		}

		datas = append(datas, data)
	}

	return dto.MemberPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (ms *memberService) GetDetail(ctx context.Context, id string) (dto.MemberResponse, error) {
	member, _, err := ms.memberRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.MemberResponse{}, dto.ErrMemberNotFound
	}

	return dto.MemberResponse{
		ID:         member.ID.String(),
		Name:       member.Name,
		Image:      member.Image,
		Major:      member.Major,
		Generation: member.Generation,
		Position: dto.PositionResponse{
			ID:   member.Position.ID.String(),
			Name: member.Position.Name,
		},
	}, nil
}

func (ms *memberService) Update(ctx context.Context, req dto.UpdateMemberRequest) (dto.MemberResponse, error) {
	// get member by id
	member, found, err := ms.memberRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.MemberResponse{}, dto.ErrGetMemberByID
	}
	if !found {
		return dto.MemberResponse{}, dto.ErrMemberNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != member.Name {
		if len(req.Name) < 3 {
			return dto.MemberResponse{}, dto.ErrNameTooShort
		}

		member.Name = req.Name
	}

	// handle image request
	if req.Image != "" && req.Image != member.Image {
		if !strings.HasPrefix(req.Image, "assets/") {
			return dto.MemberResponse{}, dto.ErrFormatImage
		}

		member.Image = req.Image
	}

	// handle major request
	if req.Major != "" && req.Major != member.Major {
		member.Major = req.Major
	}

	// handle generation request
	if req.Generation != nil && req.Generation != member.Generation {
		member.Generation = req.Generation
	}

	// handle position id request
	if req.PositionID != "" && req.PositionID != member.PositionID.String() {
		position, found, err := ms.memberRepo.GetPositionByPositionID(ctx, nil, req.PositionID)
		if err != nil {
			return dto.MemberResponse{}, dto.ErrGetPositionByID
		}
		if !found {
			return dto.MemberResponse{}, dto.ErrPositionNotFound
		}

		positionID, err := uuid.Parse(req.PositionID)
		if err != nil {
			return dto.MemberResponse{}, dto.ErrParseUUID
		}

		member.PositionID = &positionID
		member.Position = *position
	}

	err = ms.memberRepo.Update(ctx, nil, member)
	if err != nil {
		return dto.MemberResponse{}, dto.ErrUpdateMember
	}

	res := dto.MemberResponse{
		ID:         member.ID.String(),
		Name:       member.Name,
		Image:      member.Image,
		Major:      member.Major,
		Generation: member.Generation,
		Position: dto.PositionResponse{
			ID:   member.Position.ID.String(),
			Name: member.Position.Name,
		},
	}

	return res, nil
}

func (ms *memberService) Delete(ctx context.Context, id string) (dto.MemberResponse, error) {
	deletedMember, found, err := ms.memberRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.MemberResponse{}, dto.ErrMemberNotFound
	}
	if !found {
		return dto.MemberResponse{}, dto.ErrMemberNotFound
	}

	err = ms.memberRepo.DeleteByID(ctx, nil, id)
	if err != nil {
		return dto.MemberResponse{}, dto.ErrDeleteMemberByID
	}

	res := dto.MemberResponse{
		ID:         deletedMember.ID.String(),
		Name:       deletedMember.Name,
		Image:      deletedMember.Image,
		Major:      deletedMember.Major,
		Generation: deletedMember.Generation,
		Position: dto.PositionResponse{
			ID:   deletedMember.Position.ID.String(),
			Name: deletedMember.Position.Name,
		},
	}

	return res, nil
}
