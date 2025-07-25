// TODO:
// Сейчас номер страницы работает как их кол-во (3, то выведутся записи 1-3)
// Количество записей как offset (если 3, то выводим с 4-й записи и далее)
// Фильтр вроде норм работает
package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
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
	fmt.Println("Raw filter string:", filterStr)
	// Разделение условий по спец-разделителю '&$' — например, name.eq.John&$age.eq.25
	// В URL нужно писать %26$
	conditions := strings.Split(filterStr, "&$")

	// Обработка каждого условия
	for i, cond := range conditions {
		fmt.Println("Cond:", cond)
		// Условие должно быть вида "field.compare.value"
		parts := strings.Split(cond, ".")
		if len(parts) != 3 {
			return "", nil, fmt.Errorf("invalid filter format: %s", cond)
		}

		field, compare, value := parts[0], parts[1], parts[2]

		var fieldType string

		// Определение типа поля: сначала ищем в modelFields
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
		case "like":
			// Подстрока
			return fmt.Sprintf("%s ILIKE ?", key), "%" + value + "%", nil

		default:
			return "", nil, fmt.Errorf("unsupported string compare: %s", compare)
		}

	case "uint":
		switch compare {
		case "eq":
			return fmt.Sprintf("%s = ?", key), value, nil

		default:
			return "", nil, fmt.Errorf("unsupported uint compare: %s", compare)
		}

	case "bool":
		switch compare {
		case "eq":
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid boolean value: %s", value)
			}
			return fmt.Sprintf("%s = ?", key), boolVal, nil
		default:
			return "", nil, fmt.Errorf("unsupported bool compare: %s", compare)
		}

	case "Time":
		if datePattern.MatchString(value) {
			t, err := time.Parse(DATE_LAYOUT, value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid date: %v", err)
			}
			if compare == "eq" {
				return fmt.Sprintf("%s::date = ?", key), t, nil
			}
			return "", nil, fmt.Errorf("only 'eq' supported for date, got: %s", compare)

		} else if timePattern.MatchString(value) {
			t, err := time.Parse(TIME_LAYOUT, value)
			if err != nil {
				return "", nil, fmt.Errorf("invalid time: %v", err)

			}
			if compare == "eq" {
				return fmt.Sprintf("%s::time = ?", key), t.Format("15:04:05"), nil
			}
			return "", nil, fmt.Errorf("only 'eq' supported for time, got: %s", compare)

		} else {
			return "", nil, fmt.Errorf("unrecognized time/date format: %s", value)
		}

	default:
		return "", nil, fmt.Errorf("unsupported type: %s", fieldType)
	}
}

// ParseOrderString - создает отдельное условие SQL сортировки
func (s *FilterBuilder) ParseOrderString(orderStr string, modelFields map[string]string) (string, error) {
	if len(orderStr) == 0 {
		return "", nil
	}

	pairs := strings.Split(orderStr, "&$")
	if len(pairs) < 1 {
		return "", fmt.Errorf("count order parameters must be greater than 0, order string: %s", orderStr)
	}

	var query string

	for i, pair := range pairs {
		kv := strings.Split(pair, ".")
		if len(kv) != 2 {
			return "", fmt.Errorf("invalid order format: %s, expected 'key.value'", pair)
		}

		field := kv[0]

		if _, ok := modelFields[field]; !ok {
			return "", fmt.Errorf("invalid order key: %s", field)
		}

		if kv[1] != "desc" && kv[1] != "asc" {
			return "", fmt.Errorf("invalid order parameter: %s", kv[1])
		}

		// Комбинируем условия
		if i+1 < len(pairs) {
			query += fmt.Sprintf("%s %s, ", field, kv[1]) // Если это не первый и не последний элемент, то добавляем запятую
		} else {
			query += fmt.Sprintf("%s %s ", field, kv[1]) // Если это первый или последний фильтр, то запятая не нужна
		}
	}

	return query, nil
}
