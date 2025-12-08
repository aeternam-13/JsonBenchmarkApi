package main

type DummyClass struct {
	Foo1   string      `json:"foo1"`
	Foo2   string      `json:"foo2"`
	Foo3   string      `json:"foo3"`
	Foo4   string      `json:"foo4"`
	Foo5   string      `json:"foo5"`
	Foo6   string      `json:"foo6"`
	Target InnerTarget `json:"target"`
}

type InnerTarget struct {
	Foo1   string       `json:"foo1"`
	Foo2   string       `json:"foo2"`
	Foo3   string       `json:"foo3"`
	Target DeeperTarget `json:"target"`
}

// Here the target will be sent as a string instead of a struct nested object
type DummySlowerClass struct {
	Foo1   string `json:"foo1"`
	Foo2   string `json:"foo2"`
	Foo3   string `json:"foo3"`
	Foo4   string `json:"foo4"`
	Foo5   string `json:"foo5"`
	Foo6   string `json:"foo6"`
	Target string `json:"target"`
}

type InnerSlowerTarget struct {
	Foo1   string `json:"foo1"`
	Foo2   string `json:"foo2"`
	Foo3   string `json:"foo3"`
	Target string `json:"target"`
}

type DeeperTarget struct {
	Foo1         string `json:"foo1"`
	Foo2         string `json:"foo2"`
	Foo3         string `json:"foo3"`
	TargetString string `json:"targetString"`
}
