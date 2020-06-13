package language

type pt struct{}

func (p pt) Welcome() string {
	return "Bem vindis ao DevOps Recife @%v! Sou o moderadorzis do grupis e estou aqui para garantirzis que não teremos spammers. Te enviarei um desafiis e espero que você me retornis em até 60 segundis ou terei que te convidarzis para sair do grupis!"
}
func (p pt) Wrong() string {
	return "Cacildis! Resposta incorretis! Que tentar de novis?"
}
func (p pt) Correct() string {
	return "Resposta corretis! Viva a liberdadis!"
}
func (p pt) Challenge() string {
	return "Qualzis o valorzis de %v %v %v?"
}
func (p pt) ID() string {
	return "pt"
}
