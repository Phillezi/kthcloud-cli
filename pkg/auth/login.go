package auth

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := s.oauth2Conf.AuthCodeURL(generateRandomState(), oauth2.AccessTypeOffline)
	s.l.Debug("LoginHandler", zap.String("redirect", authURL))
	http.Redirect(w, r, authURL, http.StatusFound)
}

func generateRandomState() string {
	const charset = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
	randomize := func(c rune) rune {
		r := rand.IntN(16)
		if c == 'x' {
			return rune(fmt.Sprintf("%x", r)[0])
		}
		return rune(fmt.Sprintf("%x", (r&0x3)|0x8)[0])
	}
	return strings.Map(randomize, charset)
}
