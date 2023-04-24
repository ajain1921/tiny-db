package server

type Reply string

type BeginArgs struct {
	ClientId  string
	Timestamp int64
}

type UpdateArgs struct {
	ClientId  string
	Branch    string
	Account   string
	Amount    int
	Timestamp int64
}

type AbortArgs struct {
	Timestamp int64
}

type BalanceArgs struct {
	ClientId  string
	Branch    string
	Account   string
	Timestamp int64
}
