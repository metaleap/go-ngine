package assets

//	Used for the ImageDef.InitFrom field.
type ImageInitFrom struct {
	//	Initializes higher MIP levels if data does not exist in a file. Defaults to true.
	AutoMip bool

	//	Contains the embedded image data as a sequence of bytes.
	RawData []byte

	//	The URL of a file from which to take initialization data. Can be a relative path
	//	such as "walltextures/wall01.jpg".
	RefUrl string
}

//	Declares the storage for the graphical representation of an object.
type ImageDef struct {
	BaseDef
	HasSid
	//	Initializes an entire Image or portions of an Image from referenced or embedded data.
	InitFrom *ImageInitFrom
}

func (me *ImageDef) init() {
	me.InitFrom = &ImageInitFrom{AutoMip: true}
}

//	An instance referencing an Image definition.
type ImageInst struct {
	BaseInst

	//	The image definition referenced by this instance.
	Def *ImageDef
}

func (me *ImageInst) init() {
}

//#begin-gt _definstlib.gt T:Image

func newImageDef(id string) (me *ImageDef) {
	me = &ImageDef{}
	me.ID = id
	me.init()
	return
}

//	Creates and returns a new *ImageInst* instance referencing this *ImageDef* definition.
func (me *ImageDef) NewInst(id string) (inst *ImageInst) {
	inst = &ImageInst{Def: me}
	inst.init()
	return
}

var (
	//	A *map* collection that contains *LibImageDefs* libraries associated by their *ID*.
	AllImageDefLibs = LibsImageDef{}

	//	The "default" *LibImageDefs* library for *ImageDef*s.
	ImageDefs = AllImageDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllImageDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllImageDefLibs* variable: a *map* collection that contains
//	*LibImageDefs* libraries associated by their *ID*.
type LibsImageDef map[string]*LibImageDefs

//	Creates a new *LibImageDefs* library with the specified *ID*, adds it to this *LibsImageDef*, and returns it.
//	
//	If this *LibsImageDef* already contains a *LibImageDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsImageDef) AddNew(id string) (lib *LibImageDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsImageDef) new(id string) (lib *LibImageDefs) {
	lib = newLibImageDefs(id)
	return
}

//	A library that contains *ImageDef*s associated by their *ID*. To create a new *LibImageDefs* library, ONLY
//	use the *LibsImageDef.New()* or *LibsImageDef.AddNew()* methods.
type LibImageDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*ImageDef
}

func newLibImageDefs(id string) (me *LibImageDefs) {
	me = &LibImageDefs{M: map[string]*ImageDef{}}
	me.ID = id
	return
}

//	Adds the specified *ImageDef* definition to this *LibImageDefs*, and returns it.
//	
//	If this *LibImageDefs* already contains a *ImageDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibImageDefs) Add(d *ImageDef) (n *ImageDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *ImageDef* definition with the specified *ID*, adds it to this *LibImageDefs*, and returns it.
//	
//	If this *LibImageDefs* already contains a *ImageDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibImageDefs) AddNew(id string) *ImageDef { return me.Add(me.New(id)) }

//	Creates a new *ImageDef* definition with the specified *ID* and returns it, but does not add it to this *LibImageDefs*.
func (me *LibImageDefs) New(id string) (def *ImageDef) { def = newImageDef(id); return }

//	Removes the *ImageDef* with the specified *ID* from this *LibImageDefs*.
func (me *LibImageDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibImageDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibImageDefs* library or its *ImageDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibImageDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
