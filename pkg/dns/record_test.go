package dns_test

import (
	"testing"

	"github.com/lajosbencz/dogodns/pkg/dns"
)

const testRecordNameTld = "dev.localhost"
const testRecordNameSub = "dogodns"
const testRecordNameFull = testRecordNameSub + "." + testRecordNameTld
const testRecordData = "127.0.0.1"

var testRecord *dns.Record = dns.NewRecord(testRecordNameFull, testRecordData)

func TestRecord(T *testing.T) {

	// Name check
	if testRecord.Name != testRecordNameFull {
		T.Error()
	}
	// Data check
	if testRecord.Data != testRecordData {
		T.Error(testRecord.Data)
	}
	// TTL check
	if testRecord.TTL != dns.DefaultRecordTTL {
		T.Error(testRecord.TTL)
	}
	// Priority check
	if testRecord.Priority != dns.DefaultRecordPriority {
		T.Error(testRecord.Priority)
	}
	// DotCount() should be 2
	if testRecord.DotCount() != 2 {
		T.Error(testRecord.DotCount())
	}
	// HasSub() should be true
	if testRecord.HasSub() == false {
		T.Error(testRecord.HasSub())
	}
	// GetTld() check
	if testRecord.GetTld() != testRecordNameTld {
		T.Error(testRecord.GetTld())
	}
	// GetSub() check
	if testRecord.GetSub() != testRecordNameSub {
		T.Error(testRecord.GetSub())
	}

	// Change record to TLD only
	testRecord.Name = testRecordNameTld

	// DotCount() should be 1
	if testRecord.DotCount() != 1 {
		T.Error(testRecord.DotCount())
	}
	// HasSub() should be false
	if testRecord.HasSub() == true {
		T.Error(testRecord.HasSub())
	}
	// GetTld() check
	if testRecord.GetTld() != testRecordNameTld {
		T.Error(testRecord.GetTld())
	}
	// GetSub() should be empty string
	if testRecord.GetSub() != "@" {
		T.Error(testRecord.GetSub())
	}
}
