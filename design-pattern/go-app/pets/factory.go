package pets

import "github.com/dkr2909/go-advanced-projects/design-pattern/go-app/models"

func NewPet(species string) *models.Pet {
	pet := models.Pet{
		Species:     species,
		Breed:       "",
		MinWeight:   0,
		MaxWeight:   0,
		Description: "no descriptiobn entered yet",
		LifeSpan:    0,
	}
	return &pet
}
