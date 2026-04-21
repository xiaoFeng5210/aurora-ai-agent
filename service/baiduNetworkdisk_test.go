package service

import (
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestGetBaiduNetworkdiskFileList(t *testing.T) {
	resp, err := GetBaiduNetworkdiskFileList("/oss/")
	if err != nil {
		t.Fatalf("GetBaiduNetworkdiskFileList failed: %v", err)
	}
	if resp.Errno != 0 {
		t.Fatalf("GetBaiduNetworkdiskFileList failed: %v", resp.Errno)
	}
	t.Logf("resp: %v", resp.List)
}

func TestNormalizeBaiduNetworkdiskDeletePaths(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		want    []string
		wantErr bool
	}{
		{
			name:  "keeps absolute app paths",
			paths: []string{"/apps/aurora-ai-agent知识库/super/a.txt"},
			want:  []string{"/apps/aurora-ai-agent知识库/super/a.txt"},
		},
		{
			name:  "prefixes relative paths and deduplicates",
			paths: []string{" super/a.txt ", "super/a.txt", "super/b.txt"},
			want: []string{
				"/apps/aurora-ai-agent知识库/super/a.txt",
				"/apps/aurora-ai-agent知识库/super/b.txt",
			},
		},
		{
			name:    "rejects paths outside app root",
			paths:   []string{"/other/a.txt"},
			wantErr: true,
		},
		{
			name:    "rejects app root deletion",
			paths:   []string{"/apps/aurora-ai-agent知识库"},
			wantErr: true,
		},
		{
			name:    "rejects empty input",
			paths:   []string{"  "},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeBaiduNetworkdiskDeletePaths(tt.paths)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
