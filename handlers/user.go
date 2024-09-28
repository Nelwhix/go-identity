package handlers

import (
	"encoding/json"
	"go-identity/pkg"
	"go-identity/pkg/requests"
	"io"
	"net/http"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		pkg.NewUnprocessableEntityResponse(w, err.Error())
		return
	}
	defer r.Body.Close()

	var request requests.SignUpRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		pkg.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	err = h.Validator.Struct(request)
	if err != nil {
		pkg.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	err = pkg.StrictPasswordValidation(request.Password)
	if err != nil {
		pkg.NewUnprocessableEntityResponse(w, err.Error())
		return
	}

	h

}
