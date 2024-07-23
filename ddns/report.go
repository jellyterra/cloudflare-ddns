package ddns

type RecordUpdateReport struct {
	Record *Record

	Err error
}

type ZoneUpdateReport struct {
	Zone *Zone

	RecordUpdateReports []*RecordUpdateReport
}

type EnvUpdateReport struct {
	ZoneUpdateReports []*ZoneUpdateReport
}
