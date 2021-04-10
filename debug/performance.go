package debug

import (
	"fmt"
	"time"
)

/**
 * TimeStatistics
 *
 * Usage:
 * func XXX() {
 * 	 defer utils.TimeStatistics(func (delta time.Duration) { utils.Internal().Debug("time cost is ", delta) } )()
 *   ...
 * }
 */
func TimeStatistics() func(func(time.Duration)) {
	start := time.Now()
	return func(fn func(time.Duration)) {
		fn(time.Since(start))
	}
}

/**
 * TimeStatisticsAndLog
 *
 * Usage:
 * func XXX() {
 * 	 defer utils.TimeStatisticsAndPrint("CalculateSha1", nil)()
 *   ...
 * }
 *
 * Or
 * func XXX() {
 * 	 defer utils.TimeStatisticsAndPrint("CalculateSha1", func(str string) { logger.InfoF("duration= %v", str) })()
 *   ...
 * }
 */
func TimeStatisticsAndPrint(tag string, print func(str string)) func() {
	if print == nil {
		print = func(str string) { fmt.Println(str) }
	}
	calc := TimeStatistics()
	return func() {
		calc(func(delta time.Duration) { print(fmt.Sprintf("%s = %v", tag, delta)) })
	}
}
