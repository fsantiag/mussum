package language

type Pt struct{}

func (p Pt) Welcome() string {
	return "Bem vinda ao DevOps Recife! Enviei um desafio para você e espero que você me retorne em até 60 segundos ou terei que te convidar para sair do grupo! Nada pessoal, só não aceitamos spammers! :P"
}
func (p Pt) Wrong() string {
	return "Resposta incorreta!"
}
func (p Pt) Correct() string {
	return "Resposta correta!"
}
func (p Pt) Challenge() string {
	return "Qual o valor de %v %v %v?"
}
