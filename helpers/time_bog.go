package helpers

import (
	"time"
)

func GetTimeBog() map[string]interface{} {
	var respuesta map[string]interface{}

	what_time_is_it := time.Now()
	inUTC, _ := time.LoadLocation("UTC")
	inBog, _ := time.LoadLocation("America/Bogota")
	data := map[string]interface{}{
		"UNIX": what_time_is_it.Unix() * 1000,                  // Unix timestamp fixed to seconds represeted in milliseconds
		"UTC":  what_time_is_it.In(inUTC).Format(time.RFC3339), // UTC timestamp fixed to seconds
		"BOG":  what_time_is_it.In(inBog).Format(time.RFC3339), // BOG timestamp fixed to seconds
	}
	respuesta = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": data}

	return respuesta
}
