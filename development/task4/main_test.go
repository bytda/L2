package main

import "testing"

func TestAnagrams(t *testing.T) {
	got := Anagrams([]string{"стОлиК", "пятАк", "пятка", "слиток", "тяпка", "листок", "слово"})
	want := map[string][]string{"листок": {"листок", "слиток", "столик"}, "пятак": {"пятак", "пятка", "тяпка"}}
	if !Compare(got, want) {
		t.Errorf("got: %s want: %s\n", got, want)
	}

}

func Compare(x1, x2 map[string][]string) bool {
	for k, v := range x1 { //обходим первую мапу
		if len(v) != len(x2[k]) { // сравниваем длины срезов в значениях мап
			return false //если не совпало возращаем false
		}
		for i, x := range x2[k] { //перебираем элементы среза во второй мапе
			if v[i] != x { //сверяем элементы срезов в значениях мап
				return false //в случае несовпадения возращаем false
			}
		}
	}
	return true
}
