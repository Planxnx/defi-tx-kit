package enums

type TokenEvent string

const (
	TokenEventTransfer TokenEvent = "Transfer"
	TokenEventApproval TokenEvent = "Approval"
)

func (e TokenEvent) ToString() string {
	return string(e)
}

const (
	DECIMALS_WEI  = 18
	DECIMALS_GWEI = 9
)
