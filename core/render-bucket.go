package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type renderBucket interface {
	render()
}

type renderBucketBase struct {
	batch      *RenderBatch
	subBuckets []renderBucket
}

func (me *renderBucketBase) render() {
	// for _, me.batch.thrRend.tmpBucket = range me.subBuckets {
	// 	me.batch.thrRend.tmpBucket.render()
	// }
}

type renderBucketBuffer struct {
	renderBucketBase
	meshBuf *MeshBuffer
}

func (me *renderBucketBuffer) render() {
	if thrRend.curMeshBuf != me.meshBuf {
		me.meshBuf.use()
	}
	me.renderBucketBase.render()
}

type renderBucketProgram struct {
	renderBucketBase
	fx *FxEffect
}

func (me *renderBucketProgram) render() {
	thrRend.nextTech, thrRend.nextEffect = me.batch.tech, me.fx
	Core.useTechFx()
	me.renderBucketBase.render()
}

type renderBucketTexture struct {
	renderBucketBase
	glTex ugl.Texture
}

func (me *renderBucketTexture) render() {
	me.glTex.Bind()
	me.renderBucketBase.render()
}
