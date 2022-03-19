package http

import "testing"

func TestContentType(t *testing.T) {
	type args struct {
		url   string
		refer string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal test",
			args: args{
				url:   "https://google.com",
				refer: "",
			},
			want: "text/html",
		},
		{
			name: "youtube url",
			args: args{
				url: "https://www.youtube.com/watch?v=_sB2E1XnzOY",
				refer: "",
			},
			want: "text/html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentType, _ := ContentType(tt.args.url, tt.args.refer)
			if contentType != tt.want {
				t.Errorf("ContentType() = %s, want %s", contentType, tt.want)
			}
		})
	}
}

func TestSize(t *testing.T) {
	var err error
	type args struct {
		url   string
		refer string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal test",
			args: args{
				url:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				refer: "",
			},
		},
		{
			name: "velog image",
			args: args{
				url: "https://media.vlpt.us/images/lingodingo/post/86fddc0b-bbcb-434f-ae23-f150f5927aa1/400071_326705_4446.png?w=640",
				refer: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = Size(tt.args.url, tt.args.refer)
			if err != nil {
				t.Error()
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	var err error
	type args struct {
		url   string
		refer string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal test",
			args: args{
				url:   "https://google.com",
				refer: "",
			},
		},
		{
			name: "velog image",
			args: args{
				url: "https://media.vlpt.us/images/lingodingo/post/86fddc0b-bbcb-434f-ae23-f150f5927aa1/400071_326705_4446.png?w=640",
				refer: "",
			},
		},
		{
			name: "youtube url",
			args: args{
				url: "https://www.youtube.com/watch?v=_sB2E1XnzOY",
				refer: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = Headers(tt.args.url, tt.args.refer)
			if err != nil {
				t.Error()
			}
		})
	}
}
