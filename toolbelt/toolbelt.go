package toolbelt

import ("os"
        "time"
        "math"
        "strings"
        "strconv"
        "bufio"
        "fmt")

func TimeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    fmt.Printf("\n+++===%v took %v===+++\n\n", name, elapsed)
}

func DotCounter(dotCounter int) {
    fmt.Printf("\r            ")
    str := strings.Repeat(".", dotCounter)
    if dotCounter < 3 {
        dotCounter += 1
    } else {
        dotCounter = 0
    }
    fmt.Printf("\r%s", str)
    time.Sleep(1000 * time.Millisecond)
}

func ReadLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func Contains(s_list []string, e string) bool {
    for _, a := range s_list {
        if a == e {
            return true
        }
    }
    return false
}

func Round(x, unit float64) float64 {
        x = math.Floor(x/unit)*unit
	var rounded float64
	if x > 0 {
		rounded = float64(int64(x/unit+0.5)) * unit
	} else {
		rounded = float64(int64(x/unit-0.5)) * unit
	}
	formatted, err := strconv.ParseFloat(fmt.Sprintf("%.15f", rounded), 64)
	if err != nil {
		return rounded
	}
	return formatted
}
