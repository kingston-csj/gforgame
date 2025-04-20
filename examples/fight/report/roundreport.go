package report

type RoundReport struct {
	Round        int32         `json:"round"`
	SkillReports []SkillReport `json:"skill_reports"`
}

func NewRoundReport(round int32) *RoundReport {
	return &RoundReport{Round: round, SkillReports: make([]SkillReport, 0)}
}

func (r *RoundReport) AddSkillReport(skillReport *SkillReport) {
	r.SkillReports = append(r.SkillReports, *skillReport)
}
