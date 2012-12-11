package convert

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"
	util "github.com/metaleap/go-util"
	ustr "github.com/metaleap/go-util/str"
)

const ns = "http://www.collada.org/2005/11/COLLADASchema"

var (
	//	If true, conversion-logic is always run against the given input document's entire node tree;
	//	if false, conversion is only performed if the input "COLLADA" root element's "version" attribute is not "1.5" or higher.
	Force = false

	//	Set this to an image file format (BMP, JPG, PNG etc.) for Collada 1.4.1 documents
	//	that have binary images hex-encoded in <image><data> elements *without* a format attribute.
	HexFormat = ""

	//	The default logging function used by this package. Simply does a log.Printf(format, fmtArgs...)
	//	Set this to nil to disable logging.
	Log = func(format string, fmtArgs ...interface{}) {
		log.Printf(format, fmtArgs...)
	}

	//	The <shader>/<compiler> element introduced in Collada 1.5 requires a "platform" attribute, but there was no
	//	such equivalent in Collada 1.4.1. In the context of importing such 1.4.1 documents, probably any arbitrary
	//	value will do for this new attribute, but this field still gives you the choice to set this to a preferred
	//	value for all your <shader>s.
	ShaderCompilerPlatform = "PC"

	//	Conversion mode: true for strict mode, or false for lax mode.
	//	If strict, obsoleted elements and attributes are removed or rewritten so that the remaining document is (in theory) strictly 1.5-conformant and will (ideally) validate against the Collada 1.5 XML schema definition.
	//	If lax, obsoleted elements and attributes are not removed, for use-cases where your 1.5 loader consuming the conversion result is known to simply ignore or discard them quietly.
	//	Note in practice there won't be any noticeable difference in performance or output for approximately 95% of "common use-case" Collada documents...
	Strict = true

	skipped       = false
	srcDoc        *xmlx.Document
	surfaceNodes  []*xmlx.Node
	surfaceImages map[string]string
)

func logFmt(format string, fmtArgs ...interface{}) {
	if Log != nil {
		Log(format, fmtArgs...)
	}
}

func attVal(xn *xmlx.Node, name string) (val string) {
	if val = xn.As(ns, name); len(val) == 0 {
		if val = xn.As("", name); len(val) == 0 {
			val = xn.As("*", name)
		}
	}
	return
}

func attValU64(xn *xmlx.Node, name string) (val uint64) {
	if val = xn.Au64(ns, name); val == 0 {
		if val = xn.Au64("", name); val == 0 {
			val = xn.Au64("*", name)
		}
	}
	return
}

func convert(srcFile []byte, retDoc, retBytes bool) (doc *xmlx.Document, dstFile []byte, err error) {
	doc, surfaceImages, surfaceNodes, skipped = xmlx.New(), map[string]string{}, nil, false
	if err = doc.LoadBytes(srcFile, nil); err != nil {
		doc = nil
		return
	}
	srcDoc = doc
	processNode(doc.Root)
	srcDoc = nil
	if skipped {
		if retBytes {
			dstFile = srcFile
		}
	} else {
		for _, sn := range surfaceNodes {
			delNode(sn.Parent)
		}
		if retBytes {
			dstFile = doc.SaveBytes()
		}
	}
	if !retDoc {
		doc = nil
	}
	surfaceImages, surfaceNodes = nil, nil
	return
}

//	Converts the specified Collada 1.4.1 document to Collada 1.5.
func ConvertBytes(srcFile []byte) (dstFile []byte, err error) {
	_, dstFile, err = convert(srcFile, false, true)
	return
}

//	Converts the specified Collada 1.4.1 document to Collada 1.5.
func ConvertDoc(srcFile []byte) (doc *xmlx.Document, err error) {
	doc, _, err = convert(srcFile, true, false)
	return
}

func convertImage(xn *xmlx.Node) {
	hexFormat := attVal(xn, "format")
	imgHeight := attValU64(xn, "height")
	imgWidth := attValU64(xn, "width")
	imgDepth := attValU64(xn, "depth")
	if imgDepth == 0 {
		imgDepth = 1
	}
	delAtts(xn, "height", "width", "depth", "format")
	if len(hexFormat) == 0 {
		hexFormat = HexFormat
	}
	hexData, refUrl, initNode, hexNode := "", "", subNode(xn, "init_from"), subNode(xn, "data")
	if initNode != nil {
		refUrl = initNode.Value
		delNodeForce(initNode, true)
		initNode = nil
	}
	if hexNode != nil {
		hexData = hexNode.Value
		delNodeForce(hexNode, true)
		hexNode = nil
	}
	if (len(refUrl) > 0) || (len(hexData) > 0) {
		initNode = xmlx.NewNode(xn.Type)
		initNode.Name.Local = "init_from"
		if len(refUrl) > 0 {
			ensureChild(initNode, "ref").Value = refUrl
		}
		if len(hexData) > 0 {
			hex := ensureChild(initNode, "hex")
			setAttr(hex, "format", hexFormat, false)
			hex.Value = hexData
			hexData = ""
		}
	}
	if (imgWidth != 0) && (imgHeight != 0) {
		cn := xmlx.NewNode(xn.Type)
		switch imgDepth {
		case 1:
			cn.Name.Local = "create_2d"
		case 2:
			cn.Name.Local = "create_3d"
		default:
			cn.Name.Local = "create_cube"
		}
		sn := ensureChild(cn, "size_exact")
		setAttr(sn, "width", fmt.Sprintf("%v", imgWidth), false)
		setAttr(sn, "height", fmt.Sprintf("%v", imgHeight), false)
		if initNode != nil {
			cn.AddChild(initNode)
		}
		xn.AddChild(cn)
	} else if initNode != nil {
		xn.AddChild(initNode)
	}
	if oldParent := xn.Parent; oldParent.Name.Local != "library_images" {
		logFmt("!!MOVE image!!; this may be BUGGY, please report your use-case at GitHub Issues for this package!!\n")
		id := attVal(xn, "id")
		if len(id) == 0 {
			id = fmt.Sprintf("img_moved_%v", time.Now().UnixNano())
		}
		if _, pos := moveNode(xn, nil, "library_images"); pos >= 0 {
			xn = xmlx.NewNode(xn.Type)
			xn.Name.Local = "instance_image"
			setAttr(xn, "url", "#"+id, false)
			xn.Parent = oldParent
			oldParent.Children[pos] = xn
		}
	}
}

func convertShader(xn *xmlx.Node) {
	var tn, sn *xmlx.Node
	if av, tmp := attVal(xn, "stage"), "PROGRAM"; strings.HasSuffix(av, tmp) {
		setAttr(xn, "stage", strings.Replace(av, tmp, "", -1), false)
	}
	for _, tn = range subNodes(xn, "bind") {
		renameNode(tn, "bind_uniform")
	}
	for _, tn = range subNodes(xn, "annotate") {
		delNode(tn)
	}
	if tn = subNode(xn, "compiler_target"); tn != nil {
		sn = ensureChild(xn, "compiler")
		setAttr(sn, "target", tn.Value, false)
		setAttr(sn, "platform", ShaderCompilerPlatform, false)
		delNode(tn)
	}
	if tn = subNode(xn, "compiler_options"); tn != nil {
		sn = ensureChild(xn, "compiler")
		setAttr(sn, "options", tn.Value, false)
		setAttr(sn, "platform", ShaderCompilerPlatform, false)
		delNode(tn)
	}
	if tn = subNode(xn, "name"); tn != nil {
		sn = ensureChild(xn, "sources")
		setAttr(sn, "entry", tn.Value, false)
		if src := attVal(tn, "source"); len(src) > 0 {
			setAttr(ensureChild(sn, "import"), "ref", src, false)
		}
		delNode(tn)
	}
}

func convertSurface(xn *xmlx.Node) {
	var (
		imgNode, imgCreateNode, imgCreateFormatNode, rn, sn, tn *xmlx.Node
		ensureCreateNode                                        = func() *xmlx.Node {
			if imgCreateNode == nil {
				if rn = subNode(imgNode, "init_from"); rn != nil {
					imgNode.RemoveChild(rn)
				}
				imgCreateNode = ensureChild(imgNode, "create_"+strings.ToLower(attVal(xn, "type")))
				if rn != nil {
					imgCreateNode.AddChild(rn)
				}
			}
			return imgCreateNode
		}
		ensureCreateFormatNode = func(exact, hint bool) *xmlx.Node {
			imgCreateFormatNode = ensureChild(ensureCreateNode(), "format")
			if hint {
				ensureChild(imgCreateFormatNode, "hint")
			}
			if exact {
				ensureChild(imgCreateFormatNode, "exact")
			}
			return imgCreateFormatNode
		}
	)
	surfaceNodes = append(surfaceNodes, xn)
	myID, imgID, initNode := attVal(xn.Parent, "sid"), "", subNode(xn, "init_as_target")
	if len(myID) == 0 {
		myID = attVal(xn.Parent, "ref")
	}
	if myID = ustr.StripPrefix(myID, "#"); len(myID) > 0 {
		if initNode != nil {
			imgID = fmt.Sprintf("img_target_%v", time.Now().UnixNano())
			imgNode, rn = xmlx.NewNode(xn.Type), xmlx.NewNode(xn.Type)
			imgNode.Name.Local, rn.Name.Local = "image", "renderable"
			setAttr(imgNode, "id", imgID, false)
			setAttr(rn, "share", "true", false)
			imgNode.AddChild(rn)
			moveNode(imgNode, nil, "library_images")
		} else if initNode = subNode(xn, "init_from"); initNode != nil {
			imgID = initNode.Value
		} else {
			for _, sn = range xn.Children {
				if strings.HasPrefix(sn.Name.Local, "init_") {
					initNode = sn
					break
				}
			}
			if initNode != nil {
				for _, sn = range initNode.Children {
					if imgID = attVal(sn, "ref"); len(imgID) > 0 {
						break
					}
				}
			}
		}
		if imgID = ustr.StripPrefix(imgID, "#"); len(imgID) > 0 {
			surfaceImages[myID] = imgID
			if imgNode == nil {
				for _, sn = range subNode(srcDoc.Root.Children[0], "library_images").Children {
					if attVal(sn, "id") == imgID {
						imgNode = sn
						break
					}
				}
			}
		}
		if imgNode != nil {
			if tn = subNode(xn, "format"); tn != nil {
				subNode(ensureCreateFormatNode(true, false), "exact").Value = tn.Value
			}
			if tn = subNode(xn, "format_hint"); tn != nil {
				rn = subNode(ensureCreateFormatNode(false, true), "hint")
				if sn = subNode(tn, "channels"); sn != nil {
					setAttr(rn, "channels", sn.Value, false)
				}
				if sn = subNode(tn, "range"); sn != nil {
					setAttr(rn, "range", sn.Value, false)
				}
				if sn = subNode(tn, "precision"); sn != nil {
					setAttr(rn, "precision", sn.Value, false)
				}
				if sn = subNode(tn, "option"); sn != nil {
					setAttr(rn, "space", sn.Value, false)
				}
			}
			if tn = subNode(xn, "size"); tn != nil {
				vals := ustr.Split(tn.Value, " ")
				rn = ensureChild(ensureCreateNode(), "size_exact")
				if len(vals) > 0 {
					setAttr(rn, "width", vals[0], false)
				}
				if len(vals) > 1 {
					setAttr(rn, "height", vals[1], false)
				}
				if len(vals) > 2 {
					setAttr(rn, "depth", vals[2], false)
				}
			}
			if tn = subNode(xn, "viewport_ratio"); tn != nil {
				vals := ustr.Split(tn.Value, " ")
				rn = ensureChild(ensureCreateNode(), "size_ratio")
				setAttr(rn, "width", vals[0], false)
				setAttr(rn, "height", vals[1], false)
			}
			if tn = subNode(xn, "mip_levels"); tn != nil {
				setAttr(ensureChild(ensureCreateNode(), "mips"), "levels", tn.Value, false)
			}
			if tn = subNode(xn, "mipmap_generate"); tn != nil {
				setAttr(ensureChild(ensureCreateNode(), "mips"), "auto_generate", tn.Value, false)
			}
		}
	}
}

func delAttr(xn *xmlx.Node, name string) {
	if Strict {
		pos := -1
		for i, att := range xn.Attributes {
			if att.Name.Local == name {
				pos = i
				break
			}
		}
		if pos >= 0 {
			logFmt("\t\tdelAttr %s.%s\n", xn.Name.Local, name)
			nuAtts := append(xn.Attributes[:pos], xn.Attributes[pos+1:]...)
			xn.Attributes = nuAtts
		}
	}
}

func delAtts(xn *xmlx.Node, names ...string) {
	for _, name := range names {
		delAttr(xn, name)
	}
}

func delNode(xn *xmlx.Node) {
	delNodeForce(xn, false)
}

func delNodeForce(xn *xmlx.Node, force bool) {
	if (Strict || force) && (xn != nil) {
		pos := -1
		for i, sn := range xn.Parent.Children {
			if sn == xn {
				pos = i
				break
			}
		}
		if pos >= 0 {
			logFmt("\t\tdelNode %s>%s = '%s'\n", xn.Parent.Name.Local, xn.Name.Local, xn.Value)
			nuNodes := append(xn.Parent.Children[:pos], xn.Parent.Children[pos+1:]...)
			xn.Parent.Children = nuNodes
		}
	}
}

func ensureChild(xn *xmlx.Node, name string) (sn *xmlx.Node) {
	if sn = subNode(xn, name); sn == nil {
		sn = xmlx.NewNode(xn.Type)
		sn.Name.Local = name
		xn.AddChild(sn)
	}
	return
}

func ensureSiblings(xn *xmlx.Node, namesVals map[string]string) {
	var sn *xmlx.Node
	for n, v := range namesVals {
		if sn = subNode(xn.Parent, n); sn == nil {
			sn = xmlx.NewNode(xn.Type)
			sn.Name.Local = n
			xn.Parent.AddChild(sn)
		}
		sn.Value = v
	}
}

func moveNode(xn, parent *xmlx.Node, parentName string) (newParent *xmlx.Node, oldPos int) {
	oldPos = -1
	if root := srcDoc.Root.Children[0]; parent == nil {
		if parent = subNode(root, parentName); parent == nil {
			parent = xmlx.NewNode(xn.Type)
			parent.Name.Local = parentName
			parent.Parent = root
			root.Children = append(root.Children, parent)
		}
	}
	if xn.Parent != parent {
		newParent = parent
		if xn.Parent != nil {
			for i, sn := range xn.Parent.Children {
				if sn == xn {
					oldPos = i
					break
				}
			}
			if oldPos >= 0 {
				xn.Parent.Children[oldPos] = nil
			}
		}
		xn.Parent = parent
		parent.Children = append(parent.Children, xn)
	}
	return
}

func renameAttr(xn *xmlx.Node, name, newName string) {
	for _, att := range xn.Attributes {
		if att.Name.Local == name {
			logFmt("\t\trenameAttr %s.%s => %s.%s\n", xn.Name.Local, name, xn.Name.Local, newName)
			att.Name.Local = newName
			break
		}
	}
}

func renameNode(xn *xmlx.Node, newName string) {
	logFmt("\t\trenameNode %s => %s", xn.Name.Local, newName)
	xn.Name.Local = newName
}

func replaceNode(xn, nn *xmlx.Node) {
	pos := -1
	for i, sn := range xn.Parent.Children {
		if sn == xn {
			pos = i
			break
		}
	}
	if pos >= 0 {
		xn.Parent.Children[pos] = nn
		nn.Parent = xn.Parent
		xn.Parent = nil
	}
}

func restrictAttr(xn *xmlx.Node, name string, min, max int64) {
	if Strict {
		for _, att := range xn.Attributes {
			if att.Name.Local == name {
				if val, err := strconv.ParseInt(att.Value, 10, 64); err == nil {
					logFmt("\t\trestrictAttr %s.%s\n", xn.Name.Local, name)
					if val < min {
						att.Value = strconv.FormatInt(min, 10)
					}
					if val > max {
						att.Value = strconv.FormatInt(max, 10)
					}
				}
				break
			}
		}
	}
}

func setAttr(xn *xmlx.Node, name, value string, onlySetIfEmpty bool) {
	for _, att := range xn.Attributes {
		if att.Name.Local == name {
			if (!onlySetIfEmpty) || (len(att.Value) == 0) {
				logFmt("\t\tsetAttr %s.%s = %s\n", xn.Name.Local, name, value)
				att.Value = value
			}
			return
		}
	}
	logFmt("\t\taddAttr %s.%s = %s\n", xn.Name.Local, name, value)
	att := &xmlx.Attr{Value: value}
	att.Name.Local = name
	xn.Attributes = append(xn.Attributes, att)
}

func subNode(xn *xmlx.Node, name string) (sn *xmlx.Node) {
	if sn = xn.SelectNode(ns, name); sn == nil {
		if sn = xn.SelectNode("", name); sn == nil {
			sn = xn.SelectNode("*", name)
		}
	}
	return
}

func subNodes(xn *xmlx.Node, name string) (sn []*xmlx.Node) {
	if sn = xn.SelectNodes(ns, name); len(sn) == 0 {
		if sn = xn.SelectNodes("", name); len(sn) == 0 {
			sn = xn.SelectNodes("*", name)
		}
	}
	return
}

func processNode(xn *xmlx.Node) {
	if (!Force) && (xn.Name.Local == "COLLADA") {
		if _, ver := util.ParseVersion(attVal(xn, "version")); ver >= 1.5 {
			skipped = true
			return
		}
	}
	xn.Name.Space = ""
	for _, att := range xn.Attributes {
		att.Name.Space = ""
	}
	switch xn.Name.Local {
	case "COLLADA":
		setAttr(xn, "version", "1.5", false)
		setAttr(xn, "xmlns", "http://www.collada.org/2008/03/COLLADASchema", false)
	case "array":
		if !ustr.IsOneOf(xn.Parent.Name.Local, "array", "newparam", "setparam") {
			delNode(xn)
		}
	case "argument", "texenv":
		delAttr(xn, "unit")
	case "cg_value_type", "connect_param", "generator", "tapered_capsule", "tapered_cylinder", "texture_unit":
		delNode(xn)
	case "code", "include":
		if !ustr.IsOneOf(xn.Parent.Name.Local, "profile_CG", "profile_GLES2", "profile_GLSL") {
			delNode(xn)
		}
	case "color_target", "depth_target", "stencil_target":
		if val := surfaceImages[xn.Value]; len(val) == 0 {
			setAttr(ensureChild(xn, "param"), "ref", xn.Value, false)
		} else {
			setAttr(ensureChild(xn, "instance_image"), "url", "#"+val, false)
		}
		xn.Value = ""
	case "float_array":
		restrictAttr(xn, "digits", 1, 17)
		restrictAttr(xn, "magnitude", -324, 308)
	case "image":
		convertImage(xn)
	case "instance_effect":
		if xn.Parent.Name.Local == "render" {
			id := fmt.Sprintf("render_%v", time.Now().UnixNano())
			matNode := xmlx.NewNode(xn.Type)
			matNode.Name.Local = "material"
			setAttr(matNode, "id", id, false)
			matsLibNode := ensureChild(srcDoc.Root.Children[0], "library_materials")
			matsLibNode.AddChild(matNode)
			oldParent := xn.Parent
			_, pos := moveNode(xn, matNode, "")
			instNode := xmlx.NewNode(xn.Type)
			instNode.Name.Local = "instance_material"
			setAttr(instNode, "url", "#"+id, false)
			instNode.Parent = oldParent
			oldParent.Children[pos] = instNode
		}
	case "magfilter", "minfilter", "mipfilter":
		if (xn.Value == "NONE") && (xn.Name.Local != "mipfilter") {
			xn.Value = "NEAREST"
		} else {
			switch xn.Value {
			case "NEAREST_MIPMAP_NEAREST":
				ensureSiblings(xn, map[string]string{"minfilter": "NEAREST", "mipfilter": "NEAREST"})
			case "LINEAR_MIPMAP_NEAREST":
				ensureSiblings(xn, map[string]string{"minfilter": "LINEAR", "mipfilter": "NEAREST"})
			case "NEAREST_MIPMAP_LINEAR":
				ensureSiblings(xn, map[string]string{"minfilter": "NEAREST", "mipfilter": "LINEAR"})
			case "LINEAR_MIPMAP_LINEAR":
				ensureSiblings(xn, map[string]string{"minfilter": "LINEAR", "mipfilter": "LINEAR"})
			}
		}
	case "mipmap_bias":
		renameNode(xn, "mip_bias")
	case "mipmap_maxlevel":
		renameNode(xn, "mip_max_level")
	case "newparam":
		if !ustr.IsOneOf(xn.Parent.Name.Local, "effect", "profile_CG", "profile_COMMON", "profile_GLSL", "profile_GLES", "profile_GLES2") {
			delNode(xn)
		}
	case "radius":
		if vals := ustr.Split(xn.Value, " "); (xn.Parent.Name.Local == "capsule") && (len(vals) > 0) && (len(vals) < 3) {
			for len(vals) < 3 {
				vals = append(vals, "1.0")
			}
			xn.Value = strings.Join(vals, " ")
		}
	case "setparam":
		if !ustr.IsOneOf(xn.Parent.Name.Local, "instance_effect", "usertype") {
			delNode(xn)
		}
	case "shader":
		convertShader(xn)
	case "surface":
		convertSurface(xn)
	case "texture_pipeline":
		if xn.Parent.Name.Local != "states" {
			delNode(xn)
		}
	case "transparent":
		setAttr(xn, "opaque", "A_ONE", true)
	case "usertype":
		renameAttr(xn, "name", "typename")
		if !ustr.IsOneOf(xn.Parent.Name.Local, "newparam", "setparam", "array", "bind_uniform") {
			delNode(xn)
		} else {
			for _, sn := range xn.Children {
				if sn.Name.Local != "setparam" {
					delNode(sn)
				}
			}
		}
	default:
		if (xn.Parent != nil) && (xn.Parent.Name.Local == "pass") {
			switch xn.Name.Local {
			case "annotate", "extra", "evaluate", "states", "program":
				break
			case "color_target", "depth_target", "stencil_target", "color_clear", "depth_clear", "stencil_clear", "draw":
				moveNode(xn, ensureChild(xn.Parent, "evaluate"), "")
			case "shader":
				moveNode(xn, ensureChild(xn.Parent, "program"), "")
			default:
				moveNode(xn, ensureChild(xn.Parent, "states"), "")
			}
		}
		if strings.HasPrefix(xn.Name.Local, "wrap_") && (xn.Value == "NONE") {
			xn.Value = "BORDER"
		}
		if (xn.Name.Local != "sampler") && strings.HasPrefix(xn.Name.Local, "sampler") && !strings.HasPrefix(xn.Name.Local, "sampler_") {
			if sn := subNode(xn, "source"); sn != nil {
				sn.Name.Local = "instance_image"
				setAttr(sn, "url", "#"+surfaceImages[sn.Value], false)
				sn.Value = ""
			}
		}
	}
	for _, sn := range xn.Children {
		processNode(sn)
	}
}
