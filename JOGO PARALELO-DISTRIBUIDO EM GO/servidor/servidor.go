package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Jogador struct {
	Nome string
	X, Y int
	LastSequence   int 
}

type Posicao struct {
	Nome string
	X, Y int
	SequenceNumber int 
}

type InteracaoAlavanca struct {
    NomeJogador    string
    NomeAlavanca   string
    SequenceNumber int
}

type ServidorJogo struct {
	Jogadores map[string]*Jogador
	Alavancas    map[string]map[string]int
	Mutex     sync.Mutex
}

// --------------------------------------------------------------------------------

func (s *ServidorJogo) AtivarAlavanca(interacao InteracaoAlavanca, resposta *string) error {

    //O defer garante que o mutex será destravado ao final do método, 
    // mesmo que haja retorno antecipado.
    s.Mutex.Lock()
    defer s.Mutex.Unlock()

    //verfica se existe a alavanca
    if s.Alavancas[interacao.NomeAlavanca] == nil {

        //se não existir, cria um mapa para a alavanca
        s.Alavancas[interacao.NomeAlavanca] = make(map[string]int)
    }

    //verifica a ultima interação
    lastSeq := s.Alavancas[interacao.NomeAlavanca][interacao.NomeJogador]
    //se a interação for menor ou igual a ultima interação, ignora pq ela é repetida
    if interacao.SequenceNumber <= lastSeq {
        *resposta = "Comando repetido ignorado"
        return nil
    }
    //se a interação for maior que a ultima interação, atualiza o valor
    s.Alavancas[interacao.NomeAlavanca][interacao.NomeJogador] = interacao.SequenceNumber
    *resposta = "Alavanca ativada!"
    return nil
}

//gerencia atualizações de posição dos jogadores no mapa, 
// garantindo consistência,evitando colisões e ignorando comandos duplicados
func (s *ServidorJogo) EnviarPosicao(p Posicao, resposta *string) error {

    //O defer garante que o mutex será destravado ao final do método, 
    // mesmo que haja retorno antecipado.
	s.Mutex.Lock()
    defer s.Mutex.Unlock()

    //Verificação de colisão de posição
    for nome, j := range s.Jogadores {
        if nome != p.Nome && j.X == p.X && j.Y == p.Y {
            *resposta = "Posição ocupada por outro jogador"
            return nil
        }
    }


    //Verificação de comando repetido e atualização de posição
    //Procura o jogador pelo nome no mapa de jogadores.
    j, ok := s.Jogadores[p.Nome]
    if ok {

        if p.SequenceNumber <= j.LastSequence {
            *resposta = "Comando repetido ignorado"
            return nil
        }

        j.LastSequence = p.SequenceNumber
        j.X = p.X
        j.Y = p.Y
        *resposta = "Posição atualizada"
    
    //Se o jogador não existir, cria um novo jogador com a posição recebida.
    } else {
        s.Jogadores[p.Nome] = &Jogador{
            Nome:         p.Nome,
            X:            p.X,
            Y:            p.Y,
            LastSequence: p.SequenceNumber,
        }
        *resposta = "Jogador novo adicionado"
    }
    return nil
}

// Retornar o estado atual de todos os jogadores conectados ao servidor. 
func (s *ServidorJogo) ObterEstado(_ string, estado *[]Jogador) error {

    //O defer garante que o mutex será destravado ao final do método, 
    // mesmo que haja retorno antecipado.
	s.Mutex.Lock()
    defer s.Mutex.Unlock()


	fmt.Println("Obtendo estado do servidor...")

    //Percorre todos os jogadores registrados no servidor.
	for _, j := range s.Jogadores {

        //*estado é um ponteiro ([]Jogador) que será preenchido com os jogadores atuais.
        //adiciona o jogador j ao final apontado por estado.
		*estado = append(*estado, *j)
	}
	return nil
}

func main() {

	servidor := &ServidorJogo{
		Jogadores: make(map[string]*Jogador),
		Alavancas: make(map[string]map[string]int),
	}
	rpc.Register(servidor)

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	fmt.Println("Servidor do jogo ativo na porta 1234...")

	for {
		conn, _ := ln.Accept()
		go rpc.ServeConn(conn)
	}

}
