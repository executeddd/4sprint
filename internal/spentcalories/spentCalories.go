package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",") // Делим строку на слайс строк, разделитель ","

	if len(parts) != 3 { // Проверка на длину слайса, если не равен трём = вывод ошибки
		return 0, "", 0, fmt.Errorf("ошибка, слайс должен равняться трём")
	}
	steps, err := strconv.Atoi(parts[0]) // Преобразование первого элемента слайса в int
	if err != nil {                      // проверка на ошибки
		return 0, "", 0, fmt.Errorf("ошибка при преобразовании количества шагов - %v", err)
	}
	activity := parts[1]                          // вид активности и так имеет тип string, поэтому просто передаём значение
	duration, err := time.ParseDuration(parts[2]) // Преобразование третьего элемента в тип time.Duration
	if err != nil {                               // проверка на ошибки
		return 0, "", 0, fmt.Errorf("ошибка при преобразовании длительности активности - %v", err)
	}
	return steps, activity, duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	distanceInKm := (float64(steps) * lenStep) / float64(mInKm)
	return distanceInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration == 0 { // Проверка на длительность, если равна 0 = выводим в функции 0
		return 0
	}
	distanceInKm := distance(steps)          // Узнаём дистанцию в км через функцию distance
	speed := distanceInKm / duration.Hours() // Узнаём среднюю скорость путём деления дистанции на длительность тренировки в часах
	return speed
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, activity, duration, err := parseTraining(data)
	if err != nil { // Проверка на ошибки
		return fmt.Sprintf("Ошибка при получении данных: %v\n", err)
	}

	switch activity { // Проверяем вид активности
	case "Бег":
		distanceInKm := distance(steps)
		speed := meanSpeed(steps, duration)
		spentCalories := RunningSpentCalories(steps, weight, duration)
		durationInHours := duration.Hours() // Вычисляем все нужные данные по ТЗ
		return fmt.Sprintf(" Тип тренировки: %s\n Длительность: %.2f.\n Дистанция: %.2f.\n Скорость: %.2f км/ч\n Сожгли калорий: %.2f\n", activity, durationInHours, distanceInKm, speed, spentCalories)

	case "Ходьба":
		distanceInKm := distance(steps)
		speed := meanSpeed(steps, duration)
		spentCalories := WalkingSpentCalories(steps, weight, height, duration)
		durationInHours := duration.Hours() // Вычисляем все нужные данные по ТЗ
		return fmt.Sprintf(" Тип тренировки: %s\n Длительность: %.2f.\n Дистанция: %.2f.\n Скорость: %.2f км/ч\n Сожгли калорий: %.2f\n", activity, durationInHours, distanceInKm, speed, spentCalories)
	}
	return "Неизвестный тип тренировки\n"
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration) // получаем среднюю скорость
	spentCalories := ((runningCaloriesMeanSpeedMultiplier * speed) - runningCaloriesMeanSpeedShift) * weight
	return spentCalories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration) // получаем среднюю скорость
	spentCalories := ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
	return spentCalories
}
