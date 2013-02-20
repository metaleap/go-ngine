package core

type RenderBatchCriteria int

const (
	_ RenderBatchCriteria = iota
	BatchByProgram
	BatchByTexture
	BatchByBuffer
)

type nodeBatchPrep map[string][]*Node

func (me nodeBatchPrep) add(key string, node *Node) {
	me[key] = append(me[key], node)
}

func (me *nodeBatchPrep) remake() {
	*me = make(nodeBatchPrep, len(*me))
}

type nodeBatchPreps struct {
	allNodes map[*Node]bool
	byBuf    nodeBatchPrep
	byFx     nodeBatchPrep
	byImg    nodeBatchPrep

	tmat  *FxMaterial
	tfx   *FxEffect
	tfxOp FxOp
	fxId  string
	allFx map[*FxEffect]bool
}

func (me *nodeBatchPreps) prep() {
}

func (me *nodeBatchPreps) prepNodeFx(node *Node) {
	if me.tfx = Core.Libs.Effects[me.fxId]; !me.allFx[me.tfx] {
		me.allFx[me.tfx] = true
	}
}

func (me *nodeBatchPreps) prepNode(node *Node) {
	if !me.allNodes[node] {
		me.allNodes[node] = true
		if node.mesh != nil && node.mesh.meshBuffer != nil {
			me.byBuf.add(node.mesh.meshBuffer.id, node)
			if me.tmat = node.EffectiveMaterial(); me.tmat != nil {
				me.allFx = make(map[*FxEffect]bool, 1+len(me.tmat.FaceEffects.ByID)+len(me.tmat.FaceEffects.ByTag))
				me.fxId = me.tmat.DefaultEffectID
				me.prepNodeFx(node)
				for _, me.fxId = range me.tmat.FaceEffects.ByID {
					me.prepNodeFx(node)
				}
				for _, me.fxId = range me.tmat.FaceEffects.ByTag {
					me.prepNodeFx(node)
				}
				for me.tfx, _ = range me.allFx {
					me.byFx.add(me.tfx.uberName, node)
					for _, me.tfxOp = range me.tfx.Ops {
						for _, me.fxId = range me.tfxOp.fxImageIDs() {
							me.byImg.add(me.fxId, node)
						}
					}
				}
				me.allFx = nil
			}
		}
	}
}

func (me *nodeBatchPreps) reset() {
	me.allNodes, me.fxId, me.tmat, me.tfx, me.tfxOp = make(map[*Node]bool, len(me.allNodes)), "", nil, nil, nil
	me.byBuf.remake()
	me.byFx.remake()
	me.byImg.remake()
}

type RenderBatch struct {
	Enabled  bool
	Priority [3]RenderBatchCriteria

	cam  *Camera
	tech *RenderTechniqueScene
}

func (me *RenderBatch) init(tech *RenderTechniqueScene) {
	me.cam, me.tech = tech.cam, tech
	me.Priority[0] = BatchByProgram
	me.Priority[1] = BatchByTexture
	me.Priority[2] = BatchByBuffer
	return
}

func (me *RenderBatch) onPrep() {
	thrPrep.nodePreBatch.prep()
}
