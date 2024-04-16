package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (priest *Priest) registerMindSpike() {
	mbMod := priest.AddDynamicMod(core.SpellModConfig{
		ClassMask:  PriestSpellMindBlast,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 30 * core.CritRatingPerCritChance,
	})

	procAura := priest.RegisterAura(core.Aura{
		Label:     "Mind Spike Buff",
		ActionID:  core.ActionID{SpellID: 87178},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks > 0 {
				mbMod.UpdateFloatValue(float64(newStacks) * 30 * core.CritRatingPerCritChance)
				mbMod.Activate()
			} else {
				mbMod.Deactivate()
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask == PriestSpellMindBlast {
				aura.Deactivate(sim)
			}
		},
	})

	priest.MindSpike = priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 73510},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellMindSpike,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultSpellCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},
		ThreatMultiplier: 1,

		BonusCoefficient: 0.8355,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := priest.calcBaseDamage(sim, 1.178, 0.055)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Outcome.Matches(core.OutcomeLanded) {
				priest.ShadowWordPain.Dot(target).Deactivate(sim)
				priest.VampiricTouch.Dot(target).Deactivate(sim)
				priest.DevouringPlague.Dot(target).Deactivate(sim)
				procAura.Activate(sim)
				procAura.AddStack(sim)
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := priest.calcBaseDamage(sim, 1.557, 0.055)
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
