// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package mem

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson7fbae19eDecodeGithubComThefrolKyshKyshMeowInternalServerStoragev2Mem(in *jlexer.Lexer, out *FileData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Counters":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Counters = make(IntMap)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 int64
					v1 = int64(in.Int64())
					(out.Counters)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "Gauges":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Gauges = make(FloatMap)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 float64
					v2 = float64(in.Float64())
					(out.Gauges)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7fbae19eEncodeGithubComThefrolKyshKyshMeowInternalServerStoragev2Mem(out *jwriter.Writer, in FileData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Counters\":"
		out.RawString(prefix[1:])
		if in.Counters == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v3First := true
			for v3Name, v3Value := range in.Counters {
				if v3First {
					v3First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v3Name))
				out.RawByte(':')
				out.Int64(int64(v3Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"Gauges\":"
		out.RawString(prefix)
		if in.Gauges == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v4First := true
			for v4Name, v4Value := range in.Gauges {
				if v4First {
					v4First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v4Name))
				out.RawByte(':')
				out.Float64(float64(v4Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FileData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7fbae19eEncodeGithubComThefrolKyshKyshMeowInternalServerStoragev2Mem(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FileData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7fbae19eEncodeGithubComThefrolKyshKyshMeowInternalServerStoragev2Mem(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FileData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7fbae19eDecodeGithubComThefrolKyshKyshMeowInternalServerStoragev2Mem(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FileData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7fbae19eDecodeGithubComThefrolKyshKyshMeowInternalServerStoragev2Mem(l, v)
}