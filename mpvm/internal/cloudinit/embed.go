package cloudinit

import (
	"embed"
	"os"
)

//go:embed templates/*
var templates embed.FS

func GetCloudInit() ([]byte, error) {
	return templates.ReadFile("templates/cloud-init.yaml")
}

func GetVMSettings() ([]byte, error) {
	return templates.ReadFile("templates/vm-settings.json")
}

func WriteCloudInitTempFile() (string, func(), error) {
	content, err := GetCloudInit()
	if err != nil {
		return "", nil, err
	}

	tmpFile, err := os.CreateTemp("", "cloud-init-*.yaml")
	if err != nil {
		return "", nil, err
	}

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", nil, err
	}
	tmpFile.Close()

	cleanup := func() {
		os.Remove(tmpFile.Name())
	}

	return tmpFile.Name(), cleanup, nil
}

func WriteVMSettingsTempFile() (string, func(), error) {
	content, err := GetVMSettings()
	if err != nil {
		return "", nil, err
	}

	tmpFile, err := os.CreateTemp("", "vm-settings-*.json")
	if err != nil {
		return "", nil, err
	}

	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", nil, err
	}
	tmpFile.Close()

	cleanup := func() {
		os.Remove(tmpFile.Name())
	}

	return tmpFile.Name(), cleanup, nil
}
