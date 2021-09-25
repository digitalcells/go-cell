package framework

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)

	Group(string) IGroup
}

type Group struct {
	parent *Group
	core   *Core
	prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

func (group *Group) getAbsolutePrefix() string {
	if group.parent == nil {
		return group.prefix
	}

	return group.parent.getAbsolutePrefix() + group.prefix
}

func (group *Group) Get(uri string, handler ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	group.core.Get(uri, handler)
}

func (group *Group) Post(uri string, handler ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	group.core.Post(uri, handler)
}

func (group *Group) Put(uri string, handler ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	group.core.Put(uri, handler)
}

func (group *Group) Delete(uri string, handler ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	group.core.Delete(uri, handler)
}

func (group *Group) Group(uri string) IGroup {
	c := NewGroup(group.core, uri)
	c.parent = group
	return c
}
