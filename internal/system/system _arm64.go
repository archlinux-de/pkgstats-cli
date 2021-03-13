package system

func (system *System) GetCpuArchitecture() (string, error) {
	return "aarch64", nil
}
