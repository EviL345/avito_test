package entity

import (
	"github.com/EviL345/avito_test/internal/model/dto/response"
	"github.com/google/uuid"
	"time"
)

type Pvz struct {
	Id               uuid.UUID
	City             string
	RegistrationDate time.Time
}

func (p *Pvz) ToResponse() response.Pvz {
	return response.Pvz{
		Id:               p.Id,
		City:             p.City,
		RegistrationDate: p.RegistrationDate,
	}
}
