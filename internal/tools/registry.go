package tools

import "slices"

// Registry stores tool implementations by name.
type Registry struct {
	tools map[ToolName]Tool
}

// NewRegistry returns an empty tool registry.
func NewRegistry() *Registry {
	return &Registry{tools: make(map[ToolName]Tool)}
}

// NewDefaultRegistry returns the Luna v1 tool set.
func NewDefaultRegistry() *Registry {
	r := NewRegistry()
	r.Register(NewWebSearch())
	r.Register(NewWebFetch())
	for _, name := range DefaultToolNames() {
		if _, ok := r.Get(name); !ok {
			r.Register(NewStub(name))
		}
	}
	return r
}

// Register adds or replaces a tool implementation.
func (r *Registry) Register(t Tool) {
	r.tools[t.Name()] = t
}

// Get returns a registered tool by name.
func (r *Registry) Get(name ToolName) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// MustGet returns a registered tool or ErrToolNotFound.
func (r *Registry) MustGet(name ToolName) (Tool, error) {
	t, ok := r.Get(name)
	if !ok {
		return nil, ErrToolNotFound
	}
	return t, nil
}

// Names returns sorted registered tool names.
func (r *Registry) Names() []ToolName {
	names := make([]ToolName, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// DefaultToolNames is the current Luna tool set.
func DefaultToolNames() []ToolName {
	return []ToolName{
		ToolRGSearch,
		ToolReadFile,
		ToolGlob,
		ToolGrep,
		ToolWebSearch,
		ToolWebFetch,
		ToolEdit,
		ToolMultiEdit,
	}
}
