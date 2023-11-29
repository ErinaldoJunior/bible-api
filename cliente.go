package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Biblia :)
type Biblia struct {
	Livro   int    `json:"livro"`
	Capitulo int `json:"capitulo"`
	Versiculo int `json:"versiculo"`
	Texto string `json:"texto"`

}
type Capitulo struct {
	Capitulo int `json:"capitulo"`
	Versiculo int `json:"versiculo"`
	Texto string `json:"texto"`
}

type Versiculo struct {
	Capitulo int `json:"capitulo"`
	Versiculo int `json:"versiculo"`
	Texto string `json:"texto"`
}

// HandlerBiblia analisa o request e delega para função adequada
func HandlerBiblia(w http.ResponseWriter, r *http.Request) {

	// Obtendo o caminho da URL
	path := r.URL.Path

	// Removendo o prefixo "/biblia/"
	pathWithoutPrefix := strings.TrimPrefix(path, "/biblia/")
	
	// Dividindo o caminho em partes usando "/"
	parts := strings.Split(pathWithoutPrefix, "/")
	
	// Verificando o número de partes
	if pathWithoutPrefix == "" {
        todosOsLivros(w, r)
	} else if len(parts) == 1 {
		// Se houver apenas uma parte, é um pedido de livro
		idLivro := parts[0]
		id, _ := strconv.Atoi(idLivro)
		livroPorID(w, r, id)
	} else if len(parts) == 2 {
		// Se houver duas partes, é um pedido de livro-capitulo
		idLivro := parts[0]
		idCapitulo := parts[1]
		idL, _ := strconv.Atoi(idLivro)
		idC, _ := strconv.Atoi(idCapitulo)
		livroCap(w, r, idL, idC)
	} else if len(parts) == 3 {
		// Se houver 3 partes, é um pedido de livro-capitulo-versiculo
		idLivro := parts[0]
		idCapitulo := parts[1]
		idVersiculo := parts[2]
		idL, _ := strconv.Atoi(idLivro)
		idC, _ := strconv.Atoi(idCapitulo)
		idV, _ := strconv.Atoi(idVersiculo)
		livroCapVer(w, r, idL, idC, idV)
	} else {
		// Caso contrário, a URL não está no formato esperado
		http.NotFound(w, r)
	}
}

func livroPorID(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/bible_nvi")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// aqui fazemos a query passando o id
	rows, _ := db.Query("SELECT * from bible_nvi where book = ?", id)
	defer rows.Close()

	// criamos um slice vazio do tipo biblia
	var data []Biblia

	// fazemos um loop pelas linhas e adicionamos o conteudo ao slice
	for rows.Next() {
		var biblia Biblia
		rows.Scan(&biblia.Livro, &biblia.Capitulo, &biblia.Versiculo, &biblia.Texto)
		data = append(data, biblia)
	}

	// transformamos o slice em json
	json, _ := json.Marshal(data)

	// imprimimos o json
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func livroCap(w http.ResponseWriter, r *http.Request, idLivro int, idCap int ) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/bible_nvi")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// aqui fazemos a query passando o id
	rows, _ := db.Query("SELECT chapter, verse, text from bible_nvi where book = ? and chapter = ?", idLivro, idCap)
	defer rows.Close()

	// criamos um slice vazio do tipo capitulo
	var data []Capitulo

	// fazemos um loop pelas linhas e adicionamos o conteudo ao slice
	for rows.Next() {
        var cap Capitulo
		if err := rows.Scan(&cap.Capitulo, &cap.Versiculo, &cap.Texto); err != nil {
			log.Println(err)
			continue
		}
        data = append(data, cap)
    }
	/*var b Biblia
	db.QueryRow("SELECT * from bible_nvi where book = ?", id).Scan(&b.Livro) */

	// transformamos o slice em json
	json, _ := json.Marshal(data)

	// imprimimos o json
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func livroCapVer(w http.ResponseWriter, r *http.Request, idLivro int, idCap int, idVer int ) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/bible_nvi")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// aqui fazemos a query passando o id
	rows, _ := db.Query("SELECT chapter, verse, text from bible_nvi where book = ? and chapter = ? and verse = ?", idLivro, idCap, idVer)
	defer rows.Close()

	// criamos um slice vazio do tipo versiculo
	var data []Versiculo

	// fazemos um loop pelas linhas e adicionamos o conteudo ao slice
	for rows.Next() {
        var ver Versiculo
		if err := rows.Scan(&ver.Capitulo, &ver.Versiculo, &ver.Texto); err != nil {
			log.Println(err)
			continue
		}
        data = append(data, ver)
    }

	// transformamos o slice em json
	json, _ := json.Marshal(data)

	// imprimimos o json
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func todosOsLivros(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/bible_nvi")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, _ := db.Query("select * from bible_nvi")
	defer rows.Close()

	var data []Biblia
	for rows.Next() {
		var biblia Biblia
		rows.Scan(&biblia.Livro, &biblia.Capitulo, &biblia.Versiculo, &biblia.Texto)
		data = append(data, biblia)
	}

	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}