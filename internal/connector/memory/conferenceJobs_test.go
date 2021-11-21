package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConferenceJobs(t *testing.T) {
	store := NewConferenceJobs()

	err := store.Store("1", "rec1")
	assert.Nil(t, err)
	err = store.Store("1", "rec2")
	assert.Nil(t, err)
	err = store.Store("1", "rec3")
	assert.Nil(t, err)

	err = store.Store("2", "rec1")
	assert.Nil(t, err)
	err = store.Store("2", "rec2")
	assert.Nil(t, err)

	isDone, err := store.IsDone("1")
	assert.Nil(t, err)
	assert.False(t, isDone)

	err = store.DoneJob("1", "rec1")
	assert.Nil(t, err)

	isDone, err = store.IsDone("1")
	assert.Nil(t, err)
	assert.False(t, isDone)

	err = store.DoneJob("1", "rec2")
	assert.Nil(t, err)
	err = store.DoneJob("1", "rec3")
	assert.Nil(t, err)

	isDone, err = store.IsDone("1")
	assert.Nil(t, err)
	assert.True(t, isDone)
}
