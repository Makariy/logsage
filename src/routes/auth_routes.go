package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/auth"
	"main/forms"
	"main/models"
	"main/utils"
	"net/http"
)

func AddAuthRoutes(router *gin.Engine) {
	group := router.Group("/auth/")
	group.POST("/login/", handleLogin)
	group.POST("/signup/", handleSignUp)
	group.POST("/logout/", utils.LoginRequired, handleLogout)
	group.GET("/me/", utils.LoginRequired, handleMe)
}

func shouldSignUpUser(ctx *gin.Context, userForm *forms.UserForm) (*models.User, error) {
	user, err := auth.SignUpUser(ctx, userForm.Email, userForm.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "A user with this email already exists",
			})
			return nil, err
		} else {
			ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad request",
				Error:  "Cannot sign up",
			})
			return nil, err
		}
	}
	return user, nil
}

func handleLogin(ctx *gin.Context) {
	userForm, err := utils.ShouldGetForm[forms.UserForm](ctx)
	if err != nil {
		return
	}

	user, err := auth.GetUserByCredentials(userForm.Email, userForm.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
			Status: "Bad request",
			Error:  "Invalid credentials",
		})
		return
	}

	err = auth.Authorize(ctx, user)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, forms.UserResponse{
		SuccessResponse: forms.Success,
		Email:           user.Email,
		LastLogin:       user.LastLogin,
	})
}

func handleSignUp(ctx *gin.Context) {
	userForm, err := utils.ShouldGetForm[forms.UserForm](ctx)
	if err != nil {
		return
	}

	user, err := shouldSignUpUser(ctx, userForm)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusCreated, forms.UserResponse{
		SuccessResponse: forms.Success,
		Email:           user.Email,
		LastLogin:       user.LastLogin,
	})
}

func handleLogout(ctx *gin.Context) {
	err := auth.LogoutUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(http.StatusOK, forms.Success)
}

func handleMe(ctx *gin.Context) {
	user, err := utils.GetUserFromRequest(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(http.StatusOK, forms.UserResponse{
		SuccessResponse: forms.Success,
		Email:           user.Email,
		LastLogin:       user.LastLogin,
	})
}
