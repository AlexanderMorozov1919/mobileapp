package interfaces

import "time"

type Service interface {
	ParamsParserService
	FilterBuilderService
}

// ParamsParserService Сервис преобразования типов
// Парсинг строковых параметров и приведение к единому типу
type ParamsParserService interface {
	ParseDateString(dateStr string) (time.Time, error)
	ParseTimeString(timeStr string) (time.Time, error)
	ParseUintString(uintStr string) (uint64, error)

	FormatDateToString(t time.Time) string
	FormatTimeToString(t time.Time) string
}

type FilterBuilderService interface {
	ParseFilterString(filterStr string, modelFields map[string]string) (string, []interface{}, error)
}
