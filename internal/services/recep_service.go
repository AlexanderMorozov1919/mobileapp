package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/repository"
)

type receptionService struct {
	repo repository.ReceptionRepository
}

// NewReceptionService создает новый экземпляр сервиса
func NewReceptionService(repo repository.ReceptionRepository) ReceptionService {
	return &receptionService{repo: repo}
}

// CreateReception создает новую запись на прием с валидацией
func (s *receptionService) CreateReception(reception *models.Reception) error {
	// Установка статуса по умолчанию
	reception.Status = models.StatusScheduled

	// Валидация времени приема
	if reception.Date.Before(time.Now()) {
		return errors.New("cannot create reception in the past")
	}

	// Проверка на пересечение с другими записями врача
	// 	existing, err := s.repo.GetAllByDoctorAndDate(&reception.DoctorID, &reception.Date)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	for _, e := range existing {
	// 		if isTimeOverlap(e.Date, reception.Date, 30*time.Minute) {
	// 			return errors.New("doctor already has an appointment at this time")
	// 		}
	// 	}

	return s.repo.Create(reception)
}

// UpdateReception обновляет информацию о записи
func (s *receptionService) UpdateReception(reception *models.Reception) error {
	// Проверяем существование записи
	existing, err := s.repo.GetByID(reception.ID)
	if err != nil {
		return err
	}

	// Запрещаем изменять некоторые поля после создания
	reception.DoctorID = existing.DoctorID
	reception.PatientID = existing.PatientID
	reception.Status = existing.Status

	return s.repo.Update(reception)
}

// CancelReception отменяет запись
func (s *receptionService) CancelReception(id uint, reason string) error {
	reception, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if reception.Status == models.StatusCompleted {
		return errors.New("cannot cancel completed reception")
	}

	reception.Status = models.StatusCancelled
	reception.Recommendations = reason // Используем поле рекомендаций для причины отмены

	return s.repo.Update(reception)
}

// CompleteReception отмечает прием как завершенный
func (s *receptionService) CompleteReception(id uint, diagnosis string, recommendations string) error {
	reception, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if reception.Status == models.StatusCancelled {
		return errors.New("cannot complete cancelled reception")
	}

	reception.Status = models.StatusCompleted
	reception.Diagnosis = diagnosis
	reception.Recommendations = recommendations

	return s.repo.Update(reception)
}

// MarkAsNoShow отмечает, что пациент не явился
func (s *receptionService) MarkAsNoShow(id uint) error {
	reception, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if reception.Status != models.StatusScheduled {
		return errors.New("can only mark scheduled receptions as no-show")
	}

	reception.Status = models.StatusNoShow
	return s.repo.Update(reception)
}

// GetReceptionByID возвращает запись по ID
func (s *receptionService) GetReceptionByID(id uint) (*models.Reception, error) {
	return s.repo.GetByID(id)
}

// GetDoctorReceptions возвращает записи врача с фильтрацией по дате
func (s *receptionService) GetDoctorReceptions(doctorID uint, date *time.Time) ([]models.Reception, error) {
	return s.repo.GetAllByDoctorAndDate(&doctorID, date)
}

// GetPatientReceptions возвращает все записи пациента
func (s *receptionService) GetPatientReceptions(patientID uint) ([]models.Reception, error) {
	// В реальной реализации нужно добавить метод в репозиторий
	// Здесь используем GetAllByDoctorAndDate с nil doctorID
	// Это временное решение, пока в репозитории нет специального метода
	return s.repo.GetAllByDoctorAndDate(nil, nil)
}

// GetReceptionsByStatus возвращает записи по статусу
func (s *receptionService) GetReceptionsByStatus(status models.ReceptionStatus) ([]models.Reception, error) {
	// В реальной реализации нужно добавить метод в репозиторий
	// Здесь временно возвращаем пустой результат
	return []models.Reception{}, nil
}

// isTimeOverlap проверяет пересечение временных интервалов
func isTimeOverlap(t1, t2 time.Time, duration time.Duration) bool {
	end1 := t1.Add(duration)
	end2 := t2.Add(duration)
	return t1.Before(end2) && end1.After(t2)
}

// GetReceptionsByDoctorAndDate возвращает список записей для врача на конкретную дату
func (s *receptionService) GetReceptionsByDoctorAndDate(doctorID uint, date time.Time) ([]models.Reception, error) {
	// Валидация даты (не должна быть в прошлом)
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("date cannot be in the past")
	}

	// Используем существующий метод репозитория
	receptions, err := s.repo.GetAllByDoctorAndDate(&doctorID, &date)
	if err != nil {
		return nil, fmt.Errorf("failed to get receptions: %w", err)
	}

	// Сортируем записи по времени
	sort.Slice(receptions, func(i, j int) bool {
		return receptions[i].Date.Before(receptions[j].Date)
	})

	return receptions, nil
}
