//go:build windows
// +build windows

package process

// ReadHandle returns the windows handle for the process executing the command.
// WARNING: This uses reflection to read the private field of the `*os.Process`
// type. This function should never be used other than the extremely narrow
// use-case for which it was designed.
func ReadHandle(cmd *exec.Cmd) (window.Handle, error) {

	proc := cmd.Process
	if proc == nil {
		return 0, fmt.Errorf("command not started")
	}

	v := reflect.ValueOf(proc).Elem()
	field := v.FieldByName("handle")

	if field.IsZero() {
		return 0, fmt.Errorf("field not found")
	}

	// NOTE: Returning uintptr as uint64.
	handle := windows.Handle(field.Uint())

	return handle, nil
}
