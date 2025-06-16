package uid

import (
	"fmt"

	"github.com/google/uuid"
)

func GetUUID() string {
	id := uuid.New()
	fmt.Printf("uuid:%s, version:%s", id.String(), id.Version().String())
	return id.String()
}
