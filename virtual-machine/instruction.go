package virtualmachine

type RawInstruction struct {
	Line       uint16
	Expression string
}

type Instruction struct {
	Line  uint16
	Codop int16
}
