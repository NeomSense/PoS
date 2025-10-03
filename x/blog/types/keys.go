package types

// Module basics
const (
	ModuleName  = "blog"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_blog"
)

// Collection prefixes / counters
var (
	// ParamsKey is the prefix for module parameters
	ParamsKey = []byte("Params/")

	// PostKey is the prefix for individual post entries
	PostKey = []byte("Post/value/")

	// PostCountKey is the prefix for the post counter
	PostCountKey = []byte("Post/count/")
)
