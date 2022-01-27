package main

import "testing"

func Test_format(t *testing.T) {
	type args struct {
		dir       string
		file      string
		indent    int
		recursive bool
		verbose   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				dir:       "",
				file:      "",
				indent:    0,
				recursive: false,
				verbose:   false,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := formatRecursive(tt.args.dir, tt.args.file, tt.args.indent, tt.args.recursive, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("formatRecursive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
