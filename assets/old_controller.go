package assets

import (
	"math"

	unum "github.com/metaleap/go-util/num"
)

type OldController struct {
	Pos, Dir, UpAxis *unum.Vec3
	MoveSpeed, MoveSpeedupFactor float64
	TurnSpeed, TurnSpeedupFactor float64
	MaxTurnUp, MinTurnDown float64
	TickDelta func () float64
	OnUpdateMatrix func ()

	autoUpdate bool
	hAngle, vAngle float64
	posNeg, axH, axV *unum.Vec3
	Mat, MatTrans, MatLook *unum.Mat4
}

func (me *OldController) Init (onUpdateMatrix func (), tickDelta func () float64) {
	var htarget *unum.Vec3
	me.OnUpdateMatrix, me.TickDelta = onUpdateMatrix, tickDelta
	me.posNeg, me.Pos, me.Dir, me.UpAxis, me.axH, me.axV = &unum.Vec3 {}, &unum.Vec3 {}, &unum.Vec3 { 0, 0, 1 }, &unum.Vec3 { 0, 1, 0 }, &unum.Vec3 {}, &unum.Vec3 {}
	me.MoveSpeed, me.MoveSpeedupFactor, me.TurnSpeed, me.TurnSpeedupFactor = 1.38, 1, 90, 1
	me.autoUpdate, me.MaxTurnUp, me.MinTurnDown = true, 90, -90
	me.Mat, me.MatTrans, me.MatLook = unum.NewMat4Identity(), unum.NewMat4Identity(), unum.NewMat4Identity()

	htarget = &unum.Vec3 { me.Dir.X, 0, me.Dir.Z }
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

func (me *OldController) BeginUpdate () {
	me.autoUpdate = false
}

func (me *OldController) EndUpdate () {
	me.autoUpdate = true
	me.UpdateMatrixRot()
	me.UpdateMatrix()
}

func (me *OldController) MoveBackward () {
	me.Pos.SetFromAddMult1(me.Pos, me.Dir, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *OldController) MoveDown () {
	me.Pos.SetFromSubMult1(me.Pos, me.UpAxis, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *OldController) MoveForward () {
	me.Pos.SetFromSubMult1(me.Pos, me.Dir, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *OldController) MoveLeft () {
	me.Pos.SetFromAddMult1(me.Pos, me.Dir.CrossNormalized(me.UpAxis), me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *OldController) MoveRight () {
	me.Pos.SetFromAddMult1(me.Pos, me.UpAxis.CrossNormalized(me.Dir), me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *OldController) MoveUp () {
	me.Pos.SetFromAddMult1(me.Pos, me.UpAxis, me.StepSizeMove())
	me.UpdateMatrix()
}

func (me *OldController) StepSizeMove () float64 {
	return me.TickDelta() * me.MoveSpeed * me.MoveSpeedupFactor
}

func (me *OldController) StepSizeTurn () float64 {
	return me.TickDelta() * me.TurnSpeed * me.TurnSpeedupFactor
}

func (me *OldController) TurnDown () {
	if me.vAngle > me.MinTurnDown {
		me.vAngle -= me.StepSizeTurn()
		me.UpdateMatrixRot()
		me.UpdateMatrix()
	}
}

func (me *OldController) TurnLeft () {
	me.hAngle += me.StepSizeTurn()
	me.UpdateMatrixRot()
	me.UpdateMatrix()
}

func (me *OldController) TurnRight () {
	me.hAngle -= me.StepSizeTurn()
	me.UpdateMatrixRot()
	me.UpdateMatrix()
}

func (me *OldController) TurnUp () {
	if me.vAngle < me.MaxTurnUp {
		me.vAngle += me.StepSizeTurn()
		me.UpdateMatrixRot()
		me.UpdateMatrix()
	}
}

func (me *OldController) UpdateMatrix () {
	if me.autoUpdate {
		me.posNeg.SetFromNeg(me.Pos)
		me.MatLook.LookAt(me.Dir, me.UpAxis)
		me.MatTrans.Translation(me.posNeg)
		me.Mat.SetFromMult4(me.MatLook, me.MatTrans)
		me.OnUpdateMatrix()
	}
}

func (me *OldController) UpdateMatrixRot () {
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
