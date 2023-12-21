package shared

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func CreateSendCallback(start uint64, counter *uint64, c *fiber.Ctx) func(dataToSend interface{}) {
	isSecond := false
	*counter = start
	return func(dataToSend interface{}) {
		var (
			response []byte
			e        error
		)
		if isSecond {
			_, _ = c.Write([]byte(","))
		}
		if dataToSend == nil {
			response = []byte("{}")
			_, _ = c.Write(response)
		} else if response, e = json.Marshal(dataToSend); e == nil {
			_, _ = c.Write(response)
		}
		isSecond = true
		*counter++
	}
}
