package cwl

// Entry represents fs entry, it means [File|Directory|Dirent]
type Entry struct {
	Class    string
	Location string
	Path     string
	Basename string
	File
	Directory
	Dirent
}

// File represents file entry.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#File
type File struct {
	Dirname string
	Size    int64
	Format  string
}

// Directory represents directory entry.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#Directory
type Directory struct {
	Listing []Entry
}

// Dirent represents ?
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#Dirent
type Dirent struct {
	Entry     string
	EntryName string
	Writable  bool
}

// NewList constructs a list of Entry from interface
func (e Entry) NewList(i interface{}) []Entry {
	dest := []Entry{}
	switch x := i.(type) {
	case string:
		dest = append(dest, Entry{}.New(x))
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Entry{}.New(v))
		}
	}
	return dest
}

// New constructs an Entry from interface
func (e Entry) New(i interface{}) Entry {
	dest := Entry{}
	switch x := i.(type) {
	case string:
		dest.Location = x
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case fieldEntryName:
				dest.EntryName = v.(string)
			case fieldEntry:
				dest.Entry = v.(string)
			case fieldWritable:
				dest.Writable = v.(bool)
			}
		}
	}
	return dest
}
