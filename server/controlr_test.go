package main

import (
	"controlr/plugin/server/src/dto"
	"controlr/plugin/server/src/plugins/sensor"
	"controlr/plugin/server/src/plugins/ups"
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
	case "5":
		crypto = crypt.New(crypt.SHA256)
		saltPrefix = sha256_crypt.MagicPrefix
	case "6":
		crypto = crypt.New(crypt.SHA512)
		saltPrefix = sha512_crypt.MagicPrefix
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

	apc := ups.NewApc()

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

	nut := ups.NewNut()

	samplesActual := nut.Parse(lines)

	if !reflect.DeepEqual(samplesActual, samplesExpected) {
		t.Errorf("Comparing %q: expected\n %q\n but got\n %q\n", "nut", samplesExpected, samplesActual)
	}
}

func TestSystemSensor(t *testing.T) {
	lines := []string{
		"coretemp-isa-0000",
		"MB Temp:      +47.0°C  (high = +79.0°C, crit = +85.0°C)",
		"CPU Temp:     +46.0°C  (high = +79.0°C, crit = +85.0°C)",
		"Core 1:       +46.0°C  (high = +79.0°C, crit = +85.0°C)",
		"",
		"nct6776-isa-0290",
		"Vcore:         +0.87 V  (min =  +0.00 V, max =  +1.74 V)",
		"in1:           +1.86 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"AVCC:          +3.39 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"+3.3V:         +3.38 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"in4:           +1.06 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"in5:           +1.69 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"in6:           +0.86 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"3VSB:          +3.46 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"Vbat:          +3.30 V  (min =  +0.00 V, max =  +0.00 V)  ALARM",
		"Array Fan:    1934 RPM  (min =    0 RPM)",
		"fan2:         2347 RPM  (min =    0 RPM)",
		"fan3:         2011 RPM  (min =    0 RPM)",
		"SYSTIN:        +40.0°C  (high =  +0.0°C, hyst =  +0.0°C)  ALARM  sensor = thermistor",
		"CPUTIN:        +40.5°C  (high = +80.0°C, hyst = +75.0°C)  sensor = thermistor",
		"AUXTIN:        +32.0°C  (high = +80.0°C, hyst = +75.0°C)  sensor = thermistor",
		"PECI Agent 0:  +47.0°C  (high = +80.0°C, hyst = +75.0°C)",
		"						(crit = +85.0°C)",
		"intrusion0:   ALARM",
		"intrusion1:   ALARM",
		"beep_enable:  disabled",
	}

	samplesExpected := []dto.Sample{
		dto.Sample{Key: "BOARD", Value: "47", Unit: "C", Condition: "neutral"},
		dto.Sample{Key: "CPU", Value: "46", Unit: "C", Condition: "neutral"},
		dto.Sample{Key: "FAN", Value: "1934", Unit: "rpm", Condition: "neutral"},
	}

	sensor := sensor.NewSystemSensor()
	prefs := dto.Prefs{Number: ".,", Unit: "C"}

	samplesActual := sensor.Parse(prefs, lines)

	if !reflect.DeepEqual(samplesActual, samplesExpected) {
		t.Errorf("Comparing %q: expected\n %q\n but got\n %q\n", "SystemSensor C", samplesExpected, samplesActual)
	}

	samplesExpected = []dto.Sample{
		dto.Sample{Key: "BOARD", Value: "79", Unit: "F", Condition: "neutral"},
		dto.Sample{Key: "CPU", Value: "78", Unit: "F", Condition: "neutral"},
		dto.Sample{Key: "FAN", Value: "1934", Unit: "rpm", Condition: "neutral"},
	}

	prefs = dto.Prefs{Number: ".,", Unit: "F"}

	samplesActual = sensor.Parse(prefs, lines)

	if !reflect.DeepEqual(samplesActual, samplesExpected) {
		t.Errorf("Comparing %q: expected\n %q\n but got\n %q\n", "SystemSensor F", samplesExpected, samplesActual)
	}

}

func TestIpmiSensor(t *testing.T) {
	lines := []string{
		"4,CPU Temp,Temperature,Nominal,38.00,C,'OK'",
		"71,System Temp,Temperature,Nominal,30.00,C,'OK'",
		"138,Peripheral Temp,Temperature,Nominal,39.00,C,'OK'",
		"205,PCH Temp,Temperature,Nominal,45.00,C,'OK'",
		"272,P1-DIMMA1 Temp,Temperature,Nominal,30.00,C,'OK'",
		"339,P1-DIMMA2 Temp,Temperature,Nominal,31.00,C,'OK'",
		"406,P1-DIMMB1 Temp,Temperature,Nominal,30.00,C,'OK'",
		"473,P1-DIMMB2 Temp,Temperature,Nominal,30.00,C,'OK'",
		"540,FAN1,Fan,Nominal,1100.00,RPM,'OK'",
		"607,FAN2,Fan,Nominal,700.00,RPM,'OK'",
		"674,FAN3,Fan,N/A,N/A,RPM,N/A",
		"741,FAN4,Fan,Nominal,1100.00,RPM,'OK'",
		"808,FANA,Fan,Nominal,300.00,RPM,'OK'",
		"875,Vcpu,Voltage,Nominal,1.84,V,'OK'",
		"942,VDIMM,Voltage,Nominal,1.34,V,'OK'",
		"1009,12V,Voltage,Nominal,12.00,V,'OK'",
		"1076,5VCC,Voltage,Nominal,4.95,V,'OK'",
		"1143,3.3VCC,Voltage,Nominal,3.27,V,'OK'",
		"1210,VBAT,Voltage,Nominal,2.97,V,'OK'",
		"1277,5V Dual,Voltage,Nominal,5.03,V,'OK'",
		"1344,3.3V AUX,Voltage,Nominal,3.28,V,'OK'",
		"1411,1.2V BMC,Voltage,Nominal,1.26,V,'OK'",
		"1478,1.05V PCH,Voltage,Nominal,1.05,V,'OK'",
		"1545,Chassis Intru,Physical Security,Nominal,N/A,N/A,'OK'",
	}

	samplesExpected := []dto.Sample{
		dto.Sample{Key: "CPU", Value: "38", Unit: "C", Condition: "neutral"},
		dto.Sample{Key: "BOARD", Value: "30", Unit: "C", Condition: "neutral"},
		dto.Sample{Key: "FAN", Value: "1100", Unit: "rpm", Condition: "neutral"},
	}

	sensor := sensor.NewIpmiSensor()
	prefs := dto.Prefs{Number: ".,", Unit: "C"}

	samplesActual := sensor.Parse(prefs, lines)

	if !reflect.DeepEqual(samplesActual, samplesExpected) {
		t.Errorf("Comparing %q: expected\n %q\n but got\n %q\n", "IpmiSensor C", samplesExpected, samplesActual)
	}

	samplesExpected = []dto.Sample{
		dto.Sample{Key: "CPU", Value: "70", Unit: "F", Condition: "neutral"},
		dto.Sample{Key: "BOARD", Value: "62", Unit: "F", Condition: "neutral"},
		dto.Sample{Key: "FAN", Value: "1100", Unit: "rpm", Condition: "neutral"},
	}

	prefs = dto.Prefs{Number: ".,", Unit: "F"}

	samplesActual = sensor.Parse(prefs, lines)

	if !reflect.DeepEqual(samplesActual, samplesExpected) {
		t.Errorf("Comparing %q: expected\n %q\n but got\n %q\n", "IpmiSensor F", samplesExpected, samplesActual)
	}

}
