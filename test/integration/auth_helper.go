package integration

import (
	"fmt"
	"net/http"

	"github.com/charmingruby/telephony/internal/infra/security/token"
)

func authenticate(jwt *token.JWTService, req *http.Request, userID int) error {
	token, err := jwt.GenerateToken(userID)
	if err != nil {
		return err
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)

	req.Header.Add("Authorization", bearerToken)

	return nil
}
