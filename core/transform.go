package core

import (
	ugl "github.com/go3d/go-glutil"

	nga "github.com/go3d/go-ngine/assets"
)

type transform struct {
	nga.OldTransforms

	glMatModelView *ugl.GlMat4
	glMatNormal *ugl.GlMat3
}

func newTransform (owner nga.Transformable) (me *transform) {
	me = &transform {}
	me.OldTransforms.Init(owner, func () { me.glMatModelView.Load(me.MatModelView) })
	me.glMatModelView = ugl.NewGlMat4(me.MatModelView)
	return
}

func (me *transform) StepDelta (deltaPerSecond float64) float64 {
	return Loop.TickDelta * deltaPerSecond
}
