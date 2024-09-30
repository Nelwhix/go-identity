package handlers

import (
	"encoding/json"
	"go-identity/pkg"
	"go-identity/pkg/models"
	"go-identity/pkg/requests"
	"go-identity/pkg/responses"
	"go-identity/pkg/tokens"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}
	defer r.Body.Close()

	var request requests.SignUp
	err = json.Unmarshal(body, &request)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	err = pkg.StrictPasswordValidation(request.Password)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	_, err = h.Model.GetUserByEmail(r.Context(), request.Email)
	if err == nil {
		responses.NewBadRequest(w, "Email already taken")
		return
	}

	user, err := h.Model.InsertIntoUsers(r.Context(), request)
	if err != nil {
		responses.NewInternalServerErrorResponse(w, err.Error())
		return
	}

	response := responses.User{
		ID:   user.ID,
		Type: "user",
		Attributes: responses.UserAttributes{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	responses.NewCreatedResponseWithData(w, "Account created successfully.", response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}
	defer r.Body.Close()

	var request requests.Login
	err = json.Unmarshal(body, &request)
	if err != nil {
		responses.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	user, err := h.Model.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		responses.NewBadRequest(w, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		responses.NewBadRequest(w, "Email or Password is incorrect")
		return
	}

	token, err := tokens.CreateToken(h.Model, user.ID)
	if err != nil {
		responses.NewInternalServerErrorResponse(w, err.Error())
		return
	}

	response := responses.UserWithToken{
		ID:   user.ID,
		Type: "user",
		Attributes: responses.UserAttributesWithToken{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Token:     token,
		},
	}

	responses.NewOKResponseWithData(w, "Login success.", response)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewInternalServerErrorResponse(w, "User not found")
	}

	response := responses.User{
		ID:   user.ID,
		Type: "user",
		Attributes: responses.UserAttributes{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	responses.NewOKResponseWithData(w, "Get user.", response)
}
