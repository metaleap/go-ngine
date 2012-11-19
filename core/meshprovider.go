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
		nga.VertAtt3 { -1, -1, 1 }, nga.VertAtt3 { 1, -1, 1 }, nga.VertAtt3 { 1, 1, 1 },
		nga.VertAtt3 { -1, 1, 1 }, nga.VertAtt3 { -1, -1, -1 }, nga.VertAtt3 { -1, 1, -1 },
		nga.VertAtt3 { 1, 1, -1 }, nga.VertAtt3 { 1, -1, -1 })
	meshData.AddTexCoords(nga.VertAtt2 { 0, 0 }, nga.VertAtt2 { 1, 0 }, nga.VertAtt2 { 1, 1 }, nga.VertAtt2 { 0, 1 })
	meshData.AddNormals(nga.VertAtt3 { 0, 0, 1 }, nga.VertAtt3 { 0, 0, -1 }, nga.VertAtt3 { 0, 1, 0 }, nga.VertAtt3 { 0, -1, 0 }, nga.VertAtt3 { 1, 0, 0 }, nga.VertAtt3 { -1, 0, 0 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.Vert { 0, 0, 0 }, nga.Vert { 1, 1, 0 }, nga.Vert { 2, 2, 0 } },	nga.MeshFace3 { nga.Vert { 0, 0, 0 }, nga.Vert { 2, 2, 0 }, nga.Vert { 3, 3, 0 } },		//	front
		nga.MeshFace3 { nga.Vert { 4, 0, 1 }, nga.Vert { 5, 1, 1 }, nga.Vert { 6, 2, 1 } },	nga.MeshFace3 { nga.Vert { 4, 0, 1 }, nga.Vert { 6, 2, 1 }, nga.Vert { 7, 3, 1 } },		//	back
		nga.MeshFace3 { nga.Vert { 5, 0, 2 }, nga.Vert { 3, 1, 2 }, nga.Vert { 2, 2, 2 } },	nga.MeshFace3 { nga.Vert { 5, 0, 2 }, nga.Vert { 2, 2, 2 }, nga.Vert { 6, 3, 2 } },		//	top
		nga.MeshFace3 { nga.Vert { 4, 0, 3 }, nga.Vert { 7, 1, 3 }, nga.Vert { 1, 2, 3 } },	nga.MeshFace3 { nga.Vert { 4, 0, 3 }, nga.Vert { 1, 2, 3 }, nga.Vert { 0, 3, 3 } },		//	bottom
		nga.MeshFace3 { nga.Vert { 7, 0, 4 }, nga.Vert { 6, 1, 4 }, nga.Vert { 2, 2, 4 } },	nga.MeshFace3 { nga.Vert { 7, 0, 4 }, nga.Vert { 2, 2, 4 }, nga.Vert { 1, 3, 4 } },		//	right
		nga.MeshFace3 { nga.Vert { 4, 0, 5 }, nga.Vert { 0, 1, 5 }, nga.Vert { 3, 2, 5 } },	nga.MeshFace3 { nga.Vert { 4, 0, 5 }, nga.Vert { 3, 2, 5 }, nga.Vert { 5, 3, 5 } })		//	left
	return
}

func meshProviderPrefabPlane (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.VertAtt3 { -1, 0, 1 }, nga.VertAtt3 { 1, 0, 1 }, nga.VertAtt3 { -1, 0, -1 }, nga.VertAtt3 { 1, 0, -1 })
	meshData.AddTexCoords(nga.VertAtt2 { 0, 0 }, nga.VertAtt2 { 1000, 0 }, nga.VertAtt2 { 0, 1000 }, nga.VertAtt2 { 1000, 1000 })
	meshData.AddNormals(nga.VertAtt3 { 0, 1, 0 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.Vert { 0, 0, 0 }, nga.Vert { 1, 1, 0 }, nga.Vert { 2, 2, 0 } },
		nga.MeshFace3 { nga.Vert { 3, 3, 0 }, nga.Vert { 2, 2, 0 }, nga.Vert { 1, 1, 0 } })
	return
}

func meshProviderPrefabPyramid (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.VertAtt3 { 0, 1, 0 }, nga.VertAtt3 { -1, -1, 1 }, nga.VertAtt3 { 1, -1, 1 }, nga.VertAtt3 { 1, -1, -1 }, nga.VertAtt3 { -1, -1, -1 })
	meshData.AddTexCoords(nga.VertAtt2 { 0, 0 }, nga.VertAtt2 { 1, 0 }, nga.VertAtt2 { 1, 1 }, nga.VertAtt2 { 0, 1})
	meshData.AddNormals(nga.VertAtt3 { 0, 0, 1 }, nga.VertAtt3 { 1, 0, 0 }, nga.VertAtt3 { 0, 0, -1 }, nga.VertAtt3 { -1, 0, 0 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.Vert { 0, 0, 0 }, nga.Vert { 1, 1, 0 }, nga.Vert { 2, 2, 0 } },
		nga.MeshFace3 { nga.Vert { 0, 1, 1 }, nga.Vert { 2, 2, 1 }, nga.Vert { 3, 3, 1 } },
		nga.MeshFace3 { nga.Vert { 0, 1, 2 }, nga.Vert { 3, 2, 2 }, nga.Vert { 4, 3, 2 } },
		nga.MeshFace3 { nga.Vert { 0, 0, 3 }, nga.Vert { 4, 1, 3 }, nga.Vert { 1, 2, 3 } })
	return
}

func meshProviderPrefabQuad (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.VertAtt3 { 1, 1, 0 }, nga.VertAtt3 { -1, 1, 0 }, nga.VertAtt3 { -1, -1, 0 }, nga.VertAtt3 { 1, -1, 0 })
	meshData.AddTexCoords(nga.VertAtt2 { -0.125, 0 }, nga.VertAtt2 { -0.125, 3 }, nga.VertAtt2 { 1.125, 3 }, nga.VertAtt2 { 1.125, 0 })
	meshData.AddNormals(nga.VertAtt3 { 0, 0, 1 })
	meshData.AddFaces(
		nga.MeshFace3 { nga.Vert { 0, 0, 0 }, nga.Vert { 1, 1, 0 }, nga.Vert { 2, 2, 0 } },
		nga.MeshFace3 { nga.Vert { 0, 0, 0 }, nga.Vert { 2, 2, 0 }, nga.Vert { 3, 3, 0 } })
	return
}

func meshProviderPrefabTri (args ... interface {}) (meshData *nga.MeshData, err error) {
	meshData = nga.NewMeshData()
	meshData.AddPositions(nga.VertAtt3 { 0, 1, 0 }, nga.VertAtt3 { -1, -1, 0 }, nga.VertAtt3 { 1, -1, 0 })
	meshData.AddTexCoords(nga.VertAtt2 { 0, 0 }, nga.VertAtt2 { 3, 0 }, nga.VertAtt2 { 3, 2 })
	meshData.AddNormals(nga.VertAtt3 { 0, 0, 1 })
	meshData.AddFaces(nga.MeshFace3 { nga.Vert { 0, 0, 0 }, nga.Vert { 1, 1, 0 }, nga.Vert { 2, 2, 0 } })
	return
}
