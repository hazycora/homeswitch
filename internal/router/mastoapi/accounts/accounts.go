package accounts

import (
	"net/http"

	"git.gay/h/homeswitch/internal/config"
	account_model "git.gay/h/homeswitch/internal/models/account"
	"git.gay/h/homeswitch/internal/router/mastoapi/apicontext"
	"git.gay/h/homeswitch/internal/router/mastoapi/form"
	"github.com/gin-gonic/gin"
)

type RegisterAccountForm struct {
	Username  string `json:"username" form:"username" validate:"required,username,max=30"`
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required,min=8"`
	Agreement bool   `json:"agreement" form:"agreement" validate:"required"`
	Locale    string `json:"locale" form:"locale" validate:"required"`
	Reason    string `json:"reason" form:"reason"`
}

func RegisterAccountHandler(c *gin.Context) {
	if !config.RegistrationsEnabled {
		http.Error(c.Writer, "Registrations not enabled.", http.StatusForbidden)
	}

	var requestForm RegisterAccountForm
	err := c.Bind(&requestForm)
	if err != nil {
		http.Error(c.Writer, "Could not bind form", http.StatusInternalServerError)
		return
	}
	err = form.ValidateForm(requestForm)
	if err != nil {
		formError, ok := err.(form.FormError)
		if !ok {
			http.Error(c.Writer, "Error validating form", http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusUnprocessableEntity, formError)
		return
	}

	account := &account_model.Account{
		Username: requestForm.Username,
		Name:     &requestForm.Username,
		Email:    requestForm.Email,
	}
	err = account_model.CreateAccount(account, requestForm.Password)
	if err != nil {
		http.Error(c.Writer, "Error creating account", http.StatusInternalServerError)
	}

	c.Writer.WriteString("")
}

func VerifyCredentialsHandler(c *gin.Context) {
	account, ok := apicontext.GetAccount(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, account.ToMastoAccount(true))
}

func GetAccountHandler(c *gin.Context) {
	accountId := c.Param("id")
	account, ok := account_model.GetAccountByID(accountId)
	if !ok {
		// TODO: error should be same as Mastodon
		http.Error(c.Writer, "Record not found", http.StatusNotFound)
		return
	}
	mastoAccount := account.ToMastoAccount(false)
	c.JSON(http.StatusOK, mastoAccount)
}

func LookupAccountHandler(c *gin.Context) {
	acct := c.PostForm("acct")
	account, ok := account_model.GetAccountByUsername(acct)
	if !ok {
		// TODO: error should be same as Mastodon
		http.Error(c.Writer, "Record not found", http.StatusNotFound)
		return
	}
	mastoAccount := account.ToMastoAccount(false)
	c.JSON(http.StatusOK, mastoAccount)
}

// TODO: implement
func FeaturedTagsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, []string{})
}

// TODO: implement
func StatusesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, []string{})
}
