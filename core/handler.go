package core

func HandlePING(args []string) (string, error) {
	if len(args) == 0 {
		return "PONG", nil
	}
	return args[0], nil
}
