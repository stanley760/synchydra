package service

import (
	"synchydra/pkg/helper/sid"
	"synchydra/pkg/log"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
}

func NewService(logger *log.Logger, sid *sid.Sid) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
	}
}
