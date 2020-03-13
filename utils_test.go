package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNoSection(t *testing.T) {
	comment := "blah\nfoo\nbar\n"
	expected := []Section{
		{Title: "", Content: []string{"blah", "foo", "bar"}, Collapsed: false},
	}
	require.Equal(t, expected, ParseComment(comment))
}

func TestParse2OpenSections(t *testing.T) {
	comment := ":: Test1\nblah\n:: Test2\nfoo\n"
	expected := []Section{
		{Title: "Test1", Content: []string{"blah"}, Collapsed: false},
		{Title: "Test2", Content: []string{"foo"}, Collapsed: false},
	}
	require.Equal(t, expected, ParseComment(comment))
}

func TestParse1Open1Collapsed(t *testing.T) {
	comment := ":: Test1\nblah\n::- Test2\nfoo\n"
	expected := []Section{
		{Title: "Test1", Content: []string{"blah"}, Collapsed: false},
		{Title: "Test2", Content: []string{"foo"}, Collapsed: true},
	}
	require.Equal(t, expected, ParseComment(comment))
}

func TestParse1OpenEmpty1Collapsed(t *testing.T) {
	comment := ":: Test1\n::- Test2\nfoo\n"
	expected := []Section{
		{Title: "Test1", Collapsed: false},
		{Title: "Test2", Content: []string{"foo"}, Collapsed: true},
	}
	require.Equal(t, expected, ParseComment(comment))
}

func TestParse1Collapsed1Open(t *testing.T) {
	comment := "::- Test1\nblah\n:: Test2\nfoo\n"
	expected := []Section{
		{Title: "Test1", Content: []string{"blah"}, Collapsed: true},
		{Title: "Test2", Content: []string{"foo"}, Collapsed: false},
	}
	require.Equal(t, expected, ParseComment(comment))
}

func TestParse1CollapsedEmpty1Open(t *testing.T) {
	comment := "::- Test1\n:: Test2\nfoo\n"
	expected := []Section{
		{Title: "Test1", Collapsed: true},
		{Title: "Test2", Content: []string{"foo"}, Collapsed: false},
	}
	require.Equal(t, expected, ParseComment(comment))
}
