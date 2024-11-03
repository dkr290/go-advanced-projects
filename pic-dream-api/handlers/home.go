package handlers

import (
	"fmt"
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/home"
)

func (s *Handlers) HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	// account, err := s.Bun.GetAccountByUserID(user.ID)
	// if err != nil {
	// 	return err
	// }
	//
	fmt.Printf("%+v\n", user.Account)

	return home.Index().Render(r.Context(), w)
}
