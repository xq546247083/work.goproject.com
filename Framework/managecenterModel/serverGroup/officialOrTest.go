package serverGroup

type OfficialOrTest int32

const (
	// 正式服
	Con_Official OfficialOrTest = 1

	// 测试服
	Con_Test OfficialOrTest = 2
)
