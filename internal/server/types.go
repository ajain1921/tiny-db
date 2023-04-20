package server

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

// type Reply struct {
// }

type Reply string

type BeginArgs struct {
	ClientId string
}

type DepositArgs struct {
	ClientId string
	Branch   string
	Account  string
	Amount   int
}
