package cfg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func setUp() {
	os.Setenv("CFG_REQFILE", "./cfg.req")
	os.Setenv("a", "a")
	os.Setenv("b", "b")
}

func tearDown() {
	os.Unsetenv("CFG_REQFILE")
	os.Remove("./cfg.req")
	os.Unsetenv("a")
	os.Unsetenv("b")
}

// TestValidateOrFail shows that we can validate that all required configs
// have been found by returning an error listing missing values on failure
func TestValidate(t *testing.T) {
	setUp()
	defer tearDown()
	_, _ = Find("a")
	_, _ = Find("b")
	_, _ = Find("c")
	_, _ = Find("d")

	// should not validate, should have a nice error
	valid, err := Valid()
	if err.Error() != "Missing keys [\"c\", \"d\"]" {
		t.Errorf(
			"cfg.Validate() => _, %s: want, %s",
			err.Error(),
			"Missing keys [\"c\", \"d\"]",
		)
	}
	if valid {
		t.Errorf("missing c and d but valid is true")
	}
}

func TestFindRemembersKeys(t *testing.T) {
	// set up
	setUp()
	defer tearDown()
	key := "key"
	value := "value"
	os.Setenv(key, value)
	defer os.Unsetenv(key)
	missingKey := "missingKey"
	defaultVal := "dval"
	defaultVal2 := "dval2"

	// run SUT
	_, _ = Find(key)
	_ = Get(missingKey, defaultVal)
	_ = Get(missingKey, defaultVal2)

	// should only keep one copy of any default value it sees
	_ = Get(missingKey, defaultVal2)
	_ = Get(missingKey, defaultVal2)
	_ = Get(missingKey, defaultVal2)

	// confirm cfg.req has both keys and the default
	file, err := os.Open(os.Getenv("CFG_REQFILE"))
	panicOn(err)
	defer file.Close()
	reqs := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		kv := strings.Split(scanner.Text(), "=")
		reqs[kv[0]] = kv[1]
	}
	if reqs[key] != "" {
		t.Errorf("%s should be empty in reqfile\n%s", key, reqs)
	}
	if reqs[missingKey] != fmt.Sprintf("%s,%s", defaultVal, defaultVal2) {
		t.Errorf(
			"%s should be %s in reqfile\n%s",
			key,
			defaultVal,
			reqs,
		)
	}
}

func TestGetFound(t *testing.T) {
	// set up
	setUp()
	defer tearDown()
	key := "key"
	default_ := "dval"
	value := "val"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	// run SUT
	observedVal := Get(key, default_)

	// confirm
	if observedVal != value {
		t.Errorf(
			"With default of %s, should ovserve %s but saw %s",
			default_,
			value,
			observedVal,
		)
	}
}

func TestGetMissing(t *testing.T) {
	// set up
	setUp()
	defer tearDown()
	default_ := "dval"
	key := "key"

	// run SUT
	observedVal := Get(key, default_)

	//  confirm
	if observedVal != default_ {
		t.Errorf(
			"Observed value should be the default but %s != %s",
			observedVal,
			default_,
		)
	}
}

func TestFindMissing(t *testing.T) {
	// set up
	setUp()
	defer tearDown()
	missingKey := "missingKey"

	// run SUT
	observedVal, err := Find(missingKey)

	// confirm
	if err == nil {
		t.Errorf("Find should error on missing key but did not")
	}
	if observedVal != "" {
		t.Errorf("Observed value should be empty when not found")
	}
}

func TestFindExistingFromEnv(t *testing.T) {
	// set up
	setUp()
	defer tearDown()
	testVal := "testval"
	os.Setenv("CFG_TEST_KEY", testVal)
	defer os.Unsetenv("CFG_TEST_KEY")

	// run SUT
	observedVal, err := Find("CFG_TEST_KEY")

	// confirm
	if err != nil {
		t.Errorf("Error from Find not nil on get from env")
	}
	if observedVal != testVal {
		t.Errorf(
			"Observed value not equal to env value. %s != %s",
			observedVal,
			testVal,
		)
	}
}
