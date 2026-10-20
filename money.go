package meshio

import (
	"fmt"
	"math/big"
	"strings"
)

// MicrosToReais converte um valor em micros (string inteira, R$ 1 = "1000000") para reais
// em string decimal exata, sem ponto flutuante: "12500000" → "12.50", "-1" → "-0.000001".
// Mantém ao menos duas casas decimais.
func MicrosToReais(micros string) (string, error) {
	s := strings.TrimSpace(micros)
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return "", fmt.Errorf("meshio: valor em micros inválido: %q", micros)
	}
	neg := n.Sign() < 0
	if neg {
		n.Neg(n)
	}
	q, r := new(big.Int).QuoRem(n, big.NewInt(1_000_000), new(big.Int))
	frac := fmt.Sprintf("%06d", r.Int64())
	frac = strings.TrimRight(frac, "0")
	for len(frac) < 2 {
		frac += "0"
	}
	out := q.String() + "." + frac
	if neg {
		out = "-" + out
	}
	return out, nil
}

// MicrosToRat converte micros para um *big.Rat em reais, para cálculos exatos.
func MicrosToRat(micros string) (*big.Rat, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(micros), 10)
	if !ok {
		return nil, fmt.Errorf("meshio: valor em micros inválido: %q", micros)
	}
	return new(big.Rat).SetFrac(n, big.NewInt(1_000_000)), nil
}
