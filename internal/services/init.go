package services

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
)

type Service struct {
	interfaces.ParamsParserService
	interfaces.FilterBuilderService
}

func NewService() interfaces.Service {
	parser := NewParamsParser()
	return Service{
		parser,
		NewFilterBuilder(parser),
	}
}
