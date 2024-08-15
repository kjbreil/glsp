package editreader

import "github.com/kjbreil/glsp/pkg/location"

func (f *File) Replace(n string, r *location.Range) {

	f.m.Lock()
	defer f.m.Unlock()

	r.Correct()

	f.edit.gotoPoint(r.End)
	// the end is the Char right after the end point
	end := f.edit.curr

	// go to the start of the edit
	f.edit.gotoPoint(r.Start)
	// since this is the first character to be removed go back one to get the start
	f.edit.reverse()
	start := f.edit.curr
	// once the function edits use the start to calculate position of Char's after start since they have changed
	defer start.setLoc()

	rc := readString(n)
	// if the replacement is nil then connected start to end
	if rc.IsEmpty() && rc.nextIsEmpty() {
		if end == nil {
			start.setNext(rc)
			return
		}
		start.setNext(end)
		return
	}
	// delete the last -1 since we are inserting
	rc.last().delete()

	// set the starts Next to the rc
	start.setNext(rc)
	// add the end onto the last RC
	rc.last().setNext(end)
}
