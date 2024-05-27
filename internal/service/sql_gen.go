package service

import "github.com/crossevol/sqlcc/internal/models"

type SqlGenService interface {
	Gen(config models.Config)
	GenMapper(config models.Config)
}
