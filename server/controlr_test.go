package main

import (
	"fmt"
	"jbrodriguez/controlr/plugin/server/src/lib"
	"regexp"
	"testing"

	"github.com/kless/osutil/user/crypt"
	"github.com/kless/osutil/user/crypt/md5_crypt"
	"github.com/kless/osutil/user/crypt/sha256_crypt"
	"github.com/kless/osutil/user/crypt/sha512_crypt"
	"github.com/mcuadros/go-version"
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

func TestRegex(t *testing.T) {
	re := regexp.MustCompile(`root:(\$(.*?)\$(.*?)\$.*?):`)
	test := `root:$1$alpha$beta:`

	saltString := ""
	actualHash := ""
	encType := ""
	for _, match := range re.FindAllStringSubmatch(test, -1) {
		actualHash = match[1]
		encType = match[2]
		saltString = match[3]
	}

	if saltString != `alpha` {
		t.Errorf("Comparing %s : expected %s but got %s", "saltString", "alpha", saltString)
	}

	if actualHash != `$1$alpha$beta` {
		t.Errorf("Comparing %s : expected %s but got %s", "actualHash", "$1$alpha$beta", actualHash)
	}

	if encType != `1` {
		t.Errorf("Comparing %s : expected %s but got %s", "encType", "1", encType)
	}
}

func TestEncrypt(t *testing.T) {
	re := regexp.MustCompile(`family:(\$(.*?)\$(.*?)\$.*?):`)

	shadowLine := `family:$5$tZ3/aLE/9CF$N9wOHr1PsCWwJVU4XR6uKidrnf2axbaxqyXLks0Aol1:17102:0:99999:7:::`

	saltString := ""
	actualHash := ""
	encType := ""
	for _, match := range re.FindAllStringSubmatch(shadowLine, -1) {
		actualHash = match[1]
		encType = match[2]
		saltString = match[3]
	}

	var crypto crypt.Crypter
	saltPrefix := ""
	// crypto := crypt.New(crypt.SHA256)
	// saltPrefix := sha256_crypt.MagicPrefix
	switch encType {
	case "1":
		crypto = crypt.New(crypt.MD5)
		saltPrefix = md5_crypt.MagicPrefix
		break
	case "5":
		crypto = crypt.New(crypt.SHA256)
		saltPrefix = sha256_crypt.MagicPrefix
		break
	case "6":
		crypto = crypt.New(crypt.SHA512)
		saltPrefix = sha512_crypt.MagicPrefix
		break
	default:
		t.Errorf("Unknown encryption type: (%s)", encType)
	}

	saltString = fmt.Sprintf("%s%s", saltPrefix, saltString)

	password := "password"
	shadowHash, err := crypto.Generate([]byte(password), []byte(saltString))
	if err != nil {
		t.Errorf("Unable to create hash: %s", err)
	}

	if shadowHash != actualHash {
		t.Errorf("Comparing %q: expected %q but got %q", "hash", actualHash, shadowHash)
	}
}

func TestEmhttp(t *testing.T) {
	re := regexp.MustCompile(`.*?emhttp(.*)$`)

	match := re.FindStringSubmatch("/usr/local/sbin/emhttp &")
	err, secure, port := lib.GetPort(match)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}
	if secure || port != "80" {
		t.Errorf("Secure != false(%t) - Port != 80(%s)", secure, port)
	}

	match = re.FindStringSubmatch("/usr/local/sbin/emhttp -p 88 &")
	err, secure, port = lib.GetPort(match)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}
	if secure || port != "88" {
		t.Errorf("Secure != false(%t) - Port != 88(%s)", secure, port)
	}

	match = re.FindStringSubmatch("/usr/local/sbin/emhttp -p ,448 &")
	err, secure, port = lib.GetPort(match)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}
	if !secure || port != "448" {
		t.Errorf("Secure != true(%t) - Port != 448(%s)", secure, port)
	}

	match = re.FindStringSubmatch("/usr/local/sbin/emhttp -p 87,445 -r &")
	err, secure, port = lib.GetPort(match)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}
	if !secure || port != "445" {
		t.Errorf("Secure != true(%t) - Port != 445(%s)", secure, port)
	}

	match = re.FindStringSubmatch("/usr/local/sbin/emhttp -rp 87,445 &")
	err, secure, port = lib.GetPort(match)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}
	if !secure || port != "445" {
		t.Errorf("Secure != true(%t) - Port != 445(%s)", secure, port)
	}

	match = re.FindStringSubmatch("/usr/local/sbin/emhttp -p 89,447 &")
	err, secure, port = lib.GetPort(match)
	if err != nil {
		t.Errorf("Failed: %s", err)
	}
	if secure || port != "89" {
		t.Errorf("Secure != false(%t) - Port != 89(%s)", secure, port)
	}
}
