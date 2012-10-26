package core

type tModels map[string]*TModel

type TModel struct {
	matKey, meshKey string

	mesh *TMesh
}

