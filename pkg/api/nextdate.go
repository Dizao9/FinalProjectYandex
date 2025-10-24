package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	date = date.Truncate(24 * time.Hour)
	now = now.Truncate(24 * time.Hour)
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("передана пустая строка в правила повторения")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неправильно указана дата: %v", err)
	}

	params := strings.Split(repeat, " ")
	if len(params) < 1 {
		return "", fmt.Errorf("неверный формат правила повторения")
	}

	switch params[0] {
	case "d":
		if len(params) < 2 {
			return "", fmt.Errorf("для правила d необходимо указать количество дней")
		}
		interval, err := strconv.Atoi(params[1])
		if err != nil {
			return "", fmt.Errorf("неверный параметр для d: %v", err)
		}
		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("неверное количество дней для d, число должно быть в диапазоне от 1 до 400")
		}
		if date.Before(now) {
			diffDays := int(now.Sub(date).Hours() / 24)
			if diffDays > 0 {
				intervalsToSkip := diffDays / interval
				if intervalsToSkip > 0 {
					date = date.AddDate(0, 0, intervalsToSkip*interval)
				}
			}
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		if len(params) > 1 {
			return "", fmt.Errorf("для правила y не требуется дополнительных параметров")
		}

		if date.Year() < now.Year() {
			yrsToAdd := now.Year() - date.Year()
			date = date.AddDate(yrsToAdd-1, 0, 0)
		}

		for {
			isLeapDay := date.Month() == time.February && date.Day() == 29
			date = date.AddDate(1, 0, 0)
			if isLeapDay && date.Month() == time.February && date.Day() == 28 {
				date = date.AddDate(0, 0, 1)
			}

			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil
	case "w":
		if len(params) < 2 {
			return "", fmt.Errorf("для правила w необходимо указать хотя бы один день недели")
		}

		paramsForW := strings.Split(params[1], ",")

		parametrsInMap := make(map[int]bool, len(paramsForW))
		for _, v := range paramsForW {
			vInt, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("переданы неверные значения для параметра w: %w", err)
			}
			if vInt < 1 || vInt > 7 {
				return "", fmt.Errorf("переданы неверные значения для параметра w, дни недели должны находиться в диапозоне от 1 до 7")
			}
			parametrsInMap[vInt] = true

		}
		var tzWeekDay int
		for {
			date = date.AddDate(0, 0, 1)
			goWeekDay := date.Weekday()
			if goWeekDay == time.Sunday {
				tzWeekDay = 7
			} else {
				tzWeekDay = int(goWeekDay)
			}

			if parametrsInMap[tzWeekDay] && afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}
	case "m":
		if len(params) < 2 {
			return "", fmt.Errorf("для правила w необходимо указать хотя бы один день недели")
		}
		tempDaysForM := strings.Split(params[1], ",")
		var tempMounthForM []string
		var anyMonth bool = true
		var day [32]bool
		var month [13]bool
		if len(params) == 3 {
			anyMonth = false
			tempMounthForM = strings.Split(params[2], ",")
		}
		var useLastDay bool
		var useSecondDay bool
		for _, v := range tempDaysForM {
			digValDays, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("переданы неверные параметры дней для m: %w", err)
			}
			switch {
			case digValDays > 0 && digValDays <= 31:
				day[digValDays] = true
			case digValDays == -1:
				useLastDay = true
			case digValDays == -2:
				useSecondDay = true
			default:
				return "", fmt.Errorf("переданы неверные параметры дней для m, необходимо указать числа от 1 до 31 или -1, -2, передано :%d", digValDays)
			}
		}

		if !anyMonth {
			for _, v := range tempMounthForM {
				digValMounths, err := strconv.Atoi(v)
				if err != nil {
					return "", fmt.Errorf("переданы неверные параметры месяцев для m: %w", err)
				}
				if digValMounths < 0 {
					month[len(month)+digValMounths] = true
				} else {
					month[digValMounths] = true
				}
			}
		}

		for i := 0; i < 3650; i++ {
			date = date.AddDate(0, 0, 1)
			if afterNow(date, now) {
				if anyMonth || month[date.Month()] {
					dayOfMonth := date.Day()

					if day[dayOfMonth] {
						return date.Format(DateFormat), nil
					}

					if useLastDay {
						lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
						if dayOfMonth == lastDay {
							return date.Format(DateFormat), nil
						}
					}
					if useSecondDay {
						secondDay := time.Date(date.Year(), date.Month()+1, -1, 0, 0, 0, 0, time.UTC).Day()
						if dayOfMonth == secondDay {
							return date.Format(DateFormat), nil
						}
					}
				}
			}
		}
		return "", fmt.Errorf("возникла ошибка в процессе, цикл завершился без результата")

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", params[0])
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	if dateParam == "" {
		http.Error(w, "Отсутствует параметр date", http.StatusBadRequest)
		return
	}

	if repeatParam == "" {
		http.Error(w, "Отсутствует правило повтора", http.StatusBadRequest)
		return
	}
	var now time.Time
	var err error
	if nowParam == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}
