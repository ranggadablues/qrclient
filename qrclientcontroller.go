package qrclient

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ranggadablues/qrreq"

	knot "github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type QRClientController struct {
}

func getPayload(result interface{}, r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func validateBody(req *http.Request) (bytes.Buffer, error) {
	// TODO: implement signature check
	h := sha256.New()
	tee := io.TeeReader(req.Body, h)

	var b bytes.Buffer
	_, err := io.Copy(&b, tee)

	return b, err
}

// IsAuthorized is internal API
func (c *QRClientController) IsAuthorized(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := toolkit.NewResult()

	var login qrreq.LoginRequest
	body, err := validateBody(k.Request)
	if err != nil {
		log.Println(err)
		return result.SetError(qrreq.ErrorBadRequest)
	}

	err = getPayload(&login, &body)
	if err != nil {
		log.Println(err)
		return result.SetError(qrreq.ErrorBadRequest)
	}

	if !login.ValidateInput() {
		log.Println(err)
		return result.SetError(qrreq.ErrorBadRequest)
	}

	auth := _qrConfig.AuthCallback(&login)

	var token string
	if auth.Success {
		token = _loginToken.add(login.Username)

	}

	resp := qrreq.IsUserAllowedToLoginResponse{
		Status:    auth.Success,
		Mobile:    auth.Mobile,
		Email:     auth.Email,
		LastLogin: auth.LastLogin,
		Name:      auth.Name,

		Token: token,
	}
	return result.SetData(resp)
}

// ForgotPassword is internal API
func (c *QRClientController) ForgotPassword(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := toolkit.NewResult()

	var forgot qrreq.ForgotPassRequest
	body, err := validateBody(k.Request)
	if err != nil {
		log.Println(err)
		return result.SetError(qrreq.ErrorBadRequest)
	}

	err = getPayload(&forgot, &body)
	if err != nil {
		log.Println(err)
		return result.SetError(qrreq.ErrorBadRequest)
	}

	forgotResult := _qrConfig.ForgotPassCallback(&forgot)

	resp := qrreq.ForgotPassResponse{
		Username: forgot.Username,
		Email:    forgotResult.Email,
		Success:  forgotResult.Success,
	}
	return result.SetData(resp)
}

// AuthorizeUser is Public API to automatically login user that come from securelogin
func (c *QRClientController) AuthorizeUser(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	token := k.Request.URL.Query().Get("token")
	back := k.Request.URL.Query().Get("back")
	if back == "" {
		back = "https://login.clearisk.io/" // hardcoded
	}

	if token == "" {
		return result.SetError(qrreq.ErrorBadRequest)
	}

	loginname := _loginToken.get(token)
	log.Printf("qrclient: authorize user %s", loginname)
	if loginname == "" {
		http.Redirect(k.Writer, k.Request, back, http.StatusTemporaryRedirect)
		return nil
	}

	return _qrConfig.LoginCallback(loginname, k)
}
