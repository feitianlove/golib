package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestIsNum(t *testing.T) {
	if !IsNum("111") {
		t.Errorf("IsNum Err")
	}
}

func TestRandomNum(t *testing.T) {
	r := RandomNum(10, 30)
	fmt.Printf("RandomNum:%d\n", r)
}

func TestStringInSlice(t *testing.T) {
	s := []string{"a", "b"}
	r := StringInSlice("a", s)
	fmt.Printf("StringInSlice:%v\n", r)
}

func TestIsInsideDocker(t *testing.T) {
	r := IsInsideDocker(111)
	fmt.Printf("IsInsideDocker:%v\n", r)
}

func TestReadFileAsLines(t *testing.T) {
	filePath := "utils.go"
	r, err := ReadFileAsLines(filePath)
	if err != nil {
		t.Errorf("TestReadFileAsLines Err:%s", err)
	} else {
		fmt.Printf("len:%d\n", len(r))
	}
}

func TestFileExists(t *testing.T) {
	filePath := "utils.go"
	r := FileExists(filePath)
	fmt.Printf("TestFileExists:%v\n", r)
}

func TestDirExists(t *testing.T) {
	filePath := "."
	r := DirExists(filePath)
	fmt.Printf("TestDirExists:%v\n", r)
}

func TestRemoveStrInSlices(t *testing.T) {
	s := []string{"a", "b"}
	r := RemoveStrInSlices(s, "a")
	fmt.Printf("TestRemoveStrInSlices:%b\n", len(r))
}

func TestIfValidIPs(t *testing.T) {
	fmt.Printf("TestIfValidIPs:%v\n", IfValidIP("1.1.1.1"))
	fmt.Printf("TestIfValidIPs:%v\n", IfValidIP("0.1.1.1"))
	fmt.Printf("TestIfValidIPs:%v\n", IfValidIP("10.1.1.1"))
	fmt.Printf("TestIfValidIPs:%v\n", IfValidIP("01.1.1.1"))
	fmt.Printf("TestIfValidIPs:%v\n", IfValidIP("256.1.1.1"))
	fmt.Printf("TestIfValidIPs:%v\n", IfValidIP("1.0.0.0"))
}

func TestSha1Bytes(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{data: []byte("sha1 this string")},
			want: "cf23df2207d99a74fbe169e3eba035e633b65d94",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha1Bytes(tt.args.data); got != tt.want {
				t.Errorf("Sha1Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSha1(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{data: "sha1 this string"},
			want: "cf23df2207d99a74fbe169e3eba035e633b65d94",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha1(tt.args.data); got != tt.want {
				t.Errorf("Sha1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInnerIP(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "fail0",
			args: args{"xx"},
			want: false,
		},
		{
			name: "fail1",
			args: args{"99.18.167.93"},
			want: false,
		},
		{
			name: "fail2",
			args: args{"99.0.0.0"},
			want: false,
		},
		{
			name: "ok",
			args: args{"10.0.0.0"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInnerIP(tt.args.ip); got != tt.want {
				t.Errorf("IsInnerIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZeroTime(t *testing.T) {
	fmt.Println(IsZeroTime(time.Time{}))
	fmt.Println(IsZeroTime(time.Unix(631126860, 0)))
	fmt.Println(IsZeroTime(time.Unix(631126862, 0)))
}

func TestSplitSlice(t *testing.T) {
	testS := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	a := SplitSlice(testS, 1)
	b := SplitSlice(testS, 2)
	c := SplitSlice(testS, 3)
	d := SplitSlice(testS, 10)
	println(len(a))
	println(len(b))
	println(len(c))
	println(len(d))
	println(1)
}

func TestGetLocalIP(t *testing.T) {
	fmt.Println(GetLocalIP())
}
