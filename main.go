package main

import (
	_ "encoding/csv"
	"encoding/json"
	"fmt"
	"flag"
	"io/ioutil"
	"sort"
	"strings"
	"strconv"
	"os"
)

type pizza struct {
	Name        string		 `json:"name"`
	Price       int          `json:"price"`
	PriceDelta  int          `json:-`
	Ingredients []string     `ingredients:"ingredients"`
}

type ingredient struct {
	Name  string
	Price int
}

func (this *pizza) Remove(ing ingredient){
	for idx, name := range this.Ingredients {
		if name == ing.Name {
			this.PriceDelta -= ing.Price
			this.Ingredients[idx] = this.Ingredients[len(this.Ingredients)-1]
			this.Ingredients = this.Ingredients[:len(this.Ingredients)-1]
			return
		}
	}
}

func (this *pizza) Add(ing ingredient){
	for _, name := range this.Ingredients {
		if name == ing.Name {
			return
		}
	}
	this.PriceDelta += ing.Price
	this.Ingredients = append(this.Ingredients, ing.Name)
}

func (this pizza) GetPrice() int{
	if this.PriceDelta < 0 {
		return this.Price
	} else {
		return this.Price + this.PriceDelta
	}
}

// Make a sortable pizza array
type PizzaArray []*pizza
func (this PizzaArray) Len() int {
	return len(this)
}

func (this PizzaArray) Less(i, j int) bool {
	if this[i].GetPrice() == this[j].GetPrice() {
		return len(this[i].Ingredients) > len(this[j].Ingredients)
	} else {
		return this[i].GetPrice() < this[j].GetPrice()
	}
}
func (this PizzaArray) Swap(i, j int){
	tmp := this[j]
	this[j] = this[i]
	this[i] = tmp
}

func main(){
	var (
		menuPath        string
		ingredientsPath string
		list            bool
		help            bool
		add             string
		remove          string
	)
	flag.StringVar(&ingredientsPath, "i", "", "File with the ingredients information")
	flag.StringVar(&menuPath, "m", "", "File with the menu information")
	flag.BoolVar(&list, "l", false, "List the ingredients and exit")
	flag.StringVar(&add, "a", "", "Comma separated list of indices of the ingredients to add (as shown with the -l option)")
	flag.StringVar(&remove, "r", "", "Comma separated list of indices of the ingredients to remove (as shown with the -l option)")
	flag.BoolVar(&help, "h", false, "Show this message and exit")
	usage := func() {
		fmt.Printf("Usage %s [OPTIONS]\n", os.Args[0])
		fmt.Printf("\n")
		fmt.Printf("OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Printf("\n")
    }
	flag.Parse()

	if help {
		usage()
		return
	}

	if ingredientsPath == "" || (menuPath == "" && !list){
		usage()
		return
	}

	var (
		ingredients   []ingredient
		ingredientIdx map[string]int
		pizzas		  PizzaArray
	)

	ingredientIdx = make(map[string]int)
	// Ingredients
	file, err := os.Open(ingredientsPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(content, &ingredients)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}



	if list {
		for i, ing := range ingredients {
			fmt.Printf("%2d. %v %d.%d€\n", i, ing.Name, ing.Price/100, ing.Price%100)
		}
		fmt.Printf("\n")
		return
	}

	for idx, item := range ingredients {
		ingredientIdx[item.Name] = idx
	}

	// Menu
	file, err = os.Open(menuPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	content, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(content, &pizzas)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if add != "" {
		i := 0
		for _, idx := range strings.Split(add, ",") {
			idx, err := strconv.ParseInt(idx, 10, 0)
			if err != nil {
				fmt.Printf("Could not convert %s to number.\n", idx)
				return
			}
			for _, pizza := range pizzas {
				pizza.Add(ingredients[idx])
			}
		}
		i++
	}
	if remove != "" {
		i := 0
		for _, idx := range strings.Split(remove, ",") {
			idx, err := strconv.ParseInt(idx, 10, 0)
			if err != nil {
				fmt.Printf("Could not convert %s to number.\n", idx)
				return
			}
			for _, pizza := range pizzas {
				pizza.Remove(ingredients[idx])
			}
		}
		i++
	}

	sort.Sort(pizzas)
	fmt.Printf("Pizzas:\n")
	for i := 0; i<len(pizzas) && i<4; i++ {
		fmt.Printf("%s: %d.%.2d€\n", pizzas[i].Name, pizzas[i].GetPrice()/100, pizzas[i].GetPrice()%100)
		fmt.Printf("\t%d.%.2d€ %d [", pizzas[i].PriceDelta/100, -pizzas[i].PriceDelta%100, len(pizzas[i].Ingredients))
		for _, name := range pizzas[i].Ingredients {
			fmt.Printf("%s,", name)
		}
		fmt.Printf("]\n\n")
	}
}
