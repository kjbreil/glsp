package uri

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

type DocumentURI string

type URI = string

func (uri *DocumentURI) UnmarshalText(data []byte) (err error) {
	*uri, err = ParseDocumentURI(string(data))
	return
}

// Path returns the file path for the given URI.
//
// DocumentURI("").Path() returns the empty string.
//
// Path panics if called on a URI that is not a valid filename.
func (uri DocumentURI) Path() string {
	filename, err := filename(uri)
	if err != nil {
		// e.g. ParseRequestURI failed.
		//
		// This can only affect DocumentURIs created by
		// direct string manipulation; all DocumentURIs
		// received from the client pass through
		// ParseRequestURI, which ensures validity.
		panic(err)
	}
	return filepath.FromSlash(filename)
}

func (uri *DocumentURI) IsPath(path string) bool {
	uriPath := uri.Path()
	if uriPath == path {
		return true
	}
	return false
}

func (uri *DocumentURI) Schema() string {
	u, err := url.ParseRequestURI(string(*uri))
	if err != nil {
		return ""
	}
	return u.Scheme
}

// // Dir returns the URI for the directory containing the receiver.
// func (uri DocumentURI) Dir() DocumentURI {
// 	// This function could be more efficiently implemented by avoiding any call
// 	// to Path(), but at least consolidates URI manipulation.
// 	return URIFromPath(filepath.Dir(uri.Path()))
// }

func filename(uri DocumentURI) (string, error) {
	if uri == "" {
		return "", nil
	}

	// This conservative check for the common case
	// of a simple non-empty absolute POSIX filename
	// avoids the allocation of a net.URL.
	if strings.HasPrefix(string(uri), "file:///") {
		rest := string(uri)[len("file://"):] // leave one slash
		for i := 0; i < len(rest); i++ {
			b := rest[i]
			// Reject these cases:
			if b < ' ' || b == 0x7f || // control character
				b == '%' || b == '+' || // URI escape
				b == ':' || // Windows drive letter
				b == '@' || b == '&' || b == '?' { // authority or query
				goto slow
			}
		}
		return rest, nil
	}
slow:

	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return "", err
	}
	// if u.Scheme != fileScheme {
	// 	return "", fmt.Errorf("only file URIs are supported, got %q from %q", u.Scheme, uri)
	// }
	// If the URI is a Windows URI, we trim the leading "/" and uppercase
	// the drive letter, which will never be case sensitive.
	if isWindowsDriveURIPath(u.Path) {
		u.Path = strings.ToUpper(string(u.Path[1])) + u.Path[2:]
	}
	if u.Host != "" {
		return "//" + u.Host + u.Path, nil
	}

	return u.Path, nil
}

// ParseDocumentURI interprets a string as a DocumentURI, applying VS
// Code workarounds; see [DocumentURI.UnmarshalText] for details.
func ParseDocumentURI(s string) (DocumentURI, error) {
	if s == "" {
		return "", nil
	}

	split := strings.SplitN(s, ":", 2)
	if len(split) == 1 {
		return "", fmt.Errorf("DocumentURI must contain a scheme: %s", s)
	}
	schema := split[0]
	// when using fs workspace folder VSCODE returns only a single slash after schema:
	if schema != "file" {
		s = schema + ":///" + s[len(schema+":/"):]
	}

	if !strings.HasPrefix(s, schema+":/") {
		return "", fmt.Errorf("DocumentURI scheme is not '%s': %s", schema, s)
	}

	// Even though the input is a URI, it may not be in canonical form. VS Code
	// in particular over-escapes :, @, etc. Unescape and re-encode to canonicalize.
	path, err := url.PathUnescape(s[len(schema+"://"):])
	if err != nil {
		return "", err
	}

	// File URIs from Windows may have lowercase drive letters.
	// Since drive letters are guaranteed to be case insensitive,
	// we change them to uppercase to remain consistent.
	// For example, file:///c:/x/y/z becomes file:///C:/x/y/z.
	if isWindowsDriveURIPath(path) {
		path = path[:1] + strings.ToUpper(string(path[1])) + path[2:]
	}
	u := url.URL{Scheme: schema, Path: path}
	return DocumentURI(u.String()), nil
}

// URIFromPath returns DocumentURI for the supplied file path.
// Given "", it returns "".
func URIFromPath(schema string, path string) DocumentURI {
	if path == "" {
		return ""
	}
	if !isWindowsDrivePath(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}
	// Check the file path again, in case it became absolute.
	if isWindowsDrivePath(path) {
		path = "/" + strings.ToUpper(string(path[0])) + path[1:]
	}
	var host string

	// on windows if the path starts with two backslashes extract the host
	if runtime.GOOS == "windows" && strings.HasPrefix(path, "\\\\") {
		split := strings.SplitN(path, "\\", 4)
		host = split[2]
		// the split removes the start backslash so add it back in and update path
		path = "\\" + split[3]
	}

	path = filepath.ToSlash(path)
	u := url.URL{
		Scheme: schema,
		Path:   path,
		Host:   host,
	}
	return DocumentURI(u.String())
}

// isWindowsDrivePath returns true if the file path is of the form used by
// Windows. We check if the path begins with a drive letter, followed by a ":".
// For example: C:/x/y/z.
func isWindowsDrivePath(path string) bool {
	if len(path) < 3 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

// isWindowsDriveURIPath returns true if the file URI is of the format used by
// Windows URIs. The url.Parse package does not specially handle Windows paths
// (see golang/go#6027), so we check if the URI path has a drive prefix (e.g. "/C:").
func isWindowsDriveURIPath(uri string) bool {
	if len(uri) < 4 {
		return false
	}
	return uri[0] == '/' && unicode.IsLetter(rune(uri[1])) && uri[2] == ':'
}
