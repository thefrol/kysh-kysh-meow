// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package metrica

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

func easyjson6252c418DecodeGithubComThefrolKyshKyshMeowInternalMetrica(in *jlexer.Lexer, out *Metricas) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Metricas, 0, 1)
			} else {
				*out = Metricas{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Metrica
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6252c418EncodeGithubComThefrolKyshKyshMeowInternalMetrica(out *jwriter.Writer, in Metricas) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Metricas) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6252c418EncodeGithubComThefrolKyshKyshMeowInternalMetrica(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Metricas) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6252c418EncodeGithubComThefrolKyshKyshMeowInternalMetrica(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metricas) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6252c418DecodeGithubComThefrolKyshKyshMeowInternalMetrica(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Metricas) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6252c418DecodeGithubComThefrolKyshKyshMeowInternalMetrica(l, v)
}
func easyjson6252c418DecodeGithubComThefrolKyshKyshMeowInternalMetrica1(in *jlexer.Lexer, out *Metrica) {
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
		case "id":
			out.ID = string(in.String())
		case "type":
			out.MType = string(in.String())
		case "delta":
			if in.IsNull() {
				in.Skip()
				out.Delta = nil
			} else {
				if out.Delta == nil {
					out.Delta = new(int64)
				}
				*out.Delta = int64(in.Int64())
			}
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				if out.Value == nil {
					out.Value = new(float64)
				}
				*out.Value = float64(in.Float64())
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
func easyjson6252c418EncodeGithubComThefrolKyshKyshMeowInternalMetrica1(out *jwriter.Writer, in Metrica) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.MType))
	}
	if in.Delta != nil {
		const prefix string = ",\"delta\":"
		out.RawString(prefix)
		out.Int64(int64(*in.Delta))
	}
	if in.Value != nil {
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Float64(float64(*in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Metrica) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6252c418EncodeGithubComThefrolKyshKyshMeowInternalMetrica1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Metrica) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6252c418EncodeGithubComThefrolKyshKyshMeowInternalMetrica1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metrica) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6252c418DecodeGithubComThefrolKyshKyshMeowInternalMetrica1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Metrica) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6252c418DecodeGithubComThefrolKyshKyshMeowInternalMetrica1(l, v)
}
