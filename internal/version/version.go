package version

import (
	"fmt"
	"runtime"
)

var (
	// ビルド時にリンカーで設定される変数
	Version   = "dev"     // git tag または "dev"
	GitCommit = "unknown" // git commit hash
	BuildDate = "unknown" // ビルド日時
)

// Info はバージョン情報を構造体として返す
type Info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	Platform  string
}

// Get は現在のバージョン情報を返す
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String はバージョン情報を文字列として返す
func (i Info) String() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, go: %s, platform: %s)",
		i.Version, i.GitCommit, i.BuildDate, i.GoVersion, i.Platform)
}
