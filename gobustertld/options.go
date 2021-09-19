package gobustertld

import (
	"time"
)

// OptionsTLD holds all options for the tld plugin
type OptionsTLD struct {
	Domain         string
	ShowIPs        bool
	ShowCNAME      bool
	WildcardForced bool
	Resolver       string
	Timeout        time.Duration
}

// NewOptionsTLD returns a new initialized OptionsDNS
func NewOptionsTLD() *OptionsTLD {
	return &OptionsTLD{}
}
