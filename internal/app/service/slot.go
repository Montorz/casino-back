package service

import (
	"casino-back/internal/app/model"
	"math"
	"math/rand"
	"time"
)

type ISlotRepository interface {
	GetSlotData(slotName string) (model.Slot, error)
}

type SlotService struct {
	slotRepository ISlotRepository
}

func NewSlotService(slotRepository ISlotRepository) *SlotService {
	return &SlotService{slotRepository: slotRepository}
}

func (s *SlotService) GetSlotData(slotName string) (model.Slot, error) {
	return s.slotRepository.GetSlotData(slotName)
}

func (s *SlotService) GetCrashResult() (float64, error) {
	// Устанавливаем сид для генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Устанавливаем вероятность того, что выпадет 1.0
	probabilityOne := 0.2 // 20% вероятность получить 1.0

	// Генерируем случайное число от 0 до 1
	randomValue := rand.Float64()

	if randomValue < probabilityOne {
		// Возвращаем 1.0 с заданной вероятностью
		return 1.0, nil
	}

	// Параметр для распределения: чем выше lambda, тем быстрее "падает" игра
	lambda := 1.2

	// Генерируем случайный коэффициент по экспоненциальному распределению
	crashPoint := math.Exp(rand.ExpFloat64() / lambda)

	// Округляем до двух знаков после запятой
	crashPoint = math.Round(crashPoint*100) / 100

	return crashPoint, nil
}
