package core

import (
	"math"

	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type TController struct {
	Pos, Dir, UpAxis *unum.Vec3
	MoveSpeed, MoveSpeedupFactor float64
	TurnSpeed, TurnSpeedupFactor float64
	MaxTurnUp, MinTurnDown float64

	autoUpdate bool
	hAngle, vAngle float64
	posNeg, axH, axV *unum.Vec3
	mat, matTrans, matLook *unum.Mat4
	glMat *ugl.GlMat4
}

func newController () *TController {
	var htarget *unum.Vec3
	var ctl = &TController {}
	ctl.posNeg, ctl.Pos, ctl.Dir, ctl.UpAxis, ctl.axH, ctl.axV = &unum.Vec3 {}, &unum.Vec3 {}, &unum.Vec3 { 0, 0, 1 }, &unum.Vec3 { 0, 1, 0 }, &unum.Vec3 {}, &unum.Vec3 {}
	ctl.MoveSpeed, ctl.MoveSpeedupFactor, ctl.TurnSpeed, ctl.TurnSpeedupFactor = 1.38, 1, 90, 1
	ctl.autoUpdate, ctl.MaxTurnUp, ctl.MinTurnDown = true, 90, -90
	ctl.mat, ctl.matTrans, ctl.matLook = unum.NewMat4Identity(), unum.NewMat4Identity(), unum.NewMat4Identity()
	ctl.glMat = ugl.NewGlMat4(ctl.mat)

	htarget = &unum.Vec3 { ctl.Dir.X, 0, ctl.Dir.Z }
	htarget.Normalize()
	if htarget.Z >= 0 {
		if htarget.X >= 0 {
			ctl.hAngle = 360 - unum.RadToDeg(math.Asin(htarget.Z))
		} else {
			ctl.hAngle = 180 + unum.RadToDeg(math.Asin(htarget.Z))
		}
	} else {
		if htarget.X >= 0 {
			ctl.hAngle = unum.RadToDeg(math.Asin(-htarget.Z))
		} else {
			ctl.hAngle = 90 + unum.RadToDeg(math.Asin(-htarget.Z))
		}
	}
	ctl.vAngle = -unum.RadToDeg(math.Asin(ctl.Dir.Y))
	ctl.updateMatrixRot()
	ctl.updateMatrix()
	return ctl
}

func (me *TController) BeginUpdate () {
	me.autoUpdate = false
}

func (me *TController) EndUpdate () {
	me.autoUpdate = true
	me.updateMatrixRot()
	me.updateMatrix()
}

func (me *TController) MoveBackward () {
	me.Pos.SetFromAddMult1(me.Pos, me.Dir, me.StepSizeMove())
	me.updateMatrix()
}

func (me *TController) MoveDown () {
	me.Pos.SetFromSubMult1(me.Pos, me.UpAxis, me.StepSizeMove())
	me.updateMatrix()
}

func (me *TController) MoveForward () {
	me.Pos.SetFromSubMult1(me.Pos, me.Dir, me.StepSizeMove())
	me.updateMatrix()
}

func (me *TController) MoveLeft () {
	me.Pos.SetFromAddMult1(me.Pos, me.Dir.CrossNormalized(me.UpAxis), me.StepSizeMove())
	me.updateMatrix()
}

func (me *TController) MoveRight () {
	me.Pos.SetFromAddMult1(me.Pos, me.UpAxis.CrossNormalized(me.Dir), me.StepSizeMove())
	me.updateMatrix()
}

func (me *TController) MoveUp () {
	me.Pos.SetFromAddMult1(me.Pos, me.UpAxis, me.StepSizeMove())
	me.updateMatrix()
}

func (me *TController) StepSizeMove () float64 {
	return Loop.TickDelta * me.MoveSpeed * me.MoveSpeedupFactor
}

func (me *TController) StepSizeTurn () float64 {
	return Loop.TickDelta * me.TurnSpeed * me.TurnSpeedupFactor
}

func (me *TController) TurnDown () {
	if me.vAngle > me.MinTurnDown {
		me.vAngle -= me.StepSizeTurn()
		me.updateMatrixRot()
		me.updateMatrix()
	}
}

func (me *TController) TurnLeft () {
	me.hAngle += me.StepSizeTurn()
	me.updateMatrixRot()
	me.updateMatrix()
}

func (me *TController) TurnRight () {
	me.hAngle -= me.StepSizeTurn()
	me.updateMatrixRot()
	me.updateMatrix()
}

func (me *TController) TurnUp () {
	if me.vAngle < me.MaxTurnUp {
		me.vAngle += me.StepSizeTurn()
		me.updateMatrixRot()
		me.updateMatrix()
	}
}

func (me *TController) updateMatrix () {
	if me.autoUpdate {
		me.posNeg.SetFromNeg(me.Pos)
		me.matLook.LookAt(me.Dir, me.UpAxis)
		me.matTrans.Translation(me.posNeg)
		me.mat.SetFromMult4(me.matLook, me.matTrans)
		me.glMat.Load(me.mat)
	}
}

func (me *TController) updateMatrixRot () {
	if me.autoUpdate {
		me.axV.X, me.axV.Y, me.axV.Z = 0, 1, 0
		me.Dir.X, me.Dir.Y, me.Dir.Z = 1, 0, 0
		me.Dir.RotateDeg(me.hAngle, me.axV)
		me.Dir.Normalize()

		me.axH = me.axV.Cross(me.Dir)
		me.axH.Normalize()
		me.Dir.RotateDeg(me.vAngle, me.axH)
		me.Dir.Normalize()

		me.UpAxis = me.Dir.Cross(me.axH)
		me.UpAxis.Normalize()
	}
}
