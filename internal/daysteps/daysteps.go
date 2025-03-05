package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

const (
	kmInMlength = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	part := strings.Split(data, ",") // Строку делим на слайс строк, разделителем является ","
	if len(part) != 2 {              // Проверка на длину слайса, если не равен 2м, то вывод ошибки
		return 0, 0, fmt.Errorf("ошибка, длина слайса должна равняться двум")
	}
	steps, err := strconv.Atoi(part[0]) // Первую часть слайса преобразуем в int
	if err != nil {                     // проверка на возможную ошибку
		return 0, 0, fmt.Errorf("ошибка при преобразовании количества шагов - %v", err)
	}
	duration, err := time.ParseDuration(part[1]) // Вторую часть слайса преобразуем в time.Duration
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при преобразовании продолжительности прогулки - %v", err)
	}
	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data) // Получаем данные из функции parsePackage
	if err != nil {                            // Проверка на ошибку, если ошибка есть: Вывод пустой строки
		log.Println("Ошибка при получении данных:", err) // Реализовал вывод ошибки через пакет "log"
		return ""
	}
	if steps <= 0 { // Если шагов 0 или <0 : Вывод пустой строки
		log.Println("Ошибка, количество шагов равно или меньше нуля:", err) // Реализовал вывод ошибки через пакет "log"
		return ""
	}
	metersLength := float64(steps) * StepLength                                     // Вычисление пройденного растояния в метрах
	kmLength := metersLength / kmInMlength                                          // перевод метров в километры (Теперь здесь участвует константа kmInMLength)
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration) // Подсчет потраченных калорий путём вызова функции WalkingSpentCalories

	return fmt.Sprintf(" Количество шагов: %d.\n Дистанция составила %.2f км.\n Вы сожгли %.2f ккал.\n", steps, kmLength, calories)
}
