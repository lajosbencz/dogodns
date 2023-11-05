package dns

type Updater interface {
	Pull(record *Record) error
	Get(record *Record) error
	Set(record *Record) error
	Has(record *Record) bool
}
