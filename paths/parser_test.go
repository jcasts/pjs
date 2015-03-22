package paths

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

func testParsePath(t *testing.T, pStr string) *path {
  p, err := parsePath(pStr)
  testAssertNil(t, err)
  return p
}

func TestParseInvalid(t *testing.T) {
  var err error
  _, err = parsePath("")
  testAssertNotNil(t, err)
  testAssertEqual(t, "Paths can't be empty", err.Error())

  _, err = parsePath("foo=**/thing")
  testAssertNotNil(t, err)
  testAssertEqual(t, "Invalid path value '**'. Use '\\*\\*' to match characters.", err.Error())

  _, err = parsePath("foo=../thing")
  testAssertNotNil(t, err)
  testAssertEqual(t, "Invalid path value '..'. Use '\\.\\.' to match characters.", err.Error())

  _, err = parsePath("foo=blah=fskj/thing")
  testAssertNotNil(t, err)
  testAssertEqual(t, "Multiple '=' invalid. Use \\= to match character.", err.Error())
}

func TestParseAny(t *testing.T) {
  p := testParsePath(t, "foo/*/bar")
  testAssertEqual(t, 3, len(p.tokens))
  testAssertTrue(t, p.tokens[0].matches("foo", nil))
  testAssertFalse(t, p.tokens[0].matches("fo", nil))
  testAssertTrue(t, p.tokens[1].matches(1, nil))
  testAssertTrue(t, p.tokens[1].matches("anything!", nil))
  testAssertTrue(t, p.tokens[1].isAny())
  testAssertTrue(t, p.tokens[2].matches("bar", nil))
  testAssertFalse(t, p.tokens[2].matches("b", nil))
}

func TestParseValues(t *testing.T) {
  p := testParsePath(t, "foo/*=thing")
  testAssertEqual(t, 2, len(p.tokens))
  testAssertTrue(t, p.tokens[1].matches(1, "thing"))
  testAssertTrue(t, p.tokens[1].matches("anything!", "thing"))
  testAssertFalse(t, p.tokens[1].matches(1, "blah"))
}

func TestParseInt(t *testing.T) {
  p := testParsePath(t, "foo/2")
  testAssertEqual(t, 2, len(p.tokens))
  testAssertEqual(t, intMatcher, p.tokens[1].keyMatcher.matcherType)
  testAssertTrue(t, p.tokens[1].matches(2, "thing"))
  testAssertTrue(t, p.tokens[1].matches("2", "thing"))
  testAssertFalse(t, p.tokens[1].matches(1, "blah"))
}

func TestParseRanges(t *testing.T) {
  p := testParsePath(t, "foo/1..3")
  testAssertEqual(t, 2, len(p.tokens))
  testAssertTrue(t, p.tokens[1].matches(1, nil))
  testAssertTrue(t, p.tokens[1].matches(2, nil))
  testAssertTrue(t, p.tokens[1].matches(3, nil))
  testAssertFalse(t, p.tokens[1].matches(0, nil))
  testAssertFalse(t, p.tokens[1].matches(4, nil))

  p = testParsePath(t, "foo=-12..-11")
  testAssertEqual(t, 1, len(p.tokens))
  testAssertFalse(t, p.tokens[0].matches("fo", -11))
  testAssertTrue(t, p.tokens[0].matches("foo", -11))
  testAssertTrue(t, p.tokens[0].matches("foo", -12))
  testAssertFalse(t, p.tokens[0].matches("foo", -10))
  testAssertFalse(t, p.tokens[0].matches("foo", -13))
}

func TestParseWildcard(t *testing.T) {
  p := testParsePath(t, "*bar*foo/")
  testAssertEqual(t, 1, len(p.tokens))
  testAssertTrue(t, p.tokens[0].matches("barfoo", nil))
  testAssertTrue(t, p.tokens[0].matches("bar_foo", nil))
  testAssertTrue(t, p.tokens[0].matches("fizz_bar_foo", nil))
  testAssertFalse(t, p.tokens[0].matches("fizz_bar_fo", nil))
  testAssertFalse(t, p.tokens[0].matches("ar_foo", nil))

  p = testParsePath(t, "*5")
  testAssertEqual(t, 1, len(p.tokens))
  testAssertTrue(t, p.tokens[0].matches("15", nil))
  testAssertTrue(t, p.tokens[0].matches(15, nil))
}

func TestParseParent(t *testing.T) {
  p := testParsePath(t, "foo/bar/..")
  testAssertEqual(t, 3, len(p.tokens))
  testAssertTrue(t, p.tokens[2].followParent())
  testAssertFalse(t, p.tokens[2].matches(1, "blah"))
}

func TestParseRecursive(t *testing.T) {
  p := testParsePath(t, "foo/bar/**")
  testAssertEqual(t, 3, len(p.tokens))
  testAssertTrue(t, p.tokens[2].isRecursive())
  testAssertTrue(t, p.tokens[2].matches(1, "blah"))

  p = testParsePath(t, "foo/**=bar")
  testAssertEqual(t, 2, len(p.tokens))
  testAssertTrue(t, p.tokens[1].isRecursive())
  testAssertTrue(t, p.tokens[1].matches(1, "bar"))
  testAssertFalse(t, p.tokens[1].matches(1, "fizz"))

  p = testParsePath(t, "foo/**/bar")
  testAssertEqual(t, 2, len(p.tokens))
  testAssertTrue(t, p.tokens[1].isRecursive())
  testAssertTrue(t, p.tokens[1].matches("bar", nil))
  testAssertFalse(t, p.tokens[1].matches("fizz", nil))

  p = testParsePath(t, "foo/**/**/bar")
  testAssertEqual(t, 2, len(p.tokens))
  testAssertTrue(t, p.tokens[1].isRecursive())
  testAssertTrue(t, p.tokens[1].matches("bar", nil))
  testAssertFalse(t, p.tokens[1].matches("fizz", nil))
}
