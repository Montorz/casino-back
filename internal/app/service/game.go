package service

import (
	"casino-back/internal/app/model"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type IGameRepository interface {
	CreateGame(userId int, game model.Game) (int, error)
	GetGames(userId int) ([]model.Game, error)
}

type GameService struct {
	gameRepository IGameRepository
}

func NewGameService(gameRepository IGameRepository) *GameService {
	return &GameService{gameRepository: gameRepository}
}

func (s *GameService) CreateGame(userId int, game model.Game) (int, error) {
	return s.gameRepository.CreateGame(userId, game)
}

func (s *GameService) GetGames(userId int) ([]model.Game, error) {
	return s.gameRepository.GetGames(userId)
}

func (s *GameService) GetGameResult(gameName string) (interface{}, error) {
	switch gameName {
	case "crash":
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
	case "wheel":
		// Массив чисел на колесе
		numbers := []int{20, 1, 3, 1, 5, 1, 3, 1, 10, 1, 3, 1, 5, 1, 5, 3, 1, 10, 1, 3, 1, 5, 1, 3, 1}

		// Генерация случайного индекса для выбора числа из массива
		randomIndex := rand.Intn(len(numbers))

		// Возвращаем случайный индекс
		return randomIndex, nil
	case "fruits":
		// Определяем доступные символы
		fruitNames := []string{"Cherry", "Lemon", "Orange", "Plum", "Bell", "Bar", "7"}
		rand.Seed(time.Now().UnixNano())

		// Сначала выбираем один символ с повышенным шансом
		mainFruit := fruitNames[rand.Intn(len(fruitNames))]

		// Генерируем результат с повышенным шансом на 3 одинаковых
		result := []string{mainFruit, mainFruit, mainFruit} // Повышенный шанс на три одинаковых
		if rand.Float64() > 0.3 {                           // В 70% случаев добавляем разнообразие
			result[1] = fruitNames[rand.Intn(len(fruitNames))]
			result[2] = fruitNames[rand.Intn(len(fruitNames))]
		}

		// Возвращаем только результат
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported game: %s", gameName)
	}
}
