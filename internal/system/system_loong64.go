package system

func (system *System) GetCpuArchitecture() (string, error) {
	return "loong64", nil
}
