package types

const ModuleName = "florin"

var (
	OwnerKey            = []byte("owner")
	PendingOwnerKey     = []byte("pending_owner")
	SystemPrefix        = []byte("system/")
	AdminPrefix         = []byte("admin/")
	MintAllowancePrefix = []byte("mint_allowance/")
	MaxMintAllowanceKey = []byte("max_mint_allowance")
)

func SystemKey(address string) []byte {
	return append(SystemPrefix, address...)
}

func AdminKey(address string) []byte {
	return append(AdminPrefix, []byte(address)...)
}

func MintAllowanceKey(address string) []byte {
	return append(MintAllowancePrefix, []byte(address)...)
}
