package mailer

import (
	"fmt"
	"html"
	"net/url"
	"strings"
)

func VerificationLink(clientURL, token string) string {
	return strings.TrimRight(clientURL, "/") + "/auth/verify-email?token=" + url.QueryEscape(token)
}

func ResetPasswordLink(clientURL, token string) string {
	return strings.TrimRight(clientURL, "/") + "/auth/reset-password?token=" + url.QueryEscape(token)
}

func VerificationEmail(appName, name, link string) (string, string) {
	subject := "Verify your email address"
	body := fmt.Sprintf(
		`<p>Hi %s,</p><p>Please verify your email address for %s by opening the link below. The link expires in 1 hour.</p><p><a href="%s">Verify email</a></p>`,
		html.EscapeString(name), html.EscapeString(appName), html.EscapeString(link),
	)
	return subject, body
}

func ResetPasswordEmail(appName, name, link string) (string, string) {
	subject := "Reset your password"
	body := fmt.Sprintf(
		`<p>Hi %s,</p><p>A password reset was requested for your %s account. Open the link below to choose a new password. The link expires in 1 hour.</p><p><a href="%s">Reset password</a></p><p>If you did not request this, you can ignore this email.</p>`,
		html.EscapeString(name), html.EscapeString(appName), html.EscapeString(link),
	)
	return subject, body
}
