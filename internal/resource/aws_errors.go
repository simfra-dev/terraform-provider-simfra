package resource

import "strings"

func isAWSNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "NoSuchHostedZone") ||
		strings.Contains(msg, "AWSOrganizationsNotInUseException") ||
		strings.Contains(msg, "NotFoundException") ||
		strings.Contains(msg, "404")
}
