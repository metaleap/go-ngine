package core

import (
	"math"

	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type Controller struct {
	Pos, Dir, UpAxis             unum.Vec3
	MoveSpeed, MoveSpeedupFactor float64
	TurnSpeed, TurnSpeedupFactor float64
	MaxTurnUp, MinTurnDown       float64

	autoUpdate             bool
	hAngle, vAngle         float64
	posNeg, axH, axV       unum.Vec3
	mat, matTrans, matLook unum.Mat4
	glMat                  ugl.GlMat4
}

func (me *Controller) BeginUpdate() {
	me.autoUpdate = false
}

func (me *Controller) EndUpdate() {
	me.autoUpdate = true
	me.UpdateMatrixRot()
	me.UpdateMatrix()
}

func (me *Controller) init() {
	me.Dir.Z, me.UpAxis.Y = 1, 1
	me.MoveSpeed, me.MoveSpeedupFactor, me.TurnSpeed, me.TurnSpeedupFactor = 1.38, 1, 90, 1
	me.autoUpdate, me.MaxTurnUp, me.MinTurnDown = true, 90, -90
	unum.Mat4Identities(&me.mat, &me.matTrans, &me.matLook)
	me.glMat.Load(&me.mat)
	htarget := &unum.Vec3{me.Dir.X, 0, me.Dir.Z}
	htarget.Normalize()
	if htarget.Z >= 0 {
		if htarget.X >= 0 {
			me.hAngle = 360 - unum.RadToDeg(math.Asin(htarget.Z))
		} else {
			me.hAngle = 180 + unum.RadToDeg(math.Asin(htarget.Z))
		}
	} else {
		if htarget.X >= 0 {
			me.hAngle = unum.RadToDeg(math.Asin(-htarget.Z))
		} else {
			me.hAngle = 90 + unum.RadToDeg(math.Asin(-htarget.Z))
		}
	}
	me.vAngle = -unum.RadToDeg(math.Asin(me.Dir.Y))
	me.UpdateMatrixRot()
	me.UpdateMatrix()
}

func (me *Controller) MoveBackward() {
	me.Pos.SetFromAddMult1(&me.Pos, &me.Dir, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *Controller) MoveDown() {
	me.Pos.SetFromSubMult1(&me.Pos, &me.UpAxis, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *Controller) MoveForward() {
	me.Pos.SetFromSubMult1(&me.Pos, &me.Dir, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *Controller) MoveLeft() {
	me.Pos.SetFromAddMult1(&me.Pos, me.Dir.CrossNormalized(&me.UpAxis), me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *Controller) MoveRight() {
	me.Pos.SetFromAddMult1(&me.Pos, me.UpAxis.CrossNormalized(&me.Dir), me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *Controller) MoveUp() {
	me.Pos.SetFromAddMult1(&me.Pos, &me.UpAxis, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *Controller) StepSizeMove() float64 {
	return Loop.TickDelta * me.MoveSpeed * me.MoveSpeedupFactor
}

func (me *Controller) StepSizeTurn() float64 {
	return Loop.TickDelta * me.TurnSpeed * me.TurnSpeedupFactor
}

func (me *Controller) TurnDown() {
	if me.vAngle > me.MinTurnDown {
		me.vAngle -= me.StepSizeTurn()
		me.UpdateMatrixRot()
		me.UpdateMatrix()
	}
}

func (me *Controller) TurnLeft() {
	me.hAngle += me.StepSizeTurn()
	me.UpdateMatrixRot()
	me.UpdateMatrix()
}

func (me *Controller) TurnRight() {
	me.hAngle -= me.StepSizeTurn()
	me.UpdateMatrixRot()
	me.UpdateMatrix()
}

func (me *Controller) TurnUp() {
	if me.vAngle < me.MaxTurnUp {
		me.vAngle += me.StepSizeTurn()
		me.UpdateMatrixRot()
		me.UpdateMatrix()
	}
}

func (me *Controller) UpdateMatrix() {
	if me.autoUpdate {
		me.posNeg.SetFromNeg(&me.Pos)
		me.matLook.LookAt(&me.Dir, &me.UpAxis)
		me.matTrans.Translation(&me.posNeg)
		me.mat.SetFromMult4(&me.matLook, &me.matTrans)
		me.glMat.Load(&me.mat)
	}
}

func (me *Controller) UpdateMatrixRot() {
	if me.autoUpdate {
		me.axV.X, me.axV.Y, me.axV.Z = 0, 1, 0
		me.Dir.X, me.Dir.Y, me.Dir.Z = 1, 0, 0
		me.Dir.RotateDeg(me.hAngle, &me.axV)
		me.Dir.Normalize()

		me.axH = *me.axV.Cross(&me.Dir)
		me.axH.Normalize()
		me.Dir.RotateDeg(me.vAngle, &me.axH)
		me.Dir.Normalize()

		me.UpAxis = *me.Dir.Cross(&me.axH)
		me.UpAxis.Normalize()
	}
}
