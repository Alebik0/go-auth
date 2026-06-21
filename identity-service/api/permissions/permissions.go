package permissions

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func IsUser(context *gin.Context) bool {
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")
	if !slices.Contains(userRoles, "user") {
		return false
	}

	return true
}

func IsAdmin(context *gin.Context) bool {
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")
	if !slices.Contains(userRoles, "admin") {
		return false
	}

	return true
}

func GetUserID(context *gin.Context) (uint32, error) {
	userID, err := strconv.ParseInt(context.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed get user ID: %v", err)
	}

	return uint32(userID), nil
}
