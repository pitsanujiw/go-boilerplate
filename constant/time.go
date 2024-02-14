package constant

import "time"

var locThailand, _ = time.LoadLocation("Asia/Bangkok")

func GetLocationThailand() *time.Location {
	return locThailand
}
