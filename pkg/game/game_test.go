package game

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestCmpEquals(t *testing.T) {
	assert := assert.New(t)
	DefineGame(6, 4)

	z := ZeroPosition()
	assert.True(cmp.Equal(z, z))

	s := StartPosition()
	d := DiagnosticPosition()
	assert.False(cmp.Equal(s, d))
}

func TestCreatePosition(t *testing.T) {
	assert := assert.New(t)
	DefineGame(3, 4)

	d := DiagnosticPosition()
	p := CreatePosition(0, 1, 2, 3, 0, 1, 2, 3)
	assert.True(cmp.Equal(d, p))
	s := StartPosition()
	p = CreatePosition(0, 4, 4, 4, 0, 4, 4, 4)
	assert.True(cmp.Equal(s, p))
}

func TestNearFar(t *testing.T) {
	assert := assert.New(t)
	DefineGame(3, 4)

	p := CreatePosition(1, 3, 5, 7, 2, 4, 6, 8)
	n := &Side{[]int{1, 3, 5, 7}}
	f := &Side{[]int{2, 4, 6, 8}}
	assert.True(cmp.Equal(n, p.near()))
	assert.True(cmp.Equal(f, p.far()))
}

func TestCsv(t *testing.T) {
	assert := assert.New(t)
	DefineGame(3, 4)

	p := CreatePosition(1, 3, 5, 7, 2, 4, 6, 8)
	s := p.AsCsv()
	assert.Equal(s, "1,3,5,7,2,4,6,8")

	clone := CreatePositionCsv(s)
	assert.False(p == clone)
	assert.Equal(p, clone)
	assert.True(cmp.Equal(p, clone))
}
