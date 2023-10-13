package execute

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type Attributes struct {
	// Values specify which attributes the node in question should have.
	// At the moment we support strict equality only, so no `if RAM >= 16GB` types of conditions.
	Values []Parameter `json:"values,omitempty"`

	// Should we accept nodes whose attributes are not attested?
	AttestationRequired bool `json:"attestation_required,omitempty"`

	// Explicitly request specific attestors.
	Attestors AttributeAttestors `json:"attestors,omitempty"`
}

type AttributeAttestors struct {
	// Each of the listed attestors should be found.
	Each []peer.ID `json:"each,omitempty"`

	// Any one of these attestors should be found.
	OneOf []peer.ID `json:"one_of,omitempty"`
}
