package system

func (system *System) GetCpuArchitecture() (string, error) {
	return "riscv64", nil
}
