package auth

import "log"

func SendEmail(email string, message string) error {
	log.Println("Sending email to", email, "with message", message)
	return nil
}
