package lib

import "context"

// nil means unknown: profile APIs may be unavailable. Never infer capability from an ACK.
func (o *OBSClient) recordPauseSupported(ctx context.Context) *bool {
	read := func(category, name string) string {
		var p struct {
			Value string `json:"parameterValue"`
		}
		if o.request(ctx, "GetProfileParameter", map[string]any{"parameterCategory": category, "parameterName": name}, &p) != nil {
			return ""
		}
		return p.Value
	}
	mode := read("Output", "Mode")
	var supported bool
	switch mode {
	case "Simple":
		quality := read("SimpleOutput", "RecQuality")
		if quality == "" {
			return nil
		}
		supported = quality != "Stream"
	case "Advanced":
		kind := read("AdvOut", "RecType")
		if kind == "FFmpeg" {
			return nil
		}
		encoder := read("AdvOut", "RecEncoder")
		if encoder == "" {
			return nil
		}
		supported = encoder != "none"
	default:
		return nil
	}
	return &supported
}
