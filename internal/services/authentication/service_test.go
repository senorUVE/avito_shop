package authentication_test

import (
	"auth/internal/repository"
	"auth/internal/repository/entity"
	"auth/internal/services/authentication"
	"auth/internal/services/password"
	"auth/internal/services/tokenizer"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_Authenticate_NewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockInfoQuery := repository.NewMockInfoQuery(ctrl)
	mockPasswordService := password.NewMockService(ctrl)
	mockTokenizerService := tokenizer.NewMockService(ctrl)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery)
	mockDAO.EXPECT().NewInfoQuery(gomock.Any()).Return(mockInfoQuery)

	service := authentication.New(mockDAO, mockTokenizerService, mockPasswordService)

	ctx := context.Background()
	username := "testuser"
	password := "password123"
	hashedPassword := "hashedpassword"

	mockUserQuery.EXPECT().GetUserByUsername(username).Return(entity.User{}, errors.New("not found"))
	mockPasswordService.EXPECT().GenerateSalt(ctx).Return("salt")
	mockPasswordService.EXPECT().Hash(ctx, password, "salt").Return(hashedPassword, nil)
	mockUserQuery.EXPECT().InsertUser(gomock.Any()).Return(nil)
	mockInfoQuery.EXPECT().InsertUserInfo(gomock.Any(), 1000).Return(nil)
	mockTokenizerService.EXPECT().GeneratePair(ctx, gomock.Any()).Return("access_token", "refresh_token", nil)
	mockUserQuery.EXPECT().UpdateUser(gomock.Any()).Return(int64(1), nil)

	token, err := service.Authenticate(ctx, username, password)
	assert.NoError(t, err)
	assert.Equal(t, "access_token", token)
}

func TestService_Authenticate_ExistingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockPasswordService := password.NewMockService(ctrl)
	mockTokenizerService := tokenizer.NewMockService(ctrl)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery)

	service := authentication.New(mockDAO, mockTokenizerService, mockPasswordService)

	ctx := context.Background()
	username := "testuser"
	password := "password123"
	userID := uuid.New()
	hashedPassword := "hashedpassword"

	existingUser := entity.User{
		Id:       userID,
		Username: username,
		Password: hashedPassword,
	}

	mockUserQuery.EXPECT().GetUserByUsername(username).Return(existingUser, nil)
	mockPasswordService.EXPECT().GetSalt(ctx, hashedPassword).Return("salt", nil)
	mockPasswordService.EXPECT().Hash(ctx, password, "salt").Return(hashedPassword, nil)
	mockTokenizerService.EXPECT().GeneratePair(ctx, gomock.Any()).Return("access_token", "refresh_token", nil)
	mockUserQuery.EXPECT().UpdateUser(gomock.Any()).Return(int64(1), nil)

	token, err := service.Authenticate(ctx, username, password)
	assert.NoError(t, err)
	assert.Equal(t, "access_token", token)
}

func TestService_Authenticate_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockPasswordService := password.NewMockService(ctrl)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery)

	service := authentication.New(mockDAO, nil, mockPasswordService)

	ctx := context.Background()
	username := "testuser"
	password := "wrongpassword"
	userID := uuid.New()
	hashedPassword := "hashedpassword"

	existingUser := entity.User{
		Id:       userID,
		Username: username,
		Password: hashedPassword,
	}

	mockUserQuery.EXPECT().GetUserByUsername(username).Return(existingUser, nil)
	mockPasswordService.EXPECT().GetSalt(ctx, hashedPassword).Return("salt", nil)
	mockPasswordService.EXPECT().Hash(ctx, password, "salt").Return("wronghashedpassword", nil)

	token, err := service.Authenticate(ctx, username, password)
	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.EqualError(t, err, "invalid password")
}

func TestService_Authenticate_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := repository.NewMockDAO(ctrl)
	mockUserQuery := repository.NewMockUserQuery(ctrl)
	mockPasswordService := password.NewMockService(ctrl)
	mockTokenizerService := tokenizer.NewMockService(ctrl)

	mockDAO.EXPECT().NewUserQuery(gomock.Any()).Return(mockUserQuery)

	service := authentication.New(mockDAO, mockTokenizerService, mockPasswordService)

	ctx := context.Background()
	username := "testuser"
	password := "password123"

	mockUserQuery.EXPECT().GetUserByUsername(username).Return(entity.User{}, errors.New("database error"))
	mockPasswordService.EXPECT().Hash(ctx, gomock.Any(), gomock.Any()).AnyTimes()
	mockPasswordService.EXPECT().GenerateSalt(gomock.Any()).AnyTimes()
	mockUserQuery.EXPECT().InsertUser(gomock.Any()).Return(errors.New("database error")).AnyTimes()

	token, err := service.Authenticate(ctx, username, password)
	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Contains(t, err.Error(), "database error")
}
