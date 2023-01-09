package config

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/denisbrodbeck/machineid"
)

func GetMachineID() string {
	id, err := machineid.ProtectedID("tracetest")
	// fallback to hostname based machine id in case of error
	if err != nil {
		name, err := os.Hostname()
		if err != nil {
			return "default"
		}
		sum := md5.Sum([]byte(name))
		return hex.EncodeToString(sum[:])
	}
	return id
}
