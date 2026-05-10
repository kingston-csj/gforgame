package report

import "fmt"

type BattleReport struct {
	RoundReports []RoundReport `json:"round_reports"`
	Winner       int32         `json:"winner"`
}

func NewBattleReport() *BattleReport {
	return &BattleReport{RoundReports: make([]RoundReport, 0)}
}

func (r *BattleReport) AddRoundReport(roundReport *RoundReport) {
	r.RoundReports = append(r.RoundReports, *roundReport)
}

func (r *BattleReport) Display() {
	for _, roundReport := range r.RoundReports {
		fmt.Println("第", roundReport.Round, "回合开始")
		for _, skillReport := range roundReport.SkillReports {
			fmt.Println(skillReport.Attacker+"释放技能", skillReport.SkillId)
			for _, damageUnit := range skillReport.DamageUnits {
				fmt.Println(damageUnit.attacker+"对", damageUnit.defender+"造成", damageUnit.damage, "点伤害")
			}
		}
		fmt.Println("第", roundReport.Round, "回合结束")
	}
	fmt.Println("战斗结束, 胜利者是", r.Winner)
}
