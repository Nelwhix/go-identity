package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pquerna/otp/totp"
	"go-identity/pkg/models"
	"go-identity/pkg/requests"
	"go-identity/pkg/responses"
	"io"
	"net/http"
	"os"
	"time"
)

func (h *Handler) GenerateOTP(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewInternalServerErrorResponse(w, "User not found")
		return
	}

	if user.MfaVerifiedAt != nil {
		responses.NewBadRequest(w, "You already have MFA enabled.")
		return
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      os.Getenv("APP_NAME"),
		AccountName: user.Email,
	})
	if err != nil {
		responses.NewInternalServerErrorResponse(w, err.Error())
	}

	mfaSecret := secret.Secret()
	user.MfaSecret = &mfaSecret
	err = h.Model.UpdateUser(r.Context(), user)
	if err != nil {
		responses.NewInternalServerErrorResponse(w, err.Error())
	}

	otpUrl := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", os.Getenv("APP_NAME"), user.Email, *user.MfaSecret, os.Getenv("APP_NAME"))

	response := responses.GenerateOtp{
		OtpUrl: otpUrl,
	}

	responses.NewOKResponseWithData(w, "Otp generated successfully", response)
}

func (h *Handler) ValidateOTP(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewInternalServerErrorResponse(w, "User not found")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}
	defer r.Body.Close()

	var request requests.ValidateOtp
	err = json.Unmarshal(body, &request)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	isValid := totp.Validate(request.Code, *user.MfaSecret)
	if !isValid {
		responses.NewBadRequest(w, "Invalid Otp Code")
		return
	}

	verifiedAt := time.Now()
	user.MfaVerifiedAt = &verifiedAt
	err = h.Model.UpdateUser(r.Context(), user)
	if err != nil {
		responses.NewInternalServerErrorResponse(w, err.Error())
		return
	}

	responses.NewOKResponse(w, "MFA verified successfully")
}
