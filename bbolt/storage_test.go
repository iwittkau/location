package bbolt_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/iwittkau/location"

	"github.com/stretchr/testify/require"

	"github.com/iwittkau/location/bbolt"
)

func Test_Storage(t *testing.T) {
	s := bbolt.New()
	err := s.Open("testdata")
	require.NoError(t, err)
	defer s.Close()
	defer os.RemoveAll("testdata")

	c := location.CheckIn{
		ID:   "test",
		Name: "Test",
		Time: time.Unix(time.Now().Unix(), 0),
	}

	err = s.CreateCheckIn(c)
	require.NoError(t, err)
	list, err := s.ListCheckIns(100, 0)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.Len(t, list, 1)
	err = s.CreateCheckIn(c)
	require.Error(t, err)
	dbC, err := s.CheckIn(c.ID)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(c, dbC))
	s.DeleteCheckIn(c.ID)
	require.NoError(t, err)
	_, err = s.CheckIn(c.ID)
	require.Error(t, err)

	for i := 0; i < 100; i++ {
		c := location.CheckIn{
			ID:   "test",
			Name: "Test",
			Time: time.Unix(0, time.Now().UnixNano()),
		}
		err = s.CreateCheckIn(c)
		require.NoError(t, err)
	}

	list, err = s.ListCheckIns(100, 0)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.Len(t, list, 100)

	for i := range list {
		require.Exactly(t, c.ID, list[i].ID)
		require.Exactly(t, c.Name, list[i].Name)
		require.False(t, list[i].Time.IsZero())
	}

	list, err = s.ListCheckIns(100, 50)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.Len(t, list, 50)
	list, err = s.ListCheckIns(100, 1)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.Len(t, list, 99)
	list, err = s.ListCheckIns(100, 100)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.Len(t, list, 0)

}
