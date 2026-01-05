package game

// TurnAction 回合動作
type TurnAction struct {
	ActionType string // "roll_dice", "make_decision", "use_item"
	Payload    interface{}
}

// TurnResult 回合結果
type TurnResult struct {
	DiceValue        int
	OldPosition      int
	NewPosition      int
	CapitalChange    int64
	EmployeeChange   int
	CircuitCompleted bool
	DecisionRequired bool
	CellType         string // 落在的格子類型
}

// CircuitBonus 繞圈獎勵
type CircuitBonus struct {
	Capital   int64
	Employees int
}

// DefaultCircuitBonus 預設繞圈獎勵
var DefaultCircuitBonus = CircuitBonus{
	Capital:   200,
	Employees: 5,
}
