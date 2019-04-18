package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollision(t *testing.T) {

	//root := NewRootMessageFromFbs(&Root{})
	col := &CollisionMessage{
		Id: 10,
	}
	fbs := col.ToFbs(true)
	assert.NotNil(t, col)
	assert.NotNil(t, fbs)

	col = NewCollisionMessageFromFbs(fbs)
	assert.NotNil(t, col)
	assert.Equal(t, col.Id, int64(10))

}

func TestRoot(t *testing.T) {
	root := &RootMessage{
		Action: &ActionMessage{
			Type: ActionRez,
			Rez: &RezMessage{
				Id:     1,
				Name:   []uint8("aaaa"),
				ObjIds: []int64{1, 2, 3},
				Objes:  []*Vec3Message{},
			},
		},
		CurrentAt: &TimeMessage{
			Nano: int64(1),
			Unix: int64(1),
		},
		Collision: &CollisionMessage{
			Id: 10,
		},
	}
	fbs := root.ToFbs(true)
	root = NewRootMessageFromFbs(fbs)
	assert.Equal(t, root.CurrentAt.Nano, int64(1))
}
