package utils

import "os"

func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)

	return err == nil
}

func SaveFile(fileName string, contents []byte) error {
	err := os.WriteFile(fileName, contents, 0644)

	return err
}

func RemoveFile(fileName string) error {
	if FileExists(fileName) {
		err := os.Remove(fileName)

		return err
	}

	return nil
}
