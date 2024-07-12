package controllers

import (
	"ClubTennis/services"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BackupController struct {
	userService  *services.UserService
	resetService *services.ResetService
}

func NewBackupController(userservice *services.UserService, resetService *services.ResetService) *BackupController {
	return &BackupController{
		userService:  userservice,
		resetService: resetService,
	}
}

//-------------------------------------------------------------
// GET HANDLERS
/*
	GET .../backups/users

	returns a formatted .xlsx document containing all users and their details
*/
func (ctrl *BackupController) GetBackupSpreadsheet(c *gin.Context) {
	users, err := ctrl.userService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occured fetching the users"})
		return
	}

	filepath, err := services.UsersToSheet(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occured building the spreadsheet"})
		return
	}
	defer os.Remove(filepath)
	c.Status(http.StatusOK)
	c.FileAttachment(filepath, "users.xlsx")
}

//-------------------------------------------------------------
// POST HANDLERS
/*
	POST .../backups/users

	expects a .xlsx document formatted just like the one returned from a GET request
*/
func (ctrl *BackupController) LoadBackupSpreadsheet(c *gin.Context) {
	replace, err := strconv.ParseBool(c.PostForm("replace"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "replace not a bool, how did you manage that?"})
		return
	}
	file, err := c.FormFile("users")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file not present"})
		return
	}
	reg := regexp.MustCompile(`^.*\.xlsx$`)
	if !reg.MatchString(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file not .xlsx format"})
		return
	}
	filename := os.TempDir() + "/imported.xlsx"

	c.SaveUploadedFile(file, filename)
	defer os.Remove(filename)

	docUsers, err := services.SheetToUsers(filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "excel spreadsheet not properly formatted"})
		return
	}

	if replace { //delete users and replace them
		users, _ := ctrl.userService.FindAll()
		if users == nil {
			c.String(http.StatusInternalServerError, "couldnt get users")
			return
		}
		for _, u := range users {
			ctrl.userService.DeleteByID(u.ID)
		}
	}
	ctrl.userService.SaveStruct(docUsers...)

	c.JSON(http.StatusOK, gin.H{"message": "Users successfully updated"})
}

//-------------------------------------------------------------
// DELETE handlers
/*
	DELETE .../admin/

	deletes EVERYTHING from the database and restarts the server. Can only be executed by the root user. dangerous
*/
func (ctrl *BackupController) WipeServer(c *gin.Context) {
	id := c.GetUint("user_id")
	if id != 0 {
		c.String(http.StatusUnauthorized, "only the clubtennis email account user may do this")
		return
	}
	ctrl.resetService.DeleteEverything()
}
