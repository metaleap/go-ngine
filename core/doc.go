//	The *core* package provides go:ngine client-side core functionality such as rendering and user I/O.
//	
//	________
//	
//	
//	First off, the go:ngine 3 RULES of bugless(er) development:
//	
//	1. you do NOT set exported fields directly, those are provided for simplified read-access. If a field is
//	meant to be modifiable, there will be a *SetFoo()* method or it will be documented explicitly as directly
//	modifiable, but such instances will be rare.
//	
//	2. you do NOT instantiate exported struct types directly, as in "new(Material)" or "&Material {}". Many of
//	those are exported only for documentation, but are to be instantiated only inside the go:ngine package.
//	For types to be instantiated by package-external code (ie. your code), go:ngine packages provide constructor
//	functions -- often in collections, ie. mat := Materials.New("arg") instead of, say, mat := NewMaterial("arg").
//	
//	3. Those 2 are *default* assumptions and modes of operation -- there are a few "if you know what you're
//	doing" exemptions and those are (or will be) explicitly documented as such.
//	
//	________
//	
//	TODO pkg doc
package core

type disposable interface {
	dispose()
}

type initable interface {
	init()
}
