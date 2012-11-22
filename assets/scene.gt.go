package assets

//	Declares a complete, self-contained base of a Scene hierarchy or Scene graph. Currently just defined by a
//	Visual Scene, later to be augmented by optional "kinematics scenes" and/or "physics scenes".
type SceneDef struct {
	BaseDef

	//	The Visual Scene associated with this Scene.
	VisualSceneInst *VisualSceneInst
}

	func (me *SceneDef) init () {
	}

//	An instance referencing a Scene definition.
type SceneInst struct {
	BaseInst

	//	The Scene definition referenced by this instance.
	Def *SceneDef
}

	func (me *SceneInst) init () {
	}

//#begin-gt _definstlib.gt T:Scene

	func newSceneDef (id string) (me *SceneDef) {
		me = &SceneDef {}; me.BaseDef.init(id); me.init(); return
	}

	//	Creates and returns a new *SceneInst* instance referencing this *SceneDef* definition.
	func (me *SceneDef) NewInst (id string) (inst *SceneInst) {
		inst = &SceneInst { Def: me }; inst.Base.init(id); inst.init(); return
	}

var (
	//	A *map* collection that contains *LibSceneDefs* libraries associated by their *ID*.
	AllSceneDefLibs = LibsSceneDef {}

	//	The "default" *LibSceneDefs* library for *SceneDef*s.
	SceneDefs = AllSceneDefLibs.AddNew("")
)

func init () {
	syncHandlers = append(syncHandlers, func () { for _, lib := range AllSceneDefLibs { lib.SyncChanges() } })
}

//	The underlying type of the global *AllSceneDefLibs* variable: a *map* collection that contains
//	*LibSceneDefs* libraries associated by their *ID*.
type LibsSceneDef map[string]*LibSceneDefs

	//	Creates a new *LibSceneDefs* library with the specified *ID*, adds it to this *LibsSceneDef*, and returns it.
	//	
	//	If this *LibsSceneDef* already contains a *LibSceneDefs* library with the specified *ID*, does nothing and returns *nil*.
	func (me LibsSceneDef) AddNew (id string) (lib *LibSceneDefs) {
		if me[id] != nil { return }; lib = me.new(id); me[id] = lib; return
	}

	func (me LibsSceneDef) new (id string) (lib *LibSceneDefs) {
		lib = newLibSceneDefs(id); return
	}

//	A library that contains *SceneDef*s associated by their *ID*. To create a new *LibSceneDefs* library, ONLY
//	use the *LibsSceneDef.New()* or *LibsSceneDef.AddNew()* methods.
type LibSceneDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map [string] *SceneDef
}

	func newLibSceneDefs (id string) (me *LibSceneDefs) {
		me = &LibSceneDefs { M: map[string]*SceneDef {} }; me.Base.init(id); return
	}

	//	Adds the specified *SceneDef* definition to this *LibSceneDefs*, and returns it.
	//	
	//	If this *LibSceneDefs* already contains a *SceneDef* definition with the same *ID*, does nothing and returns *nil*.
	func (me *LibSceneDefs) Add (d *SceneDef) (n *SceneDef) { if me.M[d.ID] == nil { n, me.M[d.ID] = d, d; me.SetDirty() }; return }

	//	Creates a new *SceneDef* definition with the specified *ID*, adds it to this *LibSceneDefs*, and returns it.
	//	
	//	If this *LibSceneDefs* already contains a *SceneDef* definition with the specified *ID*, does nothing and returns *nil*.
	func (me *LibSceneDefs) AddNew (id string) *SceneDef { return me.Add(me.New(id)) }

	//	Creates a new *SceneDef* definition with the specified *ID* and returns it, but does not add it to this *LibSceneDefs*.
	func (me *LibSceneDefs) New (id string) (def *SceneDef) { def = newSceneDef(id); return }

	//	Removes the *SceneDef* with the specified *ID* from this *LibSceneDefs*.
	func (me *LibSceneDefs) Remove (id string) { delete(me.M, id); me.SetDirty() }

	//	Signals to *core* (or your custom package) that changes have been made to this *LibSceneDefs* that need to be picked up.
	//	Call this after you have made any number of changes to this *LibSceneDefs* library or its *SceneDef* definitions.
	//	Also called by the global *SyncChanges()* function.
	func (me *LibSceneDefs) SyncChanges () {
		me.BaseLib.Base.SyncChanges()
		for _, def := range me.M { def.BaseDef.Base.SyncChanges() }
	}

//#end-gt
