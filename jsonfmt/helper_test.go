package jsonfmt

import (
  "testing"
)

func testAssertEqual(t *testing.T, v1, v2 interface{}) {
  if v1 != v2 { t.Fatalf("Expected %v but was %v", v1, v2) }
}

func testAssertTrue(t *testing.T, value bool) {
  if !value { t.Fatalf("Expected value to be true") }
}

func testAssertFalse(t *testing.T, value bool) {
  if value { t.Fatalf("Expected value to be false") }
}

func testAssertNil(t *testing.T, value interface{}) {
  if value != nil { t.Fatalf("Expected value to be nil") }
}

func testAssertNotNil(t *testing.T, value interface{}) {
  if value == nil { t.Fatalf("Expected value to be not nil") }
}
