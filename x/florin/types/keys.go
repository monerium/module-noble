package types

const ModuleName = "florin"

var (
	AuthorityKey           = []byte("authority")
	AllowedDenomPrefix     = []byte("allowed_denom/")
	OwnerPrefix            = []byte("owner/")
	PendingOwnerPrefix     = []byte("pending_owner/")
	SystemPrefix           = []byte("system/")
	AdminPrefix            = []byte("admin/")
	MintAllowancePrefix    = []byte("mint_allowance/")
	MaxMintAllowancePrefix = []byte("max_mint_allowance/")
)

func AllowedDenomKey(denom string) []byte {
	return append(AllowedDenomPrefix, []byte(denom)...)
}

func OwnerKey(denom string) []byte {
	return append(OwnerPrefix, []byte(denom)...)
}

func PendingOwnerKey(denom string) []byte {
	return append(PendingOwnerPrefix, []byte(denom)...)
}

func SystemKey(denom string, address string) []byte {
	return append(append(SystemPrefix, []byte(denom)...), []byte(address)...)
}

func AdminKey(denom string, address string) []byte {
	return append(append(AdminPrefix, []byte(denom)...), []byte(address)...)
}

func MintAllowanceKey(denom string, address string) []byte {
	return append(append(MintAllowancePrefix, []byte(denom)...), []byte(address)...)
}

func MaxMintAllowanceKey(denom string) []byte {
	return append(MaxMintAllowancePrefix, []byte(denom)...)
}
