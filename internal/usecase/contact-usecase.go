package usecase

import (
	"context"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/chat-socio/backend/pkg/pointer"
	"github.com/chat-socio/backend/pkg/uuid"
)

type ContactUsecase interface {
	RequestFriend(ctx context.Context, request *presenter.RequestFriendRequest) (*presenter.RequestFriendResponse, error)
	RejectRequestFriend(ctx context.Context, request *presenter.RejectRequestFriendRequest) (*presenter.RejectRequestFriendResponse, error)
	AcceptRequestFriend(ctx context.Context, request *presenter.AcceptRequestFriendRequest) (*presenter.AcceptRequestFriendResponse, error)
	GetListRequestFriendSent(ctx context.Context, request *presenter.GetListRequestFriendSentRequest) ([]*presenter.GetListRequestFriendSentResponse, error)
	GetListContact(ctx context.Context, request *presenter.GetListContactRequest) ([]*presenter.GetListContactResponse, error)
	GetListRequestFriendReceived(ctx context.Context, request *presenter.GetListRequestFriendReceivedRequest) ([]*presenter.GetListRequestFriendReceivedResponse, error)
}

type contactUsecase struct {
	contactRepo domain.ContactRepository
	obs         *observability.Observability
}

// GetListRequestFriendSent implements ContactUsecase.
func (c *contactUsecase) GetListRequestFriendSent(ctx context.Context, request *presenter.GetListRequestFriendSentRequest) ([]*presenter.GetListRequestFriendSentResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ContactUsecase.GetListRequestFriendSent")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	requestFriends, err := c.contactRepo.GetListRequestFriendSentByUserID(ctx, request.UserID, request.Limit, request.LastID)
	if err != nil {
		logger.Error("failed to get list request friend sent", err)
		return nil, err
	}
	response := make([]*presenter.GetListRequestFriendSentResponse, 0)
	for _, requestFriend := range requestFriends {
		getListRequestFriendSentResponse := &presenter.GetListRequestFriendSentResponse{
			FriendRequestID: requestFriend.ID,
			CreatedAt:       requestFriend.CreatedAt,
			UpdatedAt:       requestFriend.UpdatedAt,
		}
		if requestFriend.ToUser != nil {
			getListRequestFriendSentResponse.TargetUser = &presenter.UserResponse{
				UserID:    requestFriend.ToUser.ID,
				FullName:  requestFriend.ToUser.FullName,
				Avatar:    requestFriend.ToUser.Avatar,
				UserType:  requestFriend.ToUser.Type,
				CreatedAt: requestFriend.ToUser.CreatedAt,
				UpdatedAt: requestFriend.ToUser.UpdatedAt,
			}
		}
		response = append(response, getListRequestFriendSentResponse)
	}
	return response, nil
}

// AcceptRequestFriend implements ContactUsecase.
func (c *contactUsecase) AcceptRequestFriend(ctx context.Context, request *presenter.AcceptRequestFriendRequest) (*presenter.AcceptRequestFriendResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ContactUsecase.AcceptRequestFriend")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	err := c.contactRepo.AcceptRequestFriend(ctx, request.RequestFriendID)
	if err != nil {
		logger.Error("failed to accept request friend", err)
		return nil, err
	}
	return &presenter.AcceptRequestFriendResponse{
		IsSuccess: true,
	}, nil
}

// GetListContact implements ContactUsecase.
func (c *contactUsecase) GetListContact(ctx context.Context, request *presenter.GetListContactRequest) ([]*presenter.GetListContactResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ContactUsecase.GetListContact")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	contacts, err := c.contactRepo.GetListContactByUserID(ctx, request.UserID, request.Limit, request.LastID)
	if err != nil {
		logger.Error("failed to get list contact", err)
		return nil, err
	}
	response := make([]*presenter.GetListContactResponse, 0)
	for _, contact := range contacts {
		getListContactResponse := &presenter.GetListContactResponse{
			ContactID: contact.ID,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		}
		if contact.Friend != nil {
			userResponse := &presenter.UserResponse{
				UserID:    contact.Friend.ID,
				FullName:  contact.Friend.FullName,
				Avatar:    contact.Friend.Avatar,
				UserType:  contact.Friend.Type,
				CreatedAt: contact.Friend.CreatedAt,
				UpdatedAt: contact.Friend.UpdatedAt,
			}
			getListContactResponse.User = userResponse
		}
		response = append(response, getListContactResponse)
	}
	return response, nil
}

// GetListRequestFriendReceived implements ContactUsecase.
func (c *contactUsecase) GetListRequestFriendReceived(ctx context.Context, request *presenter.GetListRequestFriendReceivedRequest) ([]*presenter.GetListRequestFriendReceivedResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ContactUsecase.GetListRequestFriendReceived")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	requestFriends, err := c.contactRepo.GetListRequestFriendReceivedByUserID(ctx, request.UserID, request.Limit, request.LastID)
	if err != nil {
		logger.Error("failed to get list request friend", err)
		return nil, err
	}
	response := make([]*presenter.GetListRequestFriendReceivedResponse, 0)
	for _, requestFriend := range requestFriends {
		getListRequestFriendResponse := &presenter.GetListRequestFriendReceivedResponse{
			FriendRequestID: requestFriend.ID,
			CreatedAt:       requestFriend.CreatedAt,
			UpdatedAt:       requestFriend.UpdatedAt,
		}
		if requestFriend.FromUser != nil {
			getListRequestFriendResponse.FromUser = &presenter.UserResponse{
				UserID:    requestFriend.FromUser.ID,
				FullName:  requestFriend.FromUser.FullName,
				Avatar:    requestFriend.FromUser.Avatar,
				UserType:  requestFriend.FromUser.Type,
				CreatedAt: requestFriend.FromUser.CreatedAt,
				UpdatedAt: requestFriend.FromUser.UpdatedAt,
			}
		}
		response = append(response, getListRequestFriendResponse)
	}
	return response, nil
}

// RejectRequestFriend implements ContactUsecase.
func (c *contactUsecase) RejectRequestFriend(ctx context.Context, request *presenter.RejectRequestFriendRequest) (*presenter.RejectRequestFriendResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ContactUsecase.RejectRequestFriend")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	err := c.contactRepo.UpdateRequestFriendStatus(ctx, request.RequestFriendID, "rejected")
	if err != nil {
		logger.Error("failed to reject request friend", err)
		return nil, err
	}
	return &presenter.RejectRequestFriendResponse{
		IsSuccess: true,
	}, nil
}

// RequestFriend implements ContactUsecase.
func (c *contactUsecase) RequestFriend(ctx context.Context, request *presenter.RequestFriendRequest) (*presenter.RequestFriendResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ContactUsecase.RequestFriend")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	id, err := uuid.NewID()
	if err != nil {
		logger.Error("failed to generate request friend id", err)
		return nil, err
	}
	err = c.contactRepo.CreateRequestFriend(ctx, &domain.RequestFriend{
		ID:         id,
		FromUserID: request.FromUserID,
		ToUserID:   request.TargetUserID,
		Status:     "pending",
		CreatedAt:  pointer.ToPtr(time.Now()),
		UpdatedAt:  pointer.ToPtr(time.Now()),
	})
	if err != nil {
		logger.Error("failed to request friend", err)
		return nil, err
	}
	return &presenter.RequestFriendResponse{
		IsSuccess: true,
	}, nil
}

var _ ContactUsecase = (*contactUsecase)(nil)

func NewContactUsecase(contactRepo domain.ContactRepository, obs *observability.Observability) *contactUsecase {
	return &contactUsecase{
		contactRepo: contactRepo,
		obs:         obs,
	}
}
