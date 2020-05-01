package attest

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestSecureBoot(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/windows_gcp_shielded_vm.json")
	if err != nil {
		t.Fatalf("reading test data: %v", err)
	}
	var dump Dump
	if err := json.Unmarshal(data, &dump); err != nil {
		t.Fatalf("parsing test data: %v", err)
	}

	el, err := ParseEventLog(dump.Log.Raw)
	if err != nil {
		t.Fatalf("parsing event log: %v", err)
	}
	events, err := el.Verify(dump.Log.PCRs)
	if err != nil {
		t.Fatalf("validating event log: %v", err)
	}

	sbState, err := ParseSecurebootState(events)
	if err != nil {
		t.Fatalf("ExtractSecurebootState() failed: %v", err)
	}

	if got, want := sbState.Enabled, true; got != want {
		t.Errorf("secureboot.Enabled = %v, want %v", got, want)
	}
}

// See: https://github.com/google/go-attestation/issues/157
func TestSecureBootBug157(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/sb_cert_eventlog")
	if err != nil {
		t.Fatalf("reading test data: %v", err)
	}
	elr, err := ParseEventLog(raw)
	if err != nil {
		t.Fatalf("parsing event log: %v", err)
	}

	pcrs := []PCR{
		{'\x00', []byte("Q\xc3#\xde\f\fiOF\x01\xcd\xd0+\xebX\xff\x13b\x9ft"), '\x03'},
		{'\x01', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x02', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x03', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x04', []byte("\xb7q\x00\x8d\x17<\x02+\xc1oKM\x1a\u007f\x8b\x99\xed\x88\xee\xb1"), '\x03'},
		{'\x05', []byte("\xd79j\xc6\xe8\x87\xda\"ޠ;@\x95/p\xb8\xdbҩ\x96"), '\x03'},
		{'\x06', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\a', []byte("E\xa8b\x1d4\xa5}\xf2\xb2\xe7\xf1L\x92\xb9\x9a\xc8\xde}X\x05"), '\x03'},
		{'\b', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\t', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\n', []byte("\x82\x84\x10>\x06\xd4\x01\"\xbcd\xa0䡉\x1a\xf9\xec\xd4\\\xf6"), '\x03'},
		{'\v', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\f', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\r', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x0e', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x0f', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x10', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x11', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x03'},
		{'\x12', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x03'},
		{'\x13', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x03'},
		{'\x14', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x03'},
		{'\x15', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x03'},
		{'\x16', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x03'},
		{'\x17', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x03'},
		{'\x00', []byte("\xfc\xec\xb5j\xcc08b\xb3\x0e\xb3Bę\v\xebP\xb5ૉr$I\xc2٧?7\xb0\x19\xfe"), '\x05'},
		{'\x01', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\x02', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\x03', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\x04', []byte("\xa9)h\x80oy_\xa3D5\xd9\xf1\x18\x13hL\xa1\xe7\x05`w\xf7\x00\xbaI\xf2o\x99b\xf8m\x89"), '\x05'},
		{'\x05', []byte("̆\x18\xb7y2\xb4\xef\xda\x12\xccX\xba\xd9>\xcdѕ\x9d\xea)\xe5\xabyE%\xa6\x19\xf5\xba\xab\xee"), '\x05'},
		{'\x06', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\a', []byte("Q\xb3\x04\x88\xc9\xe6%]\x82+\xdc\x1b ٩,2\xbd\xe6\xc3\xe7\xbc\x02\xbc\xdd2\x82^\xb5\xef\x06\x9a"), '\x05'},
		{'\b', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\t', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\n', []byte("\xc3l\x9a\xb1\x10\x9b\xa0\x8a?dX!\x18\xf8G\x1a]i[\xc9#\xa0\xa2\xbd\x04]\xb1K\x97OB9"), '\x05'},
		{'\v', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\f', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\r', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\x0e', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\x0f', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\x10', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
		{'\x11', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x05'},
		{'\x12', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x05'},
		{'\x13', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x05'},
		{'\x14', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x05'},
		{'\x15', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x05'},
		{'\x16', []byte("\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"), '\x05'},
		{'\x17', []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), '\x05'},
	}

	events, err := elr.Verify(pcrs)
	if err != nil {
		t.Errorf("failed to verify log: %v", err)
	}

	sbs, err := ParseSecurebootState(events)
	if err != nil {
		t.Errorf("failed parsing secureboot state: %v", err)
	}
	if got, want := len(sbs.PostSeparatorAuthority), 3; got != want {
		t.Errorf("len(sbs.PostSeparatorAuthority) = %d, want %d", got, want)
	}
}
