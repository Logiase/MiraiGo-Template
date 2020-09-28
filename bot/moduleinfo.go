package bot

type ModuleInfo struct {
	// ID 模块的名称
	// 应全局唯一
	ID ModuleID

	New func() Module
}

func (mi ModuleInfo) String() string {
	return string(mi.ID)
}
