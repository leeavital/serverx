package main

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/require"
)

func TestGetStructTags(t *testing.T) {
	tags := getStructTypes(reflect.TypeOf(MyArgs{}))
	require.Contains(t, tags, "json")

	getStructTypes(reflect.TypeOf(&MyArgs{}))
}
