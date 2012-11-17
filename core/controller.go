package core

import (
	ugl "github.com/go3d/go-glutil"

	nga "github.com/go3d/go-ngine/assets"
)

type Controller struct {
	nga.OldController

	glMat *ugl.GlMat4
}

func newController () (me *Controller) {
	me = &Controller {}
	me.glMat = &ugl.GlMat4 {}
	me.OldController.Init(func () { me.glMat.Load(me.Mat) }, func () float64 { return Loop.TickDelta })
	return me
}
