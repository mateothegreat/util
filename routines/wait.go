package routines

import "time"

func WaitForCondition(condition func() bool, timeout time.Duration, tick time.Duration) bool {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return false
		case <-ticker.C:
			if condition() {
				return true
			}
		}
	}
}
