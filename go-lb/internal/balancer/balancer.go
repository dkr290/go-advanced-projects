package balancer

import (
	"math/rand"
	"sync/atomic"
)

type Balancer interface {
	Next() int // returns backend index
}

type RoundRobin struct {
	n   int
	idx uint64
}

func NewRoundRobin(n int) *RoundRobin {
	return &RoundRobin{n: n}
}

func (r *RoundRobin) Next() int {
	return int(atomic.AddUint64(&r.idx, 1)-1) % r.n
}

type Weighted struct {
	weights []int
	total   int
	pos     int
	count   int
}

func NewWeighted(weights []int) *Weighted {
	total := 0
	for _, w := range weights {
		total += w
	}
	return &Weighted{weights: weights, total: total}
}

func (w *Weighted) Next() int {
	for {
		if w.count < w.weights[w.pos] {
			w.count++
			return w.pos
		}
		w.count = 1
		w.pos = (w.pos + 1) % len(w.weights)
	}
}

type Percentage struct {
	percentages []int
}

func NewPercentage(percentages []int) *Percentage {
	return &Percentage{percentages: percentages}
}

func (p *Percentage) Next() int {
	r := rand.Intn(100)
	sum := 0
	for i, pct := range p.percentages {
		sum += pct
		if r < sum {
			return i
		}
	}
	return len(p.percentages) - 1
}
