package strs

import "testing"

func TestConv2Snake(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "HTTPServer2", want: "http_server2"},
		{input: "HTTP2Server", want: "http2_server"},
		{input: "JSON2XML", want: "json2xml"},
		{input: "userID", want: "user_id"},
		{input: "already_snake", want: "already_snake"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Conv2Snake(tt.input); got != tt.want {
				t.Errorf("Conv2Snake(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConv2Camel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "http_server2", want: "HttpServer2"},
		{input: "HTTP2Server", want: "Http2Server"},
		{input: "JSON2XML", want: "Json2xml"},
		{input: "userID", want: "UserId"},
		{input: "alreadyCamel", want: "AlreadyCamel"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Conv2Camel(tt.input); got != tt.want {
				t.Errorf("Conv2Camel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConv2SnakeAndCamel(t *testing.T) {
	tests := []struct {
		input     string
		wantSnake string
		wantCamel string
	}{
		{input: "", wantSnake: "", wantCamel: ""},
		{input: "HTTPServer2", wantSnake: "http_server2", wantCamel: "HttpServer2"},
		{input: "HTTP2Server", wantSnake: "http2_server", wantCamel: "Http2Server"},
		{input: "JSON2XML", wantSnake: "json2xml", wantCamel: "Json2xml"},
		{input: "userID", wantSnake: "user_id", wantCamel: "UserId"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotSnake, gotCamel := Conv2SnakeAndCamel(tt.input)
			if gotSnake != tt.wantSnake || gotCamel != tt.wantCamel {
				t.Errorf("Conv2SnakeAndCamel(%q) = (%q, %q), want (%q, %q)", tt.input, gotSnake, gotCamel, tt.wantSnake, tt.wantCamel)
			}
		})
	}
}
