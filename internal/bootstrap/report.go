package bootstrap

// AppliedSource describes a source that was applied to the target project.
type AppliedSource struct {
	Path string
	Mode Mode
}

// FileStatus describes the result of comparing a target file to its expected
// content.
type FileStatus string

const (
	StatusOK      FileStatus = "ok"
	StatusDiff    FileStatus = "diff"
	StatusMissing FileStatus = "missing"
)

// FileReport describes a single file's status.
type FileReport struct {
	Path   string
	Status FileStatus
}

// VerifyResult aggregates the reports produced by a Verify run and provides a
// quick summary of the project health.
type VerifyResult struct {
	Reports []FileReport
	OK      int
	Diff    int
	Missing int
}

// HasIssues reports whether the verification found any missing or different
// files.
func (r *VerifyResult) HasIssues() bool {
	return r.Diff > 0 || r.Missing > 0
}
