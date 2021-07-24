package messageCreator

import (
	"digi-bot/model"
	"fmt"
)

func CreateNormalPriceChangeMsg(product model.Product, newPrice int, oldPrice int) string {
	output := createHeader(product).
		Append(fmt.Sprintf("%s -> %s", Number(oldPrice).AddComma(), Number(newPrice).AddComma()))
	return output.ToString()
}

func CreatePreviewMsg(product model.Product) string {
	output := String(product.Name).
		Bold().
		AddNewLine()
	if product.Price != product.OldPrice {
		output = output.
			Append("قیمت: ").
			Append(fmt.Sprintf("%s -> %s",
				Number(product.OldPrice).
					AddComma().
					Strike(),
				Number(product.Price).
					AddComma()))
	} else {
		output = output.
			Append("قیمت: ").
			Append(Number(product.Price).AddComma().ToString())
	}

	output = output.
		AddNewLine().
		AddNewLine().
		Append("کالا با موفقیت ذخیره شد. برای اضافه کردن کالای جدید کافی است فقط آدرس آن را وارد کنید")

	return output.ToString()
}

func CreateNotAvailableMsg(product model.Product) string {
	output := createHeader(product).
		Append("ناموجود!")

	return output.ToString()
}

func createHeader(product model.Product) String {
	output := String(product.Name).
		Bold().
		ToLink(product.Url).
		AddNewLine()

	return output
}
