package apps

import (
	"net/http"
	"strings"

	app_model "git.gay/h/homeswitch/internal/models/app"
	"git.gay/h/homeswitch/internal/router/mastoapi/apicontext"
	"git.gay/h/homeswitch/internal/router/mastoapi/form"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateAppForm struct {
	ClientName  string `form:"client_name" json:"client_name" validate:"required"`
	RedirectURI string `form:"redirect_uris" json:"redirect_uris" validate:"required"`
	Scopes      string `form:"scopes" json:"scopes"`
	Website     string `form:"website" json:"website"`
}

func CreateAppHandler(c *gin.Context) {
	var requestForm CreateAppForm
	c.Bind(&requestForm)

	err := form.ValidateForm(requestForm)
	if err != nil {
		path := c.Request.URL.Path
		formError, ok := err.(form.FormError)
		if !ok {
			log.Error().Err(err).Str("path", path).Msg("Error validating form")
			http.Error(c.Writer, "Error validating form", http.StatusInternalServerError)
			return
		}
		log.Debug().Err(formError).Str("path", path).Msg("Received invalid form")
		c.JSON(http.StatusUnprocessableEntity, formError)
		return
	}

	app := &app_model.App{
		Name:        requestForm.ClientName,
		RedirectURI: requestForm.RedirectURI,
		Scopes:      strings.Split(requestForm.Scopes, " "),
		Website:     requestForm.Website,
	}

	err = app_model.CreateApp(app)
	if err != nil {
		log.Error().Err(err).Str("path", c.Request.URL.Path).Msg("Error creating app")
		http.Error(c.Writer, "Error creating app", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":            app.ID,
		"name":          app.Name,
		"website":       app.Website,
		"redirect_uri":  app.RedirectURI,
		"client_id":     app.ClientID,
		"client_secret": app.ClientSecret,
	}

	c.JSON(http.StatusOK, response)
}

func VerifyCredentialsHandler(c *gin.Context) {
	app, ok := apicontext.GetApp(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	response := map[string]interface{}{
		"name":    app.Name,
		"website": app.Website,
	}
	c.JSON(http.StatusOK, response)
}
