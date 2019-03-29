package ventures

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVentureStore_GetAll_1(t *testing.T) {
	store := NewVentureStore()
	a := Venture{
		ID:          "1",
		Description: "1",
		State:       "1",
	}
	b := Venture{
		ID:          "2",
		Description: "2",
		State:       "2",
	}

	store.items["1"] = a
	store.items["2"] = b

	s := store.GetAll()

	require.Len(t, s, 2)
	assert.Contains(t, s, a)
	assert.Contains(t, s, b)
}

func TestVentureStore_GetAll_2(t *testing.T) {
	store := NewVentureStore()
	s := store.GetAll()
	require.Empty(t, s)
}

func TestVentureStore_GetAllAlive_1(t *testing.T) {
	store := NewVentureStore()
	a := Venture{
		ID:          "1",
		Description: "1",
		State:       "1",
		IsAlive:     true,
	}
	b := Venture{
		ID:          "2",
		Description: "2",
		State:       "2",
		IsAlive:     false,
	}

	store.items["1"] = a
	store.items["2"] = b

	s := store.GetAllAlive()

	require.Len(t, s, 1)
	assert.Contains(t, s, a)
}

func TestVentureStore_GetAllAlive_2(t *testing.T) {
	store := NewVentureStore()
	s := store.GetAllAlive()
	require.Empty(t, s)
}

func TestVentureStore_Get_1(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{
		ID:          "1",
		Description: "1",
		State:       "1",
	}
	bIn := Venture{
		ID:          "2",
		Description: "2",
		State:       "2",
	}

	store.items["1"] = aIn
	store.items["2"] = bIn

	aOut, ok := store.Get("1")
	require.True(t, ok)
	assert.Equal(t, aIn, aOut)

	bOut, ok := store.Get("2")
	require.True(t, ok)
	assert.Equal(t, bIn, bOut)
}

func TestVentureStore_Get_2(t *testing.T) {
	store := NewVentureStore()
	_, ok := store.Get("1")
	require.False(t, ok)
}

func TestVentureStore_Get_3(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{
		ID:          "1",
		Description: "1",
		State:       "1",
	}
	bIn := Venture{
		ID:          "2",
		Description: "2",
		State:       "2",
	}

	store.items["1"] = aIn
	store.items["2"] = bIn

	_, ok := store.Get("3")
	require.False(t, ok)
}

func TestVentureStor_Add_1(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{
		Description: "description",
		State:       "state",
	}

	aOut := store.Add(aIn)
	assert.Len(t, store.items, 1)
	assert.NotEmpty(t, aOut.ID)
	assert.Equal(t, "description", aOut.Description)
	assert.Equal(t, "state", aOut.State)
}

func TestVentureStor_Add_2(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{
		Description: "description",
		State:       "state",
	}
	bIn := Venture{
		Description: "description",
		State:       "state",
	}

	aOut := store.Add(aIn)
	bOut := store.Add(bIn)
	assert.Len(t, store.items, 2)
	assert.NotEmpty(t, aOut.ID)
	assert.NotEmpty(t, bOut.ID)
	assert.NotEqual(t, aOut.ID, bOut.ID)
	assert.Equal(t, "description", aOut.Description)
	assert.Equal(t, "description", bOut.Description)
	assert.Equal(t, "state", aOut.State)
	assert.Equal(t, "state", bOut.State)
}

func TestVentureStore_Update_1(t *testing.T) {
	store := NewVentureStore()
	a := Venture{
		ID:          "1",
		Description: "original",
	}
	store.items["1"] = a

	bIn := Venture{
		ID:          "1",
		Description: "new",
	}

	ok := store.Update(bIn)
	require.True(t, ok)

	bOut, ok := store.items["1"]
	require.True(t, ok)
	assert.Equal(t, bIn, bOut)
}

func TestVentureStore_Update_2(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{
		ID:          "1",
		Description: "original",
	}

	ok := store.Update(aIn)
	require.False(t, ok)
}

func TestVentureStore_GenNewID_1(t *testing.T) {
	store := NewVentureStore()
	a := store._genNewID()
	assert.Equal(t, "1", a)
}

func TestVentureStore_GenNewID_2(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{}
	store.items["1"] = aIn

	a := store._genNewID()
	assert.Equal(t, "2", a)
}

func TestVentureStore_GenNewID_3(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{}
	store.items["1"] = aIn
	store.items["2"] = aIn
	store.items["3"] = aIn

	a := store._genNewID()
	assert.Equal(t, "4", a)
}

func TestVentureStore_GenNewID_4(t *testing.T) {
	store := NewVentureStore()
	aIn := Venture{}
	store.items["3"] = aIn

	a := store._genNewID()
	assert.Equal(t, "1", a)
}
