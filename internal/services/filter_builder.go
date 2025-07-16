package services

import (
	"errors"
	"fmt"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"strings"
	"time"
)

// FilterBuilder отвечает за построение SQL-фильтров по строкам фильтра
type FilterBuilder struct {
	paramsParser interfaces.ParamsParserService
}

func NewFilterBuilder(parser interfaces.ParamsParserService) interfaces.FilterBuilderService {
	return &FilterBuilder{paramsParser: parser}
}

// ParseFilterString — разбирает строку фильтрации вида "name.eq.John&$created_at.eq.2023-07-01"
// Возвращает SQL-условие (`WHERE ...`) и список параметров для подстановки

func (s *FilterBuilder) ParseFilterString(filterStr string, modelFields map[string]string) (string, []interface{}, error) {

	// Если строка фильтра пустая — ошибка
	if len(filterStr) == 0 {
		return "", nil, errors.New("filter parameter is empty")
	}

	var query string         // Финальная SQL-строка
	var params []interface{} // Массив параметров, которые подставятся

	// Разделение условий по спец-разделителю '&$' — например, name.eq.John&$age.eq.25
	conditions := strings.Split(filterStr, "&$")

	// Обработка каждого условия
	for i, cond := range conditions {
		// Условие должно быть вида "field.compare.value"
		parts := strings.Split(cond, ".")
		if len(parts) != 3 {
			return "", nil, fmt.Errorf("invalid filter format: %s", cond)
		}

		field, compare, value := parts[0], parts[1], parts[2]

		var fieldType string

		// Определение типа поля: сначала ищем в customFields, потом в modelFields
		if t, ok := modelFields[field]; ok {
			fieldType = t
		} else {
			return "", nil, fmt.Errorf("invalid filter key: %s", field)
		}

		// Построение SQL-условия по полю, значению и типу
		subQuery, param, err := s.buildCondition(field, value, compare, fieldType)
		if err != nil {
			return "", nil, err
		}

		// Объединяем условия через AND
		if i > 0 {
			query += " AND "
		}
		query += subQuery

		// Добавляем параметры, если они есть (например, для '?')
		if param != nil {
			params = append(params, param)
		}
	}

	return query, params, nil
}

// buildCondition — создает отдельное условие SQL по типу поля и оператору сравнения
func (s *FilterBuilder) buildCondition(key, value, compare, fieldType string) (string, interface{}, error) {
	if value == "NULL" {
		switch compare {
		case "eq":
			return fmt.Sprintf("%s IS NULL", key), nil, nil
		case "ne":
			return fmt.Sprintf("%s IS NOT NULL", key), nil, nil
		default:
			return "", nil, fmt.Errorf("unsupported compare for NULL: %s", compare)
		}
	}

	switch fieldType {
	case "string":
		switch compare {
		case "eq":
			return fmt.Sprintf("%s = ?", key), value, nil
		case "sw":
			// Начинается с... (с учетом регистра)
			return fmt.Sprintf("LOWER(%s) LIKE ?", key), strings.ToLower(value) + "%", nil
		case "like":
			// Подстрока
			return fmt.Sprintf("%s LIKE ?", key), "%" + value + "%", nil

		default:
			return "", nil, fmt.Errorf("unsupported string compare: %s", compare)
		}

	case "Time":
		if datePattern.MatchString(value) {
			t, err := time.Parse(DATE_LAYOUT, value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid date: %v", err)
			}
			if compare != "eq" {
				return "", nil, fmt.Errorf("only 'eq' supported for date, got: %s", compare)
			}
			return fmt.Sprintf("%s = ?", key), t, nil

		} else if timePattern.MatchString(value) {
			t, err := time.Parse(TIME_LAYOUT, value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid time: %v", err)

			}
			if compare != "eq" {
				return "", nil, fmt.Errorf("only 'eq' supported for time, got: %s", compare)
			}
			return fmt.Sprintf("%s = ?", key), t.Format("15:04:05"), nil

		} else {
			return "", nil, fmt.Errorf("unrecognized time/date format: %s", value)
		}

	default:
		return "", nil, fmt.Errorf("unsupported type: %s", fieldType)
	}
}
