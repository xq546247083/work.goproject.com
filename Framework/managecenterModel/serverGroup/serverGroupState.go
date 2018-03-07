package serverGroup

type GroupState int32

const (
	// 正常
	Con_GroupState_Normal GroupState = 1

	// 维护
	Con_GroupState_Maintain GroupState = 2
)
