package language

type En struct{}

func (e En) Welcome() string {
	return "Welcome to the group!"
}
func (e En) Wrong() string {
	return "The answer is wrong, try again."
}
func (e En) Correct() string {
	return "That is correct, thank you!"
}
func (e En) Challenge() string {
	return "What is thre result of %v %v %v?"
}
