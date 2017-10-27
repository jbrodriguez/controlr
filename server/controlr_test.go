package main

import (
	"controlr/plugin/server/src/dto"
	"controlr/plugin/server/src/model"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/jbrodriguez/mlog"
	"github.com/kless/osutil/user/crypt"
	"github.com/kless/osutil/user/crypt/md5_crypt"
	"github.com/kless/osutil/user/crypt/sha256_crypt"
	"github.com/kless/osutil/user/crypt/sha512_crypt"
	"github.com/mcuadros/go-version"
)

func TestMain(m *testing.M) {
	mlog.Start(mlog.LevelInfo, "")

	os.Exit(m.Run())
}

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

func TestApcUps(t *testing.T) {
	lines := []string{
		"APC      : 001,048,1105",
		"DATE     : 2017-10-10 15:20:45 +0400",
		"HOSTNAME : MediaOne",
		"VERSION  : 3.14.14 (31 May 2016) slackware",
		"UPSNAME  : UPS_IDEN",
		"CABLE    : Custom Cable Smart",
		"DRIVER   : APC Smart UPS (any)",
		"UPSMODE  : Stand Alone",
		"STARTTIME: 2017-10-07 09:37:15 +0400",
		"MODEL    : Smart-UPS SC1500",
		"STATUS   : ONLINE",
		"LINEV    : 228.0 Volts",
		"LOADPCT  : 9.1 Percent",
		"BCHARGE  : 100.0 Percent",
		"TIMELEFT : 85.0 Minutes",
		"MBATTCHG : 10 Percent",
		"MINTIMEL : 10 Minutes",
		"MAXTIME  : 0 Seconds",
		"MAXLINEV : 228.0 Volts",
		"MINLINEV : 226.0 Volts",
		"OUTPUTV  : 228.0 Volts",
		"SENSE    : High",
		"DWAKE    : 0 Seconds",
		"DSHUTD   : 60 Seconds",
		"DLOWBATT : 2 Minutes",
		"LOTRANS  : 208.0 Volts",
		"HITRANS  : 253.0 Volts",
		"RETPCT   : 0.0 Percent",
		"ALARMDEL : 5 Seconds",
		"BATTV    : 26.8 Volts",
		"LINEFREQ : 50.0 Hz",
		"LASTXFER : Line voltage notch or spike",
		"NUMXFERS : 0",
		"TONBATT  : 0 Seconds",
		"CUMONBATT: 0 Seconds",
		"XOFFBATT : N/A",
		"SELFTEST : NO",
		"STESTI   : 336",
		"STATFLAG : 0x05000008",
		"REG1     : 0x00",
		"REG2     : 0x00",
		"REG3     : 0x00",
		"MANDATE  : 04/16/12",
		"SERIALNO : 5S1216T00762",
		"BATTDATE : 04/16/12",
		"NOMOUTV  : 230 Volts",
		"NOMBATTV : 24.0 Volts",
		"FIRMWARE : 738.3.I",
		"END APC  : 2017-10-10 15:20:51 +0400",
		"NOMPOWER : 865",
	}

	samplesExpected := []dto.Sample{
		dto.Sample{Key: "UPS STATUS", Value: "Online", Unit: "", Condition: "green"},
		dto.Sample{Key: "UPS LOAD", Value: "9.1", Unit: "%", Condition: "green"},
		dto.Sample{Key: "UPS CHARGE", Value: "100.0", Unit: "%", Condition: "green"},
		dto.Sample{Key: "UPS LEFT", Value: "85.0", Unit: "m", Condition: "green"},
		dto.Sample{Key: "UPS POWER", Value: "78.7", Unit: "w", Condition: "green"},
	}

	apc := model.NewApc()

	samplesActual := apc.Parse(lines)

	if !reflect.DeepEqual(samplesActual, samplesExpected) {
		t.Errorf("Comparing %q: expected\n %q\n but got\n %q\n", "apc", samplesExpected, samplesActual)
	}
}

func TestNutUps(t *testing.T) {
	lines := []string{
		"battery.charge: 100",
		"battery.charge.low: 30",
		"battery.runtime: 1000",
		"battery.type: PbAc",
		"device.mfr: MGE UPS SYSTEMS",
		"device.model: Nova 1100 AVR",
		"device.type: ups",
		"driver.name: usbhid-ups",
		"driver.parameter.pollfreq: 30",
		"driver.parameter.pollinterval: 2",
		"driver.parameter.port: auto",
		"driver.parameter.synchronous: no",
		"driver.version: 2.7.4.1",
		"driver.version.data: MGE HID 1.42",
		"driver.version.internal: 0.42",
		"outlet.1.status: on",
		"output.voltage: 230.0",
		"ups.delay.shutdown: 20",
		"ups.delay.start: 30",
		"ups.load: 6",
		"ups.mfr: MGE UPS SYSTEMS",
		"ups.model: Nova 1100 AVR",
		"ups.power.nominal: 1100",
		"ups.realpower.nominal: 300",
		"ups.productid: ffff",
		"ups.status: OL",
		"ups.timer.shutdown: -1",
		"ups.timer.start: -10",
		"ups.vendorid: 0463",
	}

	samplesExpected := []dto.Sample{
		dto.Sample{Key: "UPS CHARGE", Value: "100", Unit: "%", Condition: "green"},
		dto.Sample{Key: "UPS LEFT", Value: "16.7", Unit: "m", Condition: "green"},
		dto.Sample{Key: "UPS LOAD", Value: "6", Unit: "%", Condition: "green"},
		dto.Sample{Key: "UPS STATUS", Value: "Online", Unit: "", Condition: "green"},
		dto.Sample{Key: "UPS POWER", Value: "18.0", Unit: "w", Condition: "green"},
	}

	nut := model.NewNut()

	samplesActual := nut.Parse(lines)

	if !reflect.DeepEqual(samplesActual, samplesExpected) {
		t.Errorf("Comparing %q: expected\n %q\n but got\n %q\n", "nut", samplesExpected, samplesActual)
	}
}
