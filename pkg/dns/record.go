package dns

import "strings"

const DefaultRecordTTL = 300
const DefaultRecordPriority = 10

type Record struct {
	Name     string
	Data     string
	TTL      int
	Priority int
}

func (r *Record) DotCount() int {
	return strings.Count(r.Name, ".")
}

func (r *Record) HasSub() bool {
	return r.DotCount() > 1
}

func (r *Record) GetSub() string {
	if !r.HasSub() {
		return "@"
	}
	return strings.Join(strings.Split(r.Name, ".")[:r.DotCount()-1], ".")
}

func (r *Record) GetTld() string {
	if !r.HasSub() {
		return r.Name
	}
	return strings.Join(strings.Split(r.Name, ".")[r.DotCount()-1:], ".")
}

func NewRecord(name, data string) *Record {
	return &Record{
		Name:     name,
		Data:     data,
		TTL:      DefaultRecordTTL,
		Priority: DefaultRecordPriority,
	}
}
