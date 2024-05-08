package pathfilter

// PathFilter decides if a path has to be included in the index
type PathFilterFunc = func(path string) bool

// DefaultPathFilterFunc includes all files always
var DefaultPathFilterFunc = func(path string) bool { return true }
