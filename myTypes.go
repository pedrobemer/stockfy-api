package main

type Payload struct {
	Input actionNameArgs `json:"input"`
}

type Stock struct {
	C  float64
	H  float64
	L  float64
	O  float64
	PC float64
	T  float64
}

type SampleInput struct {
	Symbol string
}

type SampleOutput struct {
	AccessToken string
}

type Mutation struct {
	ActionName *Stock
}

type actionNameArgs struct {
	Arg1 SampleInput
}
