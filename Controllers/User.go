package Controllers

import (
	"first-api/Models"
	"first-api/Helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetUsers(c *gin.Context) {
	var user []Models.User
	err := Models.GetAllUsers(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{ 
			"code" : http.StatusOK, 
			"message": user,
		})
	}
}

func CreateUser(c *gin.Context) {
	var user Models.User
	c.BindJSON(&user)

	Helpers.Validate = validator.New()
	errValidate := Helpers.Validate.Struct(user)

	if(errValidate != nil){
		for _, fieldErr := range errValidate.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, Helpers.FieldError(fieldErr))
			return
		}
	}

	err := Models.CreateUser(&user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func GetUserByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var user Models.User
	err := Models.GetUserByID(&user, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(c *gin.Context) {
	var user Models.User
	id := c.Params.ByName("id")
	err := Models.GetUserByID(&user, id)
	if err != nil {
		c.JSON(http.StatusNotFound, user)
	}
	c.BindJSON(&user)
	err = Models.UpdateUser(&user, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(c *gin.Context) {
	var user Models.User
	id := c.Params.ByName("id")
	err := Models.DeleteUser(&user, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{"id" + id: "is deleted"})
	}
}
