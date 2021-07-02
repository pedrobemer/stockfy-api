package main

type Payload struct {
	Input actionNameArgs `json:"input"`
}

type stock struct {
	C  float32
	H  float32
	L  float32
	O  float32
	PC float32
	T  float32
}

type SampleInput struct {
	Symbol string
}

type SampleOutput struct {
	AccessToken string
}

type Mutation struct {
	ActionName *stock
}

type actionNameArgs struct {
	Arg1 SampleInput
}
