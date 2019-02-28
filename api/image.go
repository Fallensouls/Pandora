package api

import (
	"fmt"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
)

func UploadAvatar(c *gin.Context) {
	var err error
	defer func() { c.Set("error", err) }()

	file, err := c.FormFile("avatar")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	content, err := file.Open()
	if err != nil {
		err = errs.New(err)
		return
	}
	defer content.Close()

	buffer := make([]byte, 512)
	if _, err = content.Read(buffer); err != nil {
		err = errs.New(err)
		return
	}

	fileType := http.DetectContentType(buffer)
	if fileType != "image/jpeg" && fileType != "image/jpg" && fileType != "image/png" {
		c.String(http.StatusBadRequest, fmt.Sprintf("%s", "avatar must be jpg or png file"))
		return
	}

	id := c.GetInt64("user_id")
	ext := path.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", strconv.FormatInt(id, 10), ext)
	err = c.SaveUploadedFile(file, fmt.Sprintf("image/avatar/%s", filename))

	if err != nil {
		err = errs.New(err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("your avatar uploaded!"))
}
