package main

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
)

type nextDay struct {
	Name string
	time.Time
}

func (n nextDay) String() string {
	return fmt.Sprintf("Next %v : %v", n.Name, n.Format("2-January-2006"))
}

func main() {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Do: func(index int, l string) (interface{}, error) {
			if l == "Saturday" {
				return nil, fmt.Errorf("%v is example error", l)
			}
			addDays := 7
			today := int(time.Now().Weekday())
			if index > today {
				addDays = index - today
			} else if index < today {
				addDays = 7 - today + index
			}
			return nextDay{l, time.Now().AddDate(0, 0, addDays)}, nil
		},
	}

	nextdateFunc, err := prompt.Runfunc()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	nextday, err := nextdateFunc()
	if err != nil {
		fmt.Printf("func executation error: %v\n", err)
		return
	}
	fmt.Println(nextday)
}
