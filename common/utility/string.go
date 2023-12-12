package utility

import (
	"encoding/json"
	"log"
	"strings"
)

// PrettyPrint will transform struct data as json string for nicer log ...
func PrettyPrint(data interface{}) string {
	JSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return string(JSON)
}

// FormatPhoneNumber default for prefix ID.
func FormatPhoneNumber(phoneNumber string) (response string) {
	prefix := phoneNumber[0:3]
	if prefix != "+62" {
		phoneNumber = strings.Replace(phoneNumber, "0", "+62", 1)
	}

	response = phoneNumber

	return
}
