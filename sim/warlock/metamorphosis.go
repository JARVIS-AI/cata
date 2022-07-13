package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warlock *Warlock) registerMetamorphosisSpell() {

	warlock.MetamorphosisAura = warlock.RegisterAura(core.Aura{
		Label:    "Metamorphosis Aura",
		ActionID: core.ActionID{SpellID: 47241},
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})

	warlock.Metamorphosis = warlock.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 47241},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 3 * time.Second * (60.0 - 6.0*time.Duration(warlock.Talents.Nemesis)),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.MetamorphosisAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: warlock.Metamorphosis,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			MetamorphosisNumber := (float64(sim.Duration) + float64(warlock.MetamorphosisAura.Duration)) / float64(warlock.Metamorphosis.CD.Duration)
			if MetamorphosisNumber < 1 {
				if character.HasActiveAuraWithTag(core.BloodlustAuraTag) || sim.IsExecutePhase35() {
					return true
				}
			} else if warlock.Metamorphosis.CD.IsReady(sim) {
				return true
			}
			return false
		},
	})
}
