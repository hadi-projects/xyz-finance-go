package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepo.FindByID(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get User Profile",
		"data": dto.UserProfileResponse{
			UserID: user.ID,
			Email:  user.Email,
			Consumer: func() *dto.ConsumerResponse {
				if user.Consumer != nil {
					return &dto.ConsumerResponse{
						NIK:          user.Consumer.NIK,
						FullName:     user.Consumer.FullName,
						LegalName:    user.Consumer.LegalName,
						PlaceOfBirth: user.Consumer.PlaceOfBirth,
						DateOfBirth:  user.Consumer.DateOfBirth,
						Salary:       user.Consumer.Salary,
						KTPImage:     user.Consumer.KTPImage,
						SelfieImage:  user.Consumer.SelfieImage,
					}
				}
				return nil
			}(),
		},
	})
}
