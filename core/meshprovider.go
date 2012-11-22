package core

import (
	nga "github.com/go3d/go-ngine/assets"
)

type meshProviders struct {
	PrefabCube, PrefabPlane, PrefabPyramid, PrefabQuad, PrefabTri nga.MeshProvider
}

var (
	//	A collection of all "mesh providers" known to go:ngine.
	MeshProviders = &meshProviders { meshProviderPrefabCube, meshProviderPrefabPlane, meshProviderPrefabPyramid, meshProviderPrefabQuad, meshProviderPrefabTri }
)

func meshProviderPrefabCube (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(
		nga.MeshVertAtt3 { -1, -1, 1 }, nga.MeshVertAtt3 { 1, -1, 1 }, nga.MeshVertAtt3 { 1, 1, 1 },
		nga.MeshVertAtt3 { -1, 1, 1 }, nga.MeshVertAtt3 { -1, -1, -1 }, nga.MeshVertAtt3 { -1, 1, -1 },
		nga.MeshVertAtt3 { 1, 1, -1 }, nga.MeshVertAtt3 { 1, -1, -1 })
	meshData.AddTexCoords(nga.MeshVertAtt2 { 0, 0 }, nga.MeshVertAtt2 { 1, 0 }, nga.MeshVertAtt2 { 1, 1 }, nga.MeshVertAtt2 { 0, 1 })
	meshData.AddNormals(nga.MeshVertAtt3 { 0, 0, 1 }, nga.MeshVertAtt3 { 0, 0, -1 }, nga.MeshVertAtt3 { 0, 1, 0 }, nga.MeshVertAtt3 { 0, -1, 0 }, nga.MeshVertAtt3 { 1, 0, 0 }, nga.MeshVertAtt3 { -1, 0, 0 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.MeshVert { 0, 0, 0 }, nga.MeshVert { 1, 1, 0 }, nga.MeshVert { 2, 2, 0 } },	nga.MeshFace3 { nga.MeshVert { 0, 0, 0 }, nga.MeshVert { 2, 2, 0 }, nga.MeshVert { 3, 3, 0 } },		//	front
		nga.MeshFace3 { nga.MeshVert { 4, 0, 1 }, nga.MeshVert { 5, 1, 1 }, nga.MeshVert { 6, 2, 1 } },	nga.MeshFace3 { nga.MeshVert { 4, 0, 1 }, nga.MeshVert { 6, 2, 1 }, nga.MeshVert { 7, 3, 1 } },		//	back
		nga.MeshFace3 { nga.MeshVert { 5, 0, 2 }, nga.MeshVert { 3, 1, 2 }, nga.MeshVert { 2, 2, 2 } },	nga.MeshFace3 { nga.MeshVert { 5, 0, 2 }, nga.MeshVert { 2, 2, 2 }, nga.MeshVert { 6, 3, 2 } },		//	top
		nga.MeshFace3 { nga.MeshVert { 4, 0, 3 }, nga.MeshVert { 7, 1, 3 }, nga.MeshVert { 1, 2, 3 } },	nga.MeshFace3 { nga.MeshVert { 4, 0, 3 }, nga.MeshVert { 1, 2, 3 }, nga.MeshVert { 0, 3, 3 } },		//	bottom
		nga.MeshFace3 { nga.MeshVert { 7, 0, 4 }, nga.MeshVert { 6, 1, 4 }, nga.MeshVert { 2, 2, 4 } },	nga.MeshFace3 { nga.MeshVert { 7, 0, 4 }, nga.MeshVert { 2, 2, 4 }, nga.MeshVert { 1, 3, 4 } },		//	right
		nga.MeshFace3 { nga.MeshVert { 4, 0, 5 }, nga.MeshVert { 0, 1, 5 }, nga.MeshVert { 3, 2, 5 } },	nga.MeshFace3 { nga.MeshVert { 4, 0, 5 }, nga.MeshVert { 3, 2, 5 }, nga.MeshVert { 5, 3, 5 } })		//	left
	return
}

func meshProviderPrefabPlane (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.MeshVertAtt3 { -1, 0, 1 }, nga.MeshVertAtt3 { 1, 0, 1 }, nga.MeshVertAtt3 { -1, 0, -1 }, nga.MeshVertAtt3 { 1, 0, -1 })
	meshData.AddTexCoords(nga.MeshVertAtt2 { 0, 0 }, nga.MeshVertAtt2 { 1000, 0 }, nga.MeshVertAtt2 { 0, 1000 }, nga.MeshVertAtt2 { 1000, 1000 })
	meshData.AddNormals(nga.MeshVertAtt3 { 0, 1, 0 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.MeshVert { 0, 0, 0 }, nga.MeshVert { 1, 1, 0 }, nga.MeshVert { 2, 2, 0 } },
		nga.MeshFace3 { nga.MeshVert { 3, 3, 0 }, nga.MeshVert { 2, 2, 0 }, nga.MeshVert { 1, 1, 0 } })
	return
}

func meshProviderPrefabPyramid (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.MeshVertAtt3 { 0, 1, 0 }, nga.MeshVertAtt3 { -1, -1, 1 }, nga.MeshVertAtt3 { 1, -1, 1 }, nga.MeshVertAtt3 { 1, -1, -1 }, nga.MeshVertAtt3 { -1, -1, -1 })
	meshData.AddTexCoords(nga.MeshVertAtt2 { 0, 0 }, nga.MeshVertAtt2 { 1, 0 }, nga.MeshVertAtt2 { 1, 1 }, nga.MeshVertAtt2 { 0, 1})
	meshData.AddNormals(nga.MeshVertAtt3 { 0, 0, 1 }, nga.MeshVertAtt3 { 1, 0, 0 }, nga.MeshVertAtt3 { 0, 0, -1 }, nga.MeshVertAtt3 { -1, 0, 0 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.MeshVert { 0, 0, 0 }, nga.MeshVert { 1, 1, 0 }, nga.MeshVert { 2, 2, 0 } },
		nga.MeshFace3 { nga.MeshVert { 0, 1, 1 }, nga.MeshVert { 2, 2, 1 }, nga.MeshVert { 3, 3, 1 } },
		nga.MeshFace3 { nga.MeshVert { 0, 1, 2 }, nga.MeshVert { 3, 2, 2 }, nga.MeshVert { 4, 3, 2 } },
		nga.MeshFace3 { nga.MeshVert { 0, 0, 3 }, nga.MeshVert { 4, 1, 3 }, nga.MeshVert { 1, 2, 3 } })
	return
}

func meshProviderPrefabQuad (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.MeshVertAtt3 { 1, 1, 0 }, nga.MeshVertAtt3 { -1, 1, 0 }, nga.MeshVertAtt3 { -1, -1, 0 }, nga.MeshVertAtt3 { 1, -1, 0 })
	meshData.AddTexCoords(nga.MeshVertAtt2 { -0.125, 0 }, nga.MeshVertAtt2 { -0.125, 3 }, nga.MeshVertAtt2 { 1.125, 3 }, nga.MeshVertAtt2 { 1.125, 0 })
	meshData.AddNormals(nga.MeshVertAtt3 { 0, 0, 1 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.MeshVert { 0, 0, 0 }, nga.MeshVert { 1, 1, 0 }, nga.MeshVert { 2, 2, 0 } },
		nga.MeshFace3 { nga.MeshVert { 0, 0, 0 }, nga.MeshVert { 2, 2, 0 }, nga.MeshVert { 3, 3, 0 } })
	return
}

func meshProviderPrefabTri (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.MeshVertAtt3 { 0, 1, 0 }, nga.MeshVertAtt3 { -1, -1, 0 }, nga.MeshVertAtt3 { 1, -1, 0 })
	meshData.AddTexCoords(nga.MeshVertAtt2 { 0, 0 }, nga.MeshVertAtt2 { 3, 0 }, nga.MeshVertAtt2 { 3, 2 })
	meshData.AddNormals(nga.MeshVertAtt3 { 0, 0, 1 })
	meshData.AddFaces(nga.MeshFace3 { nga.MeshVert { 0, 0, 0 }, nga.MeshVert { 1, 1, 0 }, nga.MeshVert { 2, 2, 0 } })
	return
}
