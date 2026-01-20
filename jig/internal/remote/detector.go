package remote

type PluginType int

const (
	TypeUnknown PluginType = iota
	TypeClaude
	TypeGemini
	TypeMixed
)

func (t PluginType) String() string {
	switch t {
	case TypeClaude:
		return "Claude Plugin"
	case TypeGemini:
		return "Gemini Extension"
	case TypeMixed:
		return "Mixed"
	default:
		return "Unknown"
	}
}

func DetectPluginType(path string) (PluginType, error) {
	// Not implemented yet
	return TypeUnknown, nil
}
