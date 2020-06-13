package language

type en struct{}

func (e en) Welcome() string {
	return "Welcome to the group!"
}
func (e en) Wrong() string {
	return "The answer is wrong, try again."
}
func (e en) Correct() string {
	return "That is correct, thank you!"
}
func (e en) Challenge() string {
	return "What is thre result of %v %v %v?"
}
func (e en) ID() string {
	return "en"
}
