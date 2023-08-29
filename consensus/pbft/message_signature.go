package pbft

type signable interface {
	signableRecord() any
	setSignature(string)
	getSignature() string
}

var _ signable = (*PrePrepare)(nil)

// Returns the payload that is eligible to be signed. This means basically the PrePrepare struct, excluding the signature field.
func (p *PrePrepare) signableRecord() any {
	cp := *p
	cp.setSignature("")
	return cp
}

func (p *PrePrepare) setSignature(signature string) {
	p.Signature = signature
}

func (p PrePrepare) getSignature() string {
	return p.Signature
}

// Returns the payload that is eligible to be signed. This means basically the Prepare struct, excluding the signature field.
func (p *Prepare) signableRecord() any {
	cp := *p
	cp.setSignature("")
	return cp
}

func (p *Prepare) setSignature(signature string) {
	p.Signature = signature
}

func (p Prepare) getSignature() string {
	return p.Signature
}

// Returns the payload that is eligible to be signed. This means basically the Commit struct, excluding the signature field.
func (c *Commit) signableRecord() any {
	cp := *c
	cp.setSignature("")
	return cp
}

func (c *Commit) setSignature(signature string) {
	c.Signature = signature
}

func (c Commit) getSignature() string {
	return c.Signature
}
