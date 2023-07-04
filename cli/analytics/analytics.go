package analytics

import (
	"github.com/denisbrodbeck/machineid"
)

var (
	mid string
)

func ClientID() string {
	return mid
}

func Init() {
	id, err := machineid.ProtectedID("tracetest")
	if err == nil {
		// only use id if available.
		mid = id
	} // ignore errors and continue with an empty ID if necessary

}
