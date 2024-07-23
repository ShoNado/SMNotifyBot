package timeManage

import (
	"time"
)

func WaitForTime() {
	for {
		time.Sleep(time.Second)
		if time.Now().Minute() == 0 {
			callForNotfication(time.Now().Hour())
		}

	}
}

func callForNotfication(hour int) {

}
