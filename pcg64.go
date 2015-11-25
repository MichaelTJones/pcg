package pcg

// PCG Random Number Generation
// Developed by Melissa O'Neill <oneill@pcg-random.org>
// Paper and details at http://www.pcg-random.org
// Ported to Go by Michael Jones <michael.jones@gmail.com>

type PCG64 struct {
	lo *PCG32
	hi *PCG32
}

func NewPCG64() *PCG64 {
	return &PCG64{NewPCG32(), NewPCG32()}
}

func (p *PCG64) Seed(state1, state2, sequence1, sequence2 uint64) *PCG64 {
	mask := ^uint64(0) >> 1
	if sequence1&mask == sequence2&mask {
		sequence2 = ^sequence2
	}
	p.lo.Seed(state1, sequence1)
	p.hi.Seed(state2, sequence2)
	return p
}

func (p *PCG64) Random() uint64 {
	return uint64(p.hi.Random())<<32 | uint64(p.lo.Random())
}

func (p *PCG64) Bounded(bound uint64) uint64 {
	if bound == 0 {
		return 0
	}
	threshold := -bound % bound
	for {
		r := p.Random()
		if r >= threshold {
			return r % bound
		}
	}
}

func (p *PCG64) Advance(delta uint64) *PCG64 {
	p.lo.Advance(delta)
	p.hi.Advance(delta)
	return p
}

func (p *PCG64) Retreat(delta uint64) *PCG64 {
	return p.Advance(-delta)
}
