package config

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/denisbrodbeck/machineid"
)

func GetMachineID() string {
	id := getMachineID()
	if len(id) >= 10 {
		return id[:10] // limit lenght to avoid issues with GA
	}

	return id
}

func getMachineID() string {
	id, err := machineid.ProtectedID("tracetest")
	if err == nil {
		return id
	}

	// fallback to hostname based machine id in case of error
	name, err := os.Hostname()
	if err != nil {
		return "default"
	}
	sum := md5.Sum([]byte(name))
	return hex.EncodeToString(sum[:])
}
