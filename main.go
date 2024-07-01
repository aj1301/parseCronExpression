package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Command not passed correctly!! Use the run command and pass the cron string as an argument")
		os.Exit(1)
	}

	cronString := os.Args[1]
	cronMap := parseCronString(cronString)
	command := strings.Fields(cronString)[5]
	formattedOutput := formatCronOutput(cronMap, command)

	fmt.Print(formattedOutput)

}

func parseCronString(cronString string) map[string][]int {
	fields := strings.Fields(cronString)
	if len(fields) != 6 {
		fmt.Println("Invalid cron string format")
		os.Exit(1)
	}

	minuteField := fields[0]
	hourField := fields[1]
	dayOfMonthField := fields[2]
	monthField := fields[3]
	dayOfWeekField := fields[4]
	command := fields[5]

	minutes := expandField(minuteField, 0, 59)
	hours := expandField(hourField, 0, 23)
	daysOfMonth := expandField(dayOfMonthField, 1, 31)
	months := expandField(monthField, 1, 12)
	daysOfWeek := expandField(dayOfWeekField, 0, 6)

	return map[string][]int{
		"minute":       minutes,
		"hour":         hours,
		"day of month": daysOfMonth,
		"month":        months,
		"day of week":  daysOfWeek,
		"command":      []int{len(command)}, // Special handling to keep the command field
	}
}

func expandField(field string, min, max int) []int {
	var result []int

	if field == "*" {
		for i := min; i <= max; i++ {
			result = append(result, i)
		}
		return result
	}

	parts := strings.Split(field, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			start, _ := strconv.Atoi(rangeParts[0])
			end, _ := strconv.Atoi(rangeParts[1])
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else if strings.Contains(part, "/") {
			stepParts := strings.Split(part, "/")
			step, _ := strconv.Atoi(stepParts[1])
			for i := min; i <= max; i += step {
				result = append(result, i)
			}
		} else {
			val, _ := strconv.Atoi(part)
			result = append(result, val)
		}
	}

	return result
}

func formatCronOutput(cronMap map[string][]int, command string) string {
	var output strings.Builder

	for key, times := range cronMap {
		if key == "command" {
			output.WriteString(fmt.Sprintf("%-14s%s\n", key, command))
		} else {
			timesStr := make([]string, len(times))
			for i, time := range times {
				timesStr[i] = strconv.Itoa(time)
			}
			output.WriteString(fmt.Sprintf("%-14s%s\n", key, strings.Join(timesStr, " ")))
		}
	}

	return output.String()
}
