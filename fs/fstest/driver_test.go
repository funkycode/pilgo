package fstest_test

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/fs/fstest"
)

var _ fs.Driver = new(fstest.Driver)

func TestDriver(t *testing.T) {
	t.Run("MkdirAll", testDriverMkdirAll)
	t.Run("ReadDir", testDriverReadDir)
	t.Run("ReadFile", testDriverReadFile)
	t.Run("Stat", testDriverStat)
	t.Run("Symlink", testDriverSymlink)
	t.Run("WriteFile", testDriverWriteFile)
}

func testDriverMkdirAll(t *testing.T) {
	errMkdirAll := errors.New("MkdirAll")
	testCases := []struct {
		drv     fstest.Driver
		dirname string
		err     error
	}{
		{
			drv: fstest.Driver{
				MkdirAllErr: map[string]error{
					"foo": errMkdirAll,
				},
			},
			dirname: "foo",
			err:     errMkdirAll,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.dirname, func(t *testing.T) {
			err := tc.drv.MkdirAll(tc.dirname)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			hasBeenCalled, args := tc.drv.HasBeenCalled(tc.drv.MkdirAll)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			callstack := fstest.CallStack{fstest.Args{tc.dirname}}
			if want, got := callstack, args; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testDriverReadDir(t *testing.T) {
	errReadDir := errors.New("ReadDir")
	testCases := []struct {
		drv     fstest.Driver
		dirname string
		want    []fs.FileInfo
		err     error
	}{
		{
			drv: fstest.Driver{
				ReadDirReturn: map[string][]fs.FileInfo{
					"foo": {
						fstest.FileInfo{
							NameReturn:     "foo_1",
							ExistsReturn:   true,
							IsDirReturn:    false,
							LinknameReturn: "1_oof",
							PermReturn:     0o777,
						},
						fstest.FileInfo{
							NameReturn:     "foo_2",
							ExistsReturn:   true,
							IsDirReturn:    true,
							LinknameReturn: "2_oof",
							PermReturn:     0o655,
						},
					},
				},
				ReadDirErr: nil,
			},
			dirname: "foo",
			want: []fs.FileInfo{
				fstest.FileInfo{
					NameReturn:     "foo_1",
					ExistsReturn:   true,
					IsDirReturn:    false,
					LinknameReturn: "1_oof",
					PermReturn:     0o777,
				},
				fstest.FileInfo{
					NameReturn:     "foo_2",
					ExistsReturn:   true,
					IsDirReturn:    true,
					LinknameReturn: "2_oof",
					PermReturn:     0o655,
				},
			},
			err: nil,
		},
		{
			drv: fstest.Driver{
				ReadDirReturn: map[string][]fs.FileInfo{
					"bar": {
						fstest.FileInfo{
							NameReturn:     "bar_1",
							ExistsReturn:   true,
							IsDirReturn:    false,
							LinknameReturn: "1_rab",
							PermReturn:     0o777,
						},
						fstest.FileInfo{
							NameReturn:     "bar_2",
							ExistsReturn:   true,
							IsDirReturn:    true,
							LinknameReturn: "2_rab",
							PermReturn:     0o655,
						},
					},
				},
				ReadDirErr: nil,
			},
			dirname: "bar",
			want: []fs.FileInfo{
				fstest.FileInfo{
					NameReturn:     "bar_1",
					ExistsReturn:   true,
					IsDirReturn:    false,
					LinknameReturn: "1_rab",
					PermReturn:     0o777,
				},
				fstest.FileInfo{
					NameReturn:     "bar_2",
					ExistsReturn:   true,
					IsDirReturn:    true,
					LinknameReturn: "2_rab",
					PermReturn:     0o655,
				},
			},
			err: nil,
		},
		{
			drv: fstest.Driver{
				ReadDirReturn: map[string][]fs.FileInfo{
					"foo": {
						fstest.FileInfo{
							NameReturn:     "foo_1",
							ExistsReturn:   true,
							IsDirReturn:    false,
							LinknameReturn: "1_oof",
							PermReturn:     0o777,
						},
						fstest.FileInfo{
							NameReturn:     "foo_2",
							ExistsReturn:   true,
							IsDirReturn:    true,
							LinknameReturn: "2_oof",
							PermReturn:     0o655,
						},
					},
				},
				ReadDirErr: map[string]error{
					"foo": errReadDir,
				},
			},
			dirname: "foo",
			want: []fs.FileInfo{
				fstest.FileInfo{
					NameReturn:     "foo_1",
					ExistsReturn:   true,
					IsDirReturn:    false,
					LinknameReturn: "1_oof",
					PermReturn:     0o777,
				},
				fstest.FileInfo{
					NameReturn:     "foo_2",
					ExistsReturn:   true,
					IsDirReturn:    true,
					LinknameReturn: "2_oof",
					PermReturn:     0o655,
				},
			},
			err: errReadDir,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.dirname, func(t *testing.T) {
			files, err := tc.drv.ReadDir(tc.dirname)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, files; !cmp.Equal(got, want) {
				t.Fatalf("\n-want +got\n%s", cmp.Diff(want, got))
			}
			hasBeenCalled, args := tc.drv.HasBeenCalled(tc.drv.ReadDir)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			callstack := fstest.CallStack{fstest.Args{tc.dirname}}
			if want, got := callstack, args; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testDriverReadFile(t *testing.T) {
	errReadFile := errors.New("ReadFile")
	testCases := []struct {
		drv      fstest.Driver
		filename string
		want     []byte
		err      error
	}{
		{
			drv: fstest.Driver{
				ReadFileReturn: map[string][]byte{
					"foo": []byte("foo"),
				},
				ReadFileErr: nil,
			},
			filename: "foo",
			want:     []byte("foo"),
			err:      nil,
		},
		{
			drv: fstest.Driver{
				ReadFileReturn: map[string][]byte{
					"foo": []byte("foo"),
				},
				ReadFileErr: map[string]error{
					"foo": errReadFile,
				},
			},
			filename: "foo",
			want:     []byte("foo"),
			err:      errReadFile,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			b, err := tc.drv.ReadFile(tc.filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, b; string(got) != string(want) {
				t.Fatalf("want %q, got %q", want, got)
			}
			hasBeenCalled, args := tc.drv.HasBeenCalled(tc.drv.ReadFile)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			callstack := fstest.CallStack{fstest.Args{tc.filename}}
			if want, got := callstack, args; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testDriverStat(t *testing.T) {
	errStat := errors.New("Stat")
	testCases := []struct {
		fs       fstest.Driver
		filename string
		want     fs.FileInfo
		err      error
	}{
		{
			fs: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						NameReturn:     "foo_1",
						ExistsReturn:   true,
						IsDirReturn:    false,
						LinknameReturn: "1_oof",
						PermReturn:     0o777,
					},
				},
				StatErr: nil,
			},
			filename: "foo",
			want: fstest.FileInfo{
				NameReturn:     "foo_1",
				ExistsReturn:   true,
				IsDirReturn:    false,
				LinknameReturn: "1_oof",
				PermReturn:     0o777,
			},
			err: nil,
		},
		{
			fs: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						NameReturn:     "foo_1",
						ExistsReturn:   true,
						IsDirReturn:    false,
						LinknameReturn: "1_oof",
						PermReturn:     0o777,
					},
				},
				StatErr: map[string]error{
					"foo": errStat,
				},
			},
			filename: "foo",
			want: fstest.FileInfo{
				NameReturn:     "foo_1",
				ExistsReturn:   true,
				IsDirReturn:    false,
				LinknameReturn: "1_oof",
				PermReturn:     0o777,
			},
			err: errStat,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			fi, err := tc.fs.Stat(tc.filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, fi; !cmp.Equal(got, want) {
				t.Fatalf("\n-want +got\n%s", cmp.Diff(want, got))
			}
			hasBeenCalled, args := tc.fs.HasBeenCalled(tc.fs.Stat)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			callstack := fstest.CallStack{fstest.Args{tc.filename}}
			if want, got := callstack, args; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testDriverSymlink(t *testing.T) {
	errSymlink := errors.New("Symlink")
	testCases := []struct {
		drv     fstest.Driver
		oldname string
		newname string
		err     error
	}{
		{
			drv: fstest.Driver{
				SymlinkErr: map[string]error{
					"foo": errSymlink,
				},
			},
			oldname: "foo",
			newname: "bar",
			err:     errSymlink,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.oldname+" "+tc.newname, func(t *testing.T) {
			err := tc.drv.Symlink(tc.oldname, tc.newname)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			hasBeenCalled, args := tc.drv.HasBeenCalled(tc.drv.Symlink)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			callstack := fstest.CallStack{fstest.Args{tc.oldname, tc.newname}}
			if want, got := callstack, args; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testDriverWriteFile(t *testing.T) {
	errWriteFile := errors.New("WriteFile")
	testCases := []struct {
		drv      fstest.Driver
		filename string
		data     []byte
		perm     os.FileMode
		err      error
	}{
		{
			drv: fstest.Driver{
				WriteFileErr: map[string]error{
					"foo": errWriteFile,
				},
			},
			filename: "foo",
			data:     []byte("pilgrim"),
			perm:     0o777,
			err:      errWriteFile,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			err := tc.drv.WriteFile(tc.filename, tc.data, tc.perm)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			hasBeenCalled, args := tc.drv.HasBeenCalled(tc.drv.WriteFile)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			callstack := fstest.CallStack{fstest.Args{tc.filename, tc.data, tc.perm}}
			if want, got := callstack, args; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}