package convert

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dateiexplorer/attendancelist/internal/journal"
	"github.com/dateiexplorer/attendancelist/internal/timeutil"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("Test")

type errorWriter struct{}

func (e errorWriter) Write(b []byte) (int, error) {
	return 0, errTest
}

func TestToCSV(t *testing.T) {
	expected := `FirstName,LastName,Street,Number,ZipCode,City,Login,Logout
Hans,Müller,Feldweg,12,74722,Buchen,13:40:11,
Otto,Normalverbraucher,Dieselstraße,52,70376,Stuttgart,17:32:45,19:15:12
`

	// Prepare list
	list := journal.AttendanceList{
		journal.NewAttendanceEntry(journal.NewPerson("Hans", "Müller", "Feldweg", "12", "74722", "Buchen"), timeutil.NewTimestamp(2021, 10, 15, 13, 40, 11), timeutil.InvalidTimestamp),
		journal.NewAttendanceEntry(journal.NewPerson("Otto", "Normalverbraucher", "Dieselstraße", "52", "70376", "Stuttgart"), timeutil.NewTimestamp(2021, 10, 15, 17, 32, 45), timeutil.NewTimestamp(2021, 10, 15, 19, 15, 12)),
	}

	actual := new(bytes.Buffer)
	err := ToCSV(actual, list)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual.String())
}

func TestEmptyAttendanceListToCSV(t *testing.T) {
	expected := `FirstName,LastName,Street,Number,ZipCode,City,Login,Logout
`

	// Prepare list
	list := journal.AttendanceList{}

	actual := new(bytes.Buffer)
	err := ToCSV(actual, list)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual.String())
}

func TestToCSVFailedToWrite(t *testing.T) {
	actual := errorWriter{}
	err := ToCSV(actual, journal.AttendanceList{})

	assert.ErrorIs(t, err, errTest)
}