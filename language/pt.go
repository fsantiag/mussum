package language

type Pt struct{}

func (p Pt) Welcome() string {
	return "Bem vindis ao DevOps Recife @%v! Sou o moderadorzis do grupis e estouzis aqui paris garantirzis que não teremis spammers. Te enviis um desafiis e esperis que você me retorne em até 60 segundis ou terei que te convidarzis paris sairzis do grupis!"
}
func (p Pt) Wrong() string {
	return "Cacildis! Resposta incorretis! Que tentar de novis?"
}
func (p Pt) Correct() string {
	return "Respostis corretis! Viva a liberdadis!"
}
func (p Pt) Challenge() string {
	return "Qualzis o valorzis de %v %v %v?"
}
func (p Pt) Id() string {
	return "pt"
}
