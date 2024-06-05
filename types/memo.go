package types

type Memo struct {
	DepositForBurn           *DepositForBurn           `json:"deposit_for_burn"`
	DepositForBurnWithCaller *DepositForBurnWithCaller `json:"deposit_for_burn_with_caller"`
}

type DepositForBurn struct {
	// TODO: Add amount and fee_recipient fields.
	DestinationDomain uint32 `json:"destination_domain"`
	MintRecipient     []byte `json:"mint_recipient"`
}

type DepositForBurnWithCaller struct {
	DepositForBurn
	DestinationCaller []byte `json:"destination_caller"`
}
