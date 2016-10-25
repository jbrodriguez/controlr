package main

import (
	"github.com/mcuadros/go-version"
	"testing"
)

// t.deepEqual(version_compare('6.2.0-beta1', '6.1.4', '>='), true, '6.2.0-beta1 not >= than 6.1.4')
// t.deepEqual(version_compare('6.2.0-rc1', '6.2.0-beta1', '>='), true, '6.2.0-rc1  6.2.0-beta1')
// t.deepEqual(version_compare('6.2', '6.2.0-rc1', '>='), true, '6.2 not >= than 6.2.0-rc1')
// t.deepEqual(version_compare('6.2.1-2016.09.22', '6.2.0-rc1', '>='), true, '6.2.1-2016.09.22 not >= than 6.2.0-rc1')
// t.deepEqual(version_compare('6.2.1-2016.09.22', '6.2', '>='), true, '6.2.1-2016.09.22 not >= than 6.2')
// t.deepEqual(version_compare('6.2.1-2016.09.22', '6.2.1', '>='), true, '6.2.1-2016.09.22 not >= than 6.2.1')
// t.deepEqual(version_compare('6.2.1-2016-09-22', '6.2.1', '>='), true, '6.2.1-2016-09-22 not >= than 6.2.1')
// t.deepEqual(version_compare('6.3.0', '6.2.1-2016.09.22', '>='), true, '6.3.0 not >= than 6.2.1-2016.09.22')
// t.deepEqual(version_compare('6.4.0', '6.3.0-2016.10.22', '>='), true, '6.4.0 not >= than 6.3.0-2016.10.22')
// t.deepEqual(version_compare('6.4.0', '6.2.1-2016.09.22', '>='), true, '6.4.0 not >= than 6.2.1-2016.09.22')

type compareTest struct {
	vlo string
	vhi string
	res bool
}

var compareTests = []compareTest{
	{vlo: "6.2.0-beta1", vhi: "6.1.4", res: true},
	{vlo: "6.2.0-rc1", vhi: "6.2.0-beta1", res: true},
	{vlo: "6.2", vhi: "6.2.0-rc1", res: true},
	{vlo: "6.2.1-2016.09.22", vhi: "6.2.0-rc1", res: true},
	{vlo: "6.2.1-2016.09.22", vhi: "6.2", res: true},
	{vlo: "6.2.1-2016.09.22", vhi: "6.2.1", res: true},
	{vlo: "6.2.1-2016-09-22", vhi: "6.2.1", res: true},
	{vlo: "6.3.0", vhi: "6.2.1-2016.09.22", res: true},
	{vlo: "6.4.0", vhi: "6.3.0-2016.10.22", res: true},
	{vlo: "6.4.0", vhi: "6.2.1-2016.09.22", res: true},
	{vlo: "6.4.0", vhi: "6.2.1-2016-09-22", res: true},
}

func TestCompare(t *testing.T) {
	for _, test := range compareTests {
		if res := version.Compare(test.vlo, test.vhi, ">="); res != test.res {
			t.Errorf("Comparing %q : %q, expected %t but got %t", test.vlo, test.vhi, test.res, res)
		}
		//Test counterpart
		if res := version.Compare(test.vlo, test.vhi, ">="); res == !test.res {
			t.Errorf("Comparing %q : %q, expected %t but got %t", test.vhi, test.vlo, !test.res, res)
		}
	}
}
