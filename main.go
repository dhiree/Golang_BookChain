package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// this is program create by itself
type Block struct {
	Pos       int
	Data      BookCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutData string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishData string `json:"publish_data"`
	ISBN        string `json:"isbn"`
}

type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

// this function genreate the hash for first block
func (b *Block) generateHash() {

	bytes, _ := json.Marshal(b.Data)

	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()

	hash.Write([]byte(data))

	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

// this is create a new hash to user prives timestram and priveHash
func CreateBlock(prevBlock *Block, checkoutitem BookCheckout) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.TimeStamp = time.Now().String()
	block.PrevHash = prevBlock.Hash
	block.generateHash() // this is the methord to create a new hash

	return block

}

// this is a struct methord   and dc represent the blockchain

// this is create a new block to chick privous hash is match
func (bc *Blockchain) AddBlock(data BookCheckout) {
	// here we add block in blockchain

	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prevBlock, data)
	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}

}

// here chick the prevHash == t0 Hash
// this function check the block is valide or not checking privash hash
func validBlock(block, prevBlock *Block) bool {

	if prevBlock.Hash != block.PrevHash {
		return false

	}
	if !block.validataHash(block.Hash) {
		return false
	}
	if prevBlock.Pos+1 != block.Pos {
		return false
	}
	return true
}

// validate block Hash is Currect  or not
func (b *Block) validataHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

// decode the request body and and add block in blockchain
func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutitem BookCheckout

	if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could Not Write block:%v", err)
		w.Write([]byte("could not write block"))

	}
	BlockChain.AddBlock(checkoutitem)
}

// create a new book and genrate a uniq id using md5hash
func newBook(w http.ResponseWriter, r *http.Request) {
	//work is the book struct
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Book not Create:%v", err)
		w.Write([]byte("could Not Create Book"))
		return
	}

	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishData+book.Author+book.Title)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could Not Marshal Paybal:%v", err)
		w.Write([]byte("Could Not Save Data"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// create a gensi block
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(jbytes))
}

func main() {
	BlockChain = NewBlockchain()
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	go func() {
		for _, block := range BlockChain.blocks {
			//fmt.Println("prev.hash: %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data:%x\n", string(bytes))
			fmt.Println()

		}
	}()

	log.Println("lisition on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
