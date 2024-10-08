package system

func (s *System) GetCpuArchitecture() (string, error) {
	return AARCH64, nil
}
