package assignment01bca

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type block struct {
	data     []*transaction
	hash     string
	prevHash string
	nonce    int
	root     *MerkleTree
}

type transaction struct {
	ID int
}

type blockchain struct {
	blocks []*block
}

func (b *block) calculateHash() string { //Function to calcuate hash
	bytes, _ := json.Marshal(b.data)
	// concatenate the dataset
	blockData := strconv.Itoa(b.nonce) + string(bytes) + string(b.root.RootNode.Data) //+ strconv.Atoi(b.prevHash) //Appending all the block data

	//Using SHA256 to calculate hash of the block
	hashval := sha256.New()
	hashval.Write([]byte(blockData))
	b.hash = hex.EncodeToString(hashval.Sum(nil))
	return b.hash
}

func (chain *blockchain) MineBlock(numZeros int, b *block) { //Function to mine the block
	min := 1000
	max := 9999
	nonce := rand.Intn(max-min) + min  //Making a random 4 digit nonce
	y := strings.Repeat("0", numZeros) //Setting number of trailing zeros to achieve target
	//Loop to keep on calculating the hash for random nonces until target is found
	for !strings.HasPrefix(b.hash, y) {
		nonce = rand.Intn(max-min) + min
		b.nonce = nonce
		b.calculateHash()

	}
	chain.AddBlock(b) //Add the mined block to the blockchain
}

func NewBlock(t []*transaction) *block { //Function to create a new block
	//Initializing all the transactions
	Block := &block{
		data: t,
	}

	//Calculating root hash of merkel tree for the transactions of current block
	node := Block.HashTransactions()
	var p *MerkleTree = new(MerkleTree)
	p.RootNode = node
	Block.root = p
	return Block
}

func (chain *blockchain) AddBlock(b *block) { //Function to add block to the blockchain
	prevBlock := chain.blocks[len(chain.blocks)-1] //Storing the previous block in the blockchain
	b.prevHash = prevBlock.hash                    //Storing the hash of the previous block
	chain.blocks = append(chain.blocks, b)         //Appending this block to the blockchain
}

func Genesis(b *block) *block { //function to initialize the genesis block

	return NewBlock(b.data)
}

func (chain *blockchain) DisplayBlocks() { //Displaying all the block data]
	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %s\n", block.prevHash)
		for _, b := range block.data {
			fmt.Printf("Data in Block: %v\n", b.ID)
		}
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("Nonce: %d\n", block.nonce)
		fmt.Printf("Root: %x\n", block.root.RootNode.Data)
		fmt.Printf("------------------------------------------------------------------------\n")
	}
}

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode { //function to make a merkle node and hash it
	newNode := MerkleNode{}

	if left == nil && right == nil {
		hashval := sha256.Sum256(data)
		newNode.Data = hashval[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hashVal := sha256.Sum256(prevHashes)
		newNode.Data = hashVal[:]
	}

	newNode.Left = left
	newNode.Right = right

	return &newNode
}

func NewMerkleTree(data [][]byte) *MerkleTree { //Function to make the merkle tree
	var parents []MerkleNode

	for float64(int(math.Log2(float64(len(data))))) != math.Log2(float64(len(data))) {
		data = append(data, data[len(data)-1])
	}

	for _, dat := range data {
		newNode := NewMerkleNode(nil, nil, dat)
		parents = append(parents, *newNode)
	}

	for i := 0; i < int(math.Log2(float64(len(data)))); i++ {
		var children []MerkleNode

		for j := 0; j < len(parents); j += 2 {
			newNode := NewMerkleNode(&parents[j], &parents[j+1], nil)
			children = append(children, *newNode)
		}

		parents = children
	}

	merkletree := MerkleTree{&parents[0]}

	return &merkletree
}

func (b *block) DisplayMerkelTree() {
	var txHashes [][]byte

	for _, tx := range b.data {
		txHashes = append(txHashes, tx.Serialize())
	}
	for _, tx := range txHashes {
		fmt.Printf("%x\n", tx)
	}
}

func (b *block) AddItem(item *transaction) []*transaction { //adding transaction to list of transactions
	b.data = append(b.data, item)
	return b.data
}

func (b *transaction) Serialize() []byte { //function to serialize
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)
	return res.Bytes()
}

func (b *block) HashTransactions() *MerkleNode { //function to serialize the transactions
	var totalhashes [][]byte

	for _, hashval := range b.data {
		totalhashes = append(totalhashes, hashval.Serialize())
	}
	tree := NewMerkleTree(totalhashes)
	return tree.RootNode
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func changeBlock(b *block, t *transaction, index int) {
	b.data[index] = t            //Changing the block
	node := b.HashTransactions() //Hashing the transactions again

	//Forming the tree again for new value of root node
	var p *MerkleTree = new(MerkleTree)
	p.RootNode = node
	b.root = p
}

func (chain *blockchain) verifyChain() {
	for _, block := range chain.blocks { //traversing the entire chain
		prevHash := block.hash           //storing curret blockHash
		newHash := block.calculateHash() //storing hash after recalculating

		if prevHash != newHash { //if recalculated hash not same, block has been changed
			fmt.Println("Chain Invalid")
		}
	}
}

func DisplayMerkelRec(m *MerkleNode, space int) {
	if m == nil {
		return
	}

	space = space + 2

	DisplayMerkelRec(m.Right, space)

	fmt.Println()

	for i := 2; i < space; i++ {
		fmt.Print(" ")
	}

	fmt.Printf("%x\n", m.Data)

	DisplayMerkelRec(m.Left, space)

}

func (b *block) DisplayMerkel() {

	fmt.Println("-------------------Merkle Tree-----------------------")

	DisplayMerkelRec(b.root.RootNode, 0)

	fmt.Println("-----------------------------------------------------")

}

func MainFunc() {
	//Creating a list of transactions
	item1 := &transaction{ID: 0}
	item2 := &transaction{ID: 1}
	item3 := &transaction{ID: 2}
	item4 := &transaction{ID: 3}
	item5 := &transaction{ID: 4}
	item6 := &transaction{ID: 5}

	//Adding item to first block
	items := []*transaction{}
	block1 := &block{data: items}
	block1.AddItem(item1)
	block1.AddItem(item2)

	//adding items to second block
	block2 := &block{data: items}
	block2.AddItem(item3)
	block2.AddItem(item4)
	// block2.AddItem(item5)
	// block2.AddItem(item6)
	// block2.AddItem(item3)
	// block2.AddItem(item4)
	// block2.AddItem(item5)
	// block2.AddItem(item6)

	//adding items to third block
	block3 := &block{data: items}
	block3.AddItem(item5)
	block3.AddItem(item6)

	//Making the genesis block
	gen := Genesis(block1)
	gen.prevHash = "0"
	gen.nonce = 0
	gen.calculateHash()

	//Initiailizing the chain
	chain := blockchain{[]*block{gen}}

	//Mining first block
	block4 := NewBlock(block2.data)
	chain.MineBlock(1, block4)

	//Mining second block
	block5 := NewBlock(block3.data)
	chain.MineBlock(1, block5)

	//Displayig the list after mining
	chain.DisplayBlocks()

	//Displaying Merkle Tree (block 2)
	block4.DisplayMerkel()

	//Changing a transaction of the block
	index := 1
	item7 := &transaction{ID: 8}
	changeBlock(block3, item7, index)

	//Verifying the chain after changing block
	chain.verifyChain()
}
