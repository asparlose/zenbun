package zenbun_test

import (
	"fmt"

	"github.com/asparlose/zenbun"
)

func Example() {
	db := zenbun.New()

	db.Index("neko", "吾輩は猫である")
	db.Index("ningen", "人間失格")

	candidates := db.Find("猫")
	for _, c := range candidates {
		fmt.Printf("%s: %.4f\n", c.DocumentName, c.Score)
	}

	// Output: neko: 1.0000
}
