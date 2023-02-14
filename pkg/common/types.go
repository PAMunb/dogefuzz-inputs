package common

type Block struct {
	PC               string                 `json:"pc"`
	Range            BlockRange             `json:"range"`
	Predecessors     []string               `json:"predecessors"`
	Successors       []string               `json:"successors"`
	EntryStack       []string               `json:"entryStack"`
	StackPops        uint64                 `json:"stackPops"`
	StackAdditions   []string               `json:"stackAdditions"`
	ExitStack        []string               `json:"exitStack"`
	Instructions     map[string]Instruction `json:"instructions"`
	InstructionOrder []string               `json:"instructionOrder"`
}

type BlockRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Instruction struct {
	Op      string   `json:"op"`
	Args    []string `json:"args"`
	StackOp string   `json:"stackOp"`
}

type Function struct {
	Signature  string   `json:"signature"`
	EntryBlock string   `json:"entryBlock"`
	ExitBlock  string   `json:"exitBlock"`
	Body       []string `json:"body"`
}
