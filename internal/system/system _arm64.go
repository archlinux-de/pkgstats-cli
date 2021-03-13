package system

import "runtime"

func (system *System) GetCpuArchitecture() (string, error) {
	return "aarch64", nil
}
