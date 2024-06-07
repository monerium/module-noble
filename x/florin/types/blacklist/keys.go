package blacklist

const SubmoduleName = "florin-blacklist"

var (
	OwnerKey        = []byte("blacklist/owner")
	PendingOwnerKey = []byte("blacklist/pending_owner")
	AdminPrefix     = []byte("blacklist/admin/")
	AdversaryPrefix = []byte("blacklist/adversary/")
)

func AdminKey(address string) []byte {
	return append(AdminPrefix, []byte(address)...)
}

func AdversaryKey(address string) []byte {
	return append(AdversaryPrefix, []byte(address)...)
}
