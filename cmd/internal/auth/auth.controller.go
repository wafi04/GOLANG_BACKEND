package auth

import (
	"golang/cmd/internal/auth/dto"
	"golang/cmd/internal/auth/middleware"
	"golang/cmd/internal/pkg/response"
	"golang/cmd/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService  *AuthService
	userService  *user.UserService
}

func  NewAuthController  (authService *AuthService,userService *user.UserService)  * AuthController	{
	return &AuthController{authService: authService,userService: userService}
}


func  (C *AuthController)  Register(ctx  *gin.Context){
	var dto dto.RegisterDTO

	if  err  :=  ctx.ShouldBindBodyWithJSON(&dto);  err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error"  : err.Error(),
		})
		return
	}

	err :=  C.authService.Register(dto)


	if  err != nil {
		response.InternalServerError(ctx,"Internal Server Error")
		return
	}
		response.Created(ctx,"User registered successfully",dto)

}

func (c *AuthController) Login(ctx *gin.Context) {
	var dto dto.LoginDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		response.BadRequest(ctx,err.Error())
		return
	}

	token, err := c.authService.Login(dto)
	if err != nil {
		response.Unauthorized(ctx,err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}


func  (c  *AuthController)  GetProfile(ctx *gin.Context){
	userID := middleware.GetUserIDFromContext(ctx)
	
    
    if userID == "" {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "User ID not found",
        })
        return
    }

    // Example: Using userID to fetch user profile
    user, err := c.userService.GetProfileUser(userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve user profile",
        })
        return
    }

    // Return user profile
    ctx.JSON(http.StatusOK, user)
}