package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	nga "github.com/go3d/go-ngine/assets"
	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
)

type libMap struct {
	xnPlural, xnSingular, tn string
}

var (
	outDirPath = flag.String("dst", ugo.GopathSrcGithub("go3d", "go-ngine", "assets", "collada", "imp-1.5"), "out dir path")
	libs       = []libMap{
		libMap{"animation_clips", "animation_clip", "AnimationClip"},
		libMap{"animations", "animation", "Animation"},
		libMap{"cameras", "camera", "Camera"},
		libMap{"controllers", "controller", "Controller"},
		libMap{"formulas", "formula", "Formula"},
		libMap{"geometries", "geometry", "Geometry"},
		libMap{"lights", "light", "Light"},
		libMap{"nodes", "node", "Node"},
		libMap{"visual_scenes", "visual_scene", "VisualScene"},
		libMap{"force_fields", "force_field", "PxForceField"},
		libMap{"physics_materials", "physics_material", "PxMaterial"},
		libMap{"physics_models", "physics_model", "PxModel"},
		libMap{"physics_scenes", "physics_scene", "PxScene"},
		libMap{"effects", "effect", "FxEffect"},
		libMap{"images", "image", "FxImage"},
		libMap{"materials", "material", "FxMaterial"},
		libMap{"articulated_systems", "articulated_system", "KxArticulatedSystem"},
		libMap{"joints", "joint", "KxJoint"},
		libMap{"kinematics_models", "kinematics_model", "KxModel"},
		libMap{"kinematics_scenes", "kinematics_scene", "KxScene"},
	}
)

func main() {
	const (
		srcImpLib = `
func libs_%s(xn *xmlx.Node) {
	var (
		lib *nga.Lib%sDefs
		def *nga.%sDef
		id  string
	)
	for _, ln := range xcns(xn, "library_%s") {
		id = xas(ln, "id")
		if lib = nga.All%sDefLibs[id]; lib == nil {
			lib = nga.All%sDefLibs.AddNew(id)
		}
		for _, def = range objs_%sDef(ln, "%s") {
			if def != nil {
				lib.Add(def)
			}
		}
		lib.SetDirty()
	}
}
`
		srcImpObj = `
func obj_%s(xn *xmlx.Node, n string) (obj *nga.%s) {
	if (xn != nil) && (len(n) > 0) {
		xn = xcn(xn, n)
	}
	if xn != nil {
		obj = init_%s(xn)
	}
	return
}
`
		srcImpInitCtor = `
func init_%s(xn *xmlx.Node) (obj *nga.%s) {
	obj = nga.New%s()
`
		srcImpInitNew = `
func init_%s(xn *xmlx.Node) (obj *nga.%s) {
	obj = new(nga.%s)
`
		srcImpN = `
func objs_%s(xn *xmlx.Node, n string) (objs []*nga.%s) {
	xns := xcns(xn, n)
	objs = make([]*nga.%s, len(xns))
	for i, xn := range xns {
		objs[i] = obj_%s(xn, "")
	}
	return
}
`
		srcLoad = `
func load_%s(xn *xmlx.Node, obj *nga.%s) {

}
`
	)
	var (
		i                                    int
		ok, canDirty                         bool
		srcLibs, srcInits, srcObjs, srcLoads string
	)
	has := []string{"Asset", "Extras", "FxParamDefs", "ID", "Inputs", "Name", "ParamDefs", "ParamInsts", "Sid", "Sources", "Techniques"}
	flag.Parse()
	for n, t := range nga.Types {
		if canDirty = false; !(strings.HasPrefix(n, "Lib") || strings.HasPrefix(n, "Mesh") || strings.HasPrefix(n, "Base") || strings.HasSuffix(n, "Base") || strings.HasPrefix(n, "Has") || strings.HasPrefix(n, "Ref")) {
			srcObjs += fmt.Sprintf(srcImpObj, n, n, n)
			if _, ok = nga.Functions["New"+n]; ok && (n != "FxImageInitFrom") {
				srcInits += fmt.Sprintf(srcImpInitCtor, n, n, n)
			} else {
				srcInits += fmt.Sprintf(srcImpInitNew, n, n, n)
			}
			if t.Kind() == reflect.Struct {
				for i = 0; i < t.NumField(); i++ {
					if strings.HasPrefix(t.Field(i).Name, "Base") {
						canDirty = true
					}
				}
				if canDirty {
					if strings.HasSuffix(n, "Def") || strings.HasSuffix(n, "Inst") {
						srcInits += "\tobj.Init()\n"
					}
					if strings.HasSuffix(n, "Inst") {
						srcInits += "\tsetInstDefRef(xn, &obj.BaseInst)\n"
					}
				}
				for _, h := range has {
					if _, ok = t.FieldByName("Has" + h); ok {
						srcInits += fmt.Sprintf("\thas_%s(xn, &obj.Has%s)\n", h, h)
					}
				}
			}
			srcInits += fmt.Sprintf("\n\tload_%s(xn, obj)", n)
			if canDirty {
				srcInits += "\n\tobj.SetDirty()"
			}
			srcInits += "\n\treturn\n}\n"
			srcObjs += fmt.Sprintf(srcImpN, n, n, n, n)
			srcLoads += fmt.Sprintf(srcLoad, n, n)
		}
	}
	for _, lm := range libs {
		//	animations Animation Animation animations Animation Animation Animation animation
		srcLibs += fmt.Sprintf(srcImpLib, lm.xnPlural, lm.tn, lm.tn, lm.xnPlural, lm.tn, lm.tn, lm.tn, lm.xnSingular)
	}
	srcLibs += "\nfunc libs_All (xn *xmlx.Node) {\n"
	for _, lm := range libs {
		srcLibs += fmt.Sprintf("\tlibs_%s(xn)\n", lm.xnPlural)
	}
	srcLibs += "}\n"
	uio.WriteTextFile(filepath.Join(*outDirPath, "skel_libs.txt"), srcLibs)
	uio.WriteTextFile(filepath.Join(*outDirPath, "skel_inits.txt"), srcInits)
	uio.WriteTextFile(filepath.Join(*outDirPath, "skel_objs.txt"), srcObjs)
	uio.WriteTextFile(filepath.Join(*outDirPath, "skel_loads.txt"), srcLoads)
}
