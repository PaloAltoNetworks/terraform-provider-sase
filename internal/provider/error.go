package provider

import (
	"github.com/paloaltonetworks/sase-go/api"
)

func IsObjectNotFound(e error) bool {
	return e == api.ObjectNotFoundError
}
