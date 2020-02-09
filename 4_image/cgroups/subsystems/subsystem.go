package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var (
	//SubsystemIns = []Subsystem{
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{}, // 只要定义了就有
		&CpuSubSystem{},
	}
)
