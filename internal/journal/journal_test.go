// This source file is part of the attendance list project
// as a part of the go lecture by H. Neemann.
// For this reason you have no permission to use, modify or
// share this code without the agreement of the authors.
//
// Matriculation numbers of the authors: 5703004, 5736465

// Package journal provides functionality for writing text based journal files.
package journal

import (
	"os"
	"path"
	"testing"

	"github.com/dateiexplorer/attendancelist/internal/timeutil"
	"github.com/stretchr/testify/assert"
)

// Data

var persons = map[string]Person{
	"HM": {"Hans", "Müller", Address{"Feldweg", "12", "74722", "Buchen"}},
	"GM": {"Gisela", "Musterfrau", Address{"Musterstraße", "10", "74821", "Mosbach"}},
	"MM": {"Max", "Mustermann", Address{"Musterstraße", "20", "74821", "Mosbach"}},
	"AM": {"Anne", "Meier", Address{"Hauptstraße", "18", "74821", "Mosbach"}},
	"LM": {"Lieschen", "Müller", Address{"Lindenstraße", "15", "10115", "Berlin"}},
	"ON": {"Otto", "Normalverbraucher", Address{"Dieselstraße", "52", "70376", "Stuttgart"}},
}

var locs = map[string]Location{
	"DH": "DHBW Mosbach",
	"AM": "Alte Mälzerei",
}

// Functions

func TestReadJournal(t *testing.T) {
	expected := Journal{timeutil.NewDate(2021, 10, 15), []JournalEntry{
		{timeutil.NewTimestamp(2021, 10, 15, 6, 20, 13), "d61ec70b78628e15", Login, locs["DH"], persons["HM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 9, 15, 20), "989ce491d5df53c9", Login, locs["DH"], persons["GM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 12, 15, 30), "f797f342aebab436", Login, locs["DH"], persons["MM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 12, 17, 20), "1ce7549a51133e9f", Login, locs["DH"], persons["AM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 13, 30, 0), "68e7faee906ffd4c", Login, locs["DH"], persons["LM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 13, 40, 10), "d61ec70b78628e15", Logout, locs["DH"], persons["HM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 13, 40, 11), "5faacdf0e6e7b44a", Login, locs["AM"], persons["HM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 15, 42, 23), "68e7faee906ffd4c", Logout, locs["DH"], persons["LM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 16, 48, 21), "f797f342aebab436", Logout, locs["DH"], persons["MM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 16, 52, 0), "989ce491d5df53c9", Logout, locs["DH"], persons["GM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 17, 15, 22), "1ce7549a51133e9f", Logout, locs["DH"], persons["AM"]},
		{timeutil.NewTimestamp(2021, 10, 15, 17, 32, 45), "848dc86c0b5e62a0", Login, locs["AM"], persons["ON"]},
		{timeutil.NewTimestamp(2021, 10, 15, 19, 15, 12), "848dc86c0b5e62a0", Logout, locs["AM"], persons["ON"]},
	}}

	journal, err := ReadJournal("testdata", timeutil.NewDate(2021, 10, 15))
	assert.NoError(t, err)
	assert.Equal(t, expected, journal)
}

func TestReadNotExistingJournal(t *testing.T) {
	expected := Journal{timeutil.NewDate(2020, 9, 4), []JournalEntry{}}
	journal, err := ReadJournal("testdata", timeutil.NewDate(2020, 9, 4))
	assert.Error(t, err)
	assert.Equal(t, expected, journal)
}

func TestReadMalformedJournal(t *testing.T) {
	expected := []Journal{
		{timeutil.NewDate(2020, 1, 1), []JournalEntry{}},
		{timeutil.NewDate(2020, 1, 2), []JournalEntry{}},
	}

	for _, journal := range expected {
		actual, err := ReadJournal("testdata", journal.date)
		assert.NotErrorIs(t, err, os.ErrNotExist)
		assert.Error(t, err)
		assert.Equal(t, journal, actual)
	}
}

func TestGetVisitedLocationsForPerson(t *testing.T) {
	expected := []Location{
		"DHBW Mosbach",
	}

	journal, err := ReadJournal("testdata", timeutil.NewDate(2021, 10, 15))
	assert.NoError(t, err)
	p := persons["MM"]

	actual := journal.GetVisitedLocationsForPerson(&p)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetVisitedLocationForPersonMultipleLocations(t *testing.T) {
	expected := []Location{
		"DHBW Mosbach",
		"Alte Mälzerei",
	}

	journal, err := ReadJournal("testdata", timeutil.NewDate(2021, 10, 15))
	assert.NoError(t, err)
	p := persons["HM"]

	actual := journal.GetVisitedLocationsForPerson(&p)
	assert.NotNil(t, actual)

	for _, location := range expected {
		assert.Contains(t, actual, location)
	}
}

func TestGetVisitedLocationsForPersonNotExistingPerson(t *testing.T) {
	expected := []Location{}
	person := Person{"Susi", "Sorglos", Address{"Musterstraße", "20", "72327", "Musterstadt"}}

	journal, err := ReadJournal("testdata", timeutil.NewDate(2021, 10, 15))
	assert.NoError(t, err)

	actual := journal.GetVisitedLocationsForPerson(&person)
	assert.Equal(t, expected, actual)
}

func TestGetAttendanceEntriesForLocation(t *testing.T) {
	expected := AttendanceList{
		NewAttendanceEntry(persons["HM"], timeutil.NewTimestamp(2021, 10, 15, 13, 40, 11), timeutil.InvalidTimestamp),
		NewAttendanceEntry(persons["ON"], timeutil.NewTimestamp(2021, 10, 15, 17, 32, 45), timeutil.NewTimestamp(2021, 10, 15, 19, 15, 12)),
	}

	journal, err := ReadJournal("testdata", timeutil.NewDate(2021, 10, 15))
	assert.NoError(t, err)

	actual := journal.GetAttendanceListForLocation(locs["AM"])
	assert.Equal(t, expected, actual)
}

func TestGetAttendanceListForLocationNotExistingLocation(t *testing.T) {
	expected := AttendanceList{}
	journal, err := ReadJournal("testdata", timeutil.NewDate(2021, 10, 15))
	assert.NoError(t, err)

	actual := journal.GetAttendanceListForLocation("Night Club")
	assert.Equal(t, expected, actual)
}

func TestNextEntry(t *testing.T) {
	expected := [][]string{
		{"Hans", "Müller", "Feldweg", "12", "74722", "Buchen", "13:40:11", ""},
		{"Otto", "Normalverbraucher", "Dieselstraße", "52", "70376", "Stuttgart", "17:32:45", "19:15:12"},
	}

	list := AttendanceList{
		NewAttendanceEntry(persons["HM"], timeutil.NewTimestamp(2021, 10, 15, 13, 40, 11), timeutil.InvalidTimestamp),
		NewAttendanceEntry(persons["ON"], timeutil.NewTimestamp(2021, 10, 15, 17, 32, 45), timeutil.NewTimestamp(2021, 10, 15, 19, 15, 12)),
	}

	counter := 0
	for actual := range list.NextEntry() {
		assert.Equal(t, expected[counter], actual)
		counter++
	}

	assert.Equal(t, len(list), counter)
}

func TestHeader(t *testing.T) {
	expected := []string{"FirstName", "LastName", "Street", "Number", "ZipCode", "City", "Login", "Logout"}
	list := AttendanceList{}
	actual := list.Header()
	assert.Equal(t, expected, actual)
}

func TestNewPerson(t *testing.T) {
	p := NewPerson("Max", "Mustermann", "Musterstraße", "20", "74821", "Mosbach")
	assert.Equal(t, persons["MM"], p)
}

func TestWriteToJournalFile(t *testing.T) {
	timestamp := timeutil.NewTimestamp(2021, 10, 16, 15, 30, 0)
	expected := JournalEntry{timestamp, "aabbccddeeff", Login, locs["DH"], persons["MM"]}

	err := WriteToJournalFile("testdata", &expected)
	assert.NoError(t, err)

	journal, err := ReadJournal("testdata", timestamp.Date())
	assert.NoError(t, err)

	actual := journal.entries
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, expected, actual[0])

	err = os.Remove(path.Join("testdata", timestamp.Date().String()+journalFileExtension))
	assert.NoError(t, err)
}

func TestWriteToJournalFileAppendEntry(t *testing.T) {
	date := timeutil.NewDate(2021, 10, 16)
	expected := []JournalEntry{
		{timeutil.NewTimestamp(2021, 10, 16, 15, 30, 0), "aabbccddeeff", Login, locs["DH"], persons["MM"]},
		{timeutil.NewTimestamp(2021, 10, 16, 17, 20, 0), "aabbccddeeff", Logout, locs["DH"], persons["MM"]},
	}

	for _, e := range expected {
		err := WriteToJournalFile("testdata", &e)
		assert.NoError(t, err)
	}

	journal, err := ReadJournal("testdata", date)
	assert.NoError(t, err)

	actual := journal.entries
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, expected, actual)

	err = os.Remove(path.Join("testdata", date.String()+journalFileExtension))
	assert.NoError(t, err)
}

func TestEntries(t *testing.T) {
	expected := []JournalEntry{
		{timeutil.NewTimestamp(2021, 10, 16, 15, 30, 0), "aabbccddeeff", Login, locs["DH"], persons["MM"]},
		{timeutil.NewTimestamp(2021, 10, 16, 17, 20, 0), "aabbccddeeff", Logout, locs["DH"], persons["MM"]},
	}

	journal := Journal{timeutil.NewDate(2021, 10, 16), expected}
	actual := journal.Entries()
	assert.Equal(t, expected, actual)
}

func TestNewJournalEntry(t *testing.T) {
	expected := JournalEntry{timeutil.NewTimestamp(2021, 10, 16, 15, 30, 0), "aabbccddeeff", Login, locs["DH"], persons["MM"]}
	actual := NewJournalEntry(timeutil.NewTimestamp(2021, 10, 16, 15, 30, 0), "aabbccddeeff", Login, locs["DH"], persons["MM"])
	assert.Equal(t, expected, actual)
}
